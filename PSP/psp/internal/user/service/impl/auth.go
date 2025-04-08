package impl

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/go-ldap/ldap/v3"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	boot "github.com/yuansuan/ticp/common/go-kit/gin-boot"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"google.golang.org/grpc/status"

	"github.com/yuansuan/ticp/PSP/psp/internal/common"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/cache"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/errcode"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/jwt"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/rbac"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/storage"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/sysconfig"
	"github.com/yuansuan/ticp/PSP/psp/internal/user/config"
	"github.com/yuansuan/ticp/PSP/psp/internal/user/dao"
	"github.com/yuansuan/ticp/PSP/psp/internal/user/dao/model"
	"github.com/yuansuan/ticp/PSP/psp/internal/user/dto"
	"github.com/yuansuan/ticp/PSP/psp/internal/user/service"
	"github.com/yuansuan/ticp/PSP/psp/internal/user/service/client"
	"github.com/yuansuan/ticp/PSP/psp/internal/user/util"
	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
	"github.com/yuansuan/ticp/PSP/psp/pkg/tracelog"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/strutil"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/timeutil"
	"github.com/yuansuan/ticp/PSP/psp/pkg/xtype"
)

type AuthServiceImpl struct {
	userDao dao.UserDao
}

func NewAuthService() service.AuthService {
	return &AuthServiceImpl{
		userDao: dao.NewUserDaoImpl(),
	}
}

func (s AuthServiceImpl) CheckUserPass(ctx context.Context, in dto.UserRequest) error {
	id, err := s.userDao.GetIdByName(in.Name)
	if err != nil || id < 0 {
		return status.Error(errcode.ErrUserNotExist, "user not found")
	}
	user, exists, err := s.userDao.Get(id)
	if err != nil {
		return status.Error(errcode.ErrInternalServer, err.Error())
	}
	if !exists {
		return status.Error(errcode.ErrUserNotExist, "user not found")
	}

	//if user.IsDeleted {
	//	return status.Error(errcode.ErrUserDeleted, "user deleted")
	//}

	if util.PasswdCrypto(in.Password) != user.Password {
		return status.Error(errcode.ErrUserInvalidPassword, "invalid password")
	}

	return nil
}

func (s AuthServiceImpl) CheckLdapUserPass(ctx context.Context, req dto.UserRequest) (bool, error) {
	conf := config.GetConfig()

	if !conf.LdapConf.Enable {
		return true, nil
	}

	user, err := s.userDao.GetUserByName(req.Name)
	if err != nil && status.Code(err) != errcode.ErrUserNotFound {
		return false, nil
	}

	if user != nil && user.IsInternal {
		// 如果登录用户名为系统内置用户(超级管理员)，走内部登录逻辑
		return true, nil
	}

	logger := logging.Default()
	// 连接ldap
	var conn *ldap.Conn
	if conf.LdapConf.Encryption == "ssl" {
		conn, err = ldap.DialTLS("tcp",
			conf.LdapConf.Server, &tls.Config{InsecureSkipVerify: true},
		)
	} else {
		conn, err = s.GetConn(conf.LdapConf.Server)

		if err == nil && conf.LdapConf.Encryption == "starttls" {
			err = conn.StartTLS(&tls.Config{InsecureSkipVerify: true})
		}
	}

	if err != nil {
		return false, status.Error(errcode.ErrLDAPConnectFailed, err.Error())
	}
	defer conn.Close()

	// bind with bind dn and password
	err = conn.Bind(conf.LdapConf.AdminBindDn, conf.LdapConf.AdminBindPassword)
	if err != nil {
		logger.Errorf("LDAP bind admin user failed with username %s", conf.LdapConf.AdminBindDn)
		if ldap.IsErrorWithCode(err, ldap.LDAPResultInvalidCredentials) {
			return false, status.Error(errcode.ErrAuthUserFailed, err.Error())
		}
		return false, status.Error(errcode.ErrLDAPConnectFailed, err.Error())
	}

	// filter with login user
	filter := fmt.Sprintf("(%s=%s)", conf.LdapConf.UID, req.Name)
	if conf.LdapConf.UserFilter != "" {
		filter = fmt.Sprintf("(&(%s)%s)", conf.LdapConf.UserFilter, filter)
	}

	// 查询登录用户
	searchRequest := ldap.NewSearchRequest(
		conf.LdapConf.BaseDN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0, 0, false,
		filter,
		[]string{},
		nil)

	sr, err := conn.Search(searchRequest)

	if err != nil {
		logger.Errorf("ldapsearch-filter: %v", filter)
		return false, status.Error(errcode.ErrAuthUserFailed, err.Error())
	}

	if len(sr.Entries) != 1 {
		logger.Errorf("User does not exist or too many entries returned, search result len:%v", len(sr.Entries))
		return false, status.Error(errcode.ErrAuthUserFailed, "User does not exist or too many entries returned")
	}

	// 身份检验
	err = conn.Bind(sr.Entries[0].DN, req.Password)
	if err != nil {
		logger.Errorf("LDAP bind login user failed with userdn %s", sr.Entries[0].DN)
		if ldap.IsErrorWithCode(err, ldap.LDAPResultInvalidCredentials) {
			return true, status.Error(errcode.ErrAuthUserFailed, err.Error())
		}
		return false, status.Error(errcode.ErrLDAPConnectFailed, err.Error())
	}

	return false, s.syncLdapUser(ctx, req.Name)
}

// syncLdapUser 处理ldap用户，user表不存在则新增并赋予默认角色
func (s AuthServiceImpl) syncLdapUser(ctx context.Context, userName string) error {
	id, err := s.userDao.GetIdByName(userName)

	if err != nil && status.Code(err) == errcode.ErrUserNotFound {
		defRoleRsp, err := client.GetInstance().SysConfig.GetRBACDefaultRoleId(ctx, &sysconfig.GetRBACDefaultRoleIdRequest{})
		if err != nil {
			return err
		}
		if defRoleRsp != nil && defRoleRsp.RoleId == 0 {
			return status.Error(errcode.ErrUserDefaultRoleNotExist, "")
		}

		if userName == common.ProjectFolderPath || userName == common.PublicFolderPath {
			return status.Error(errcode.ErrUserNameCollideBuiltin, "")
		}

		var node *snowflake.Node
		node, err = snowflake.GetInstance()
		if err != nil {
			logging.Default().Errorf("new snowflake node err: %v", err)
			return err
		}
		id = node.Generate().Int64()

		user := model.User{
			Id:      id,
			Name:    userName,
			Enabled: true,
		}
		_, err = s.userDao.InsertUser(user)
		if err != nil {
			return err
		}

		tracelog.Info(ctx, fmt.Sprintf("sync ldap user success, user[%+v]", user))
		_, err = client.GetInstance().Role.AddObjectRoles(ctx, &rbac.ObjectRoles{
			Id: &rbac.ObjectID{
				Id:   snowflake.ID(user.Id).String(),
				Type: rbac.ObjectType_USER,
			},
			Roles: []int64{defRoleRsp.RoleId},
		})

		_, _ = client.GetInstance().Storage.InitUserHome(ctx, &storage.InitUserHomeReq{
			UserName: userName,
		})
	} else {
		return err
	}

	return nil
}

//
//// ldapSearch ldap查询，根据不同ldap服务配置不同(开启身份验证)，调用ldapSearch前可能需要先进行conn.Bind操作
//func ldapSearch(conn *ldap.Conn, baseDN string, attr config.ExtraAttrs, filter string) (map[string]string, error) {
//
//	attributes := []string{attr.UIDKey, attr.EmailKey, attr.MobileKey, attr.RealNameKey}
//	request := ldap.NewSearchRequest(baseDN, ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false, filter, attributes, nil)
//
//	cur, err := conn.Search(request)
//	if err != nil || cur.Entries == nil || len(cur.Entries) <= 0 {
//		return nil, err
//	}
//
//	entry := cur.Entries[0]
//
//	v := reflect.ValueOf(attr)
//	t := reflect.TypeOf(attr)
//	result := make(map[string]string, 0)
//
//	for i := 0; i < t.NumField(); i++ {
//		key := v.Field(i)
//		if key.IsValid() {
//			result[key.String()] = entry.GetAttributeValue(key.String())
//		}
//	}
//	return result, nil
//}

func (s AuthServiceImpl) GetConn(serverAddr string) (*ldap.Conn, error) {
	return ldap.DialURL(serverAddr, ldap.DialWithTLSConfig(&tls.Config{InsecureSkipVerify: true}))
}

func (s AuthServiceImpl) GetOnlineList(ctx *gin.Context, req dto.OnlineListRequest) (*dto.OnlineUserListResponse, error) {

	var redisSuccess bool
	searchResultList := make([]string, 0)

	// 从redis获取存活token
	redisClient := boot.Middleware.DefaultRedis()
	if redisClient != nil {
		// PSP:WHITE_LIST_TOKEN:*{filterName}*:*
		// PSP:WHITE_LIST_TOKEN:**:*
		key := jwt.GetWhiteListKey(fmt.Sprintf("*%s*", req.FilterName), "*")

		result := redisClient.Scan(0, key, 0)
		if result.Err() == nil {
			redisSuccess = true
		}
		iterator := result.Iterator()
		for iterator.Next() {
			searchResultList = append(searchResultList, iterator.Val())
		}
	}

	if !redisSuccess {
		// 如果redis服务不可用，降级使用go-cache
		for key := range cache.Cache.Items() {
			// PSP:WHITE_LIST_TOKEN:(.*{filterName}.*):
			regexpStr := regexp.MustCompile(jwt.GetWhiteListKey(fmt.Sprintf(".*%s.*", req.FilterName), ""))
			if regexpStr.MatchString(key) {
				searchResultList = append(searchResultList, key)
			}
		}
	}

	// 汇总，将查找到的searchResult转化成map,key:{userName} value:{*dto.OnlineUserResponse}
	// PSP:WHITE_LIST_TOKEN:(.*):
	regexpStr := regexp.MustCompile(jwt.GetWhiteListKey("(.*)", ""))
	userMap := make(map[string]*dto.OnlineUserResponse, 0)

	// token = PSP:WHITE_LIST_TOKEN:{userName}:{jti}
	for _, token := range searchResultList {
		match := regexpStr.FindStringSubmatch(token)

		if len(match) > 1 {
			username := match[1]
			if v, ok := userMap[username]; ok {
				v.Count += 1
			} else {
				userMap[username] = &dto.OnlineUserResponse{
					Name:  username,
					Count: 1,
				}
			}
		}
	}

	// map转list，便于排序分页
	list := make([]*dto.OnlineUserResponse, 0)
	for _, v := range userMap {
		list = append(list, v)
	}

	if len(list) == 0 {
		return &dto.OnlineUserListResponse{
			List: list,
			Page: &xtype.PageResp{
				Index: 1,
				Size:  10,
				Total: 0,
			},
		}, nil
	}

	// 排序
	if req.OrderBy == "name" {
		sort.Slice(list, func(i, j int) bool {
			if req.SortByAsc {
				return list[i].Name < list[j].Name
			} else {
				return list[i].Name > list[j].Name
			}
		})
	}

	// 分页
	pageSize := int64(10)
	pageIndex := int64(1)
	if req.Page != nil {
		pageSize = req.Page.Size
		pageIndex = req.Page.Index
	}

	return &dto.OnlineUserListResponse{
		List: getOnlineUserPageList(list, pageSize, pageIndex),
		Page: &xtype.PageResp{
			Index: pageIndex,
			Size:  pageSize,
			Total: int64(len(userMap)),
		},
	}, nil
}

func getOnlineUserPageList[T *dto.OnlineUserResponse | *dto.OnlineUserInfoResponse](list []T, pageSize int64, pageIndex int64) []T {
	length := int64(len(list))
	startIndex := (pageIndex - 1) * pageSize
	endIndex := pageIndex * pageSize

	if startIndex >= length {
		return []T{}
	}

	if endIndex >= length {
		endIndex = length
	}

	return list[startIndex:endIndex]
}

func (s AuthServiceImpl) GetOnlineListByUser(ctx *gin.Context, req dto.OnlineListByUserRequest) (*dto.OnlineUserInfoListResponse, error) {

	var redisSuccess bool
	searchResultList := make([]string, 0)

	// 从redis获取存活token
	redisClient := boot.Middleware.DefaultRedis()
	if redisClient != nil {
		// PSP:WHITE_LIST_TOKEN:*{filterName}*:*
		// PSP:WHITE_LIST_TOKEN:**:*
		key := jwt.GetWhiteListKey(req.Name, "*")

		result := redisClient.Scan(0, key, 0)
		if result.Err() == nil {
			redisSuccess = true
		}
		iterator := result.Iterator()
		for iterator.Next() {
			searchResultList = append(searchResultList, iterator.Val())
		}
	}

	if !redisSuccess {
		// 如果redis服务不可用，降级使用go-cache
		for key := range cache.Cache.Items() {
			// PSP:WHITE_LIST_TOKEN:yskj:2131231
			regexpStr := regexp.MustCompile(jwt.GetWhiteListKey(req.Name, ""))
			if regexpStr.MatchString(key) {
				searchResultList = append(searchResultList, key)
			}
		}
	}

	list := make([]*dto.OnlineUserInfoResponse, 0)
	for _, token := range searchResultList {
		userInfo := &dto.OnlineUserInfoResponse{
			Jti: token[strings.LastIndex(token, ":")+1:],
		}

		if redisSuccess {
			val := redisClient.Get(token).Val()
			if strutil.IsEmpty(val) {
				continue
			}

			userInfo.IP = redisClient.Get(token).Val()

			userInfo.ExpireTime = timeutil.DefaultFormatTime(time.Now().Add(redisClient.TTL(token).Val()))
		} else {
			val, expireTime, ok := cache.Cache.GetWithExpiration(token)
			if !ok {
				continue
			}

			userInfo.IP = val.(string)
			userInfo.ExpireTime = timeutil.DefaultFormatTime(expireTime)
		}
		list = append(list, userInfo)
	}

	// 分页
	pageSize := int64(10)
	pageIndex := int64(1)
	if req.Page != nil {
		pageSize = req.Page.Size
		pageIndex = req.Page.Index
	}

	return &dto.OnlineUserInfoListResponse{
		List: getOnlineUserPageList(list, pageSize, pageIndex),
		Page: &xtype.PageResp{
			Index: pageIndex,
			Size:  pageSize,
			Total: int64(len(list)),
		},
	}, nil

}
