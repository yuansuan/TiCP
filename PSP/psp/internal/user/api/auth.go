package api

import (
	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/http"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/errcode"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/jwt"
	"github.com/yuansuan/ticp/PSP/psp/internal/user/dto"
	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/ginutil"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/strutil"
)

// Login
//
//	@Summary		登录接口
//	@Description	登录接口
//	@Tags			用户-auth
//	@Accept			json
//	@Produce		json
//	@Param			param	body		dto.UserRequest				true	"入参"
//	@Success		200		{object}	dto.LoginSuccessResponse	"正常返回表示成功，否则返回错误码和错误信息"
//	@Router			/auth/login [post]
func (s *RouteService) Login(ctx *gin.Context) {
	var req dto.UserRequest
	if err := ctx.BindJSON(&req); err != nil {
		http.Errf(ctx, errcode.ErrFileFailToBindRequest, "failed to bind request, err: %v", err)
		return
	}

	// 检查当前系统的 license 是否过期
	//err := s.LicenseService.CheckLicenseExpired(ctx)
	//if err != nil {
	//	errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrAuthLicenseHasExpiredOrNotExist)
	//	return
	//}

	// 先进行ldap/ad域身份校验，如果校验开关未开启或用户名密码不匹配，则再进行系统内部用户校验
	needInternalAuth, err := s.AuthService.CheckLdapUserPass(ctx, req)
	if needInternalAuth {
		err = s.AuthService.CheckUserPass(ctx, req)
	}

	if err != nil {
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrAuthUserFailed)
		return
	}

	loginSuccess, err := s.UserService.LoginCheck(ctx, snowflake.MustParseString(req.Id).Int64(), req.Name)
	if err != nil {
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrAuthUserFailed)
		return
	}

	// 生成并设置token
	_, _, err = jwt.SetToken(snowflake.MustParseString(loginSuccess.User.Id).Int64(), loginSuccess.User.Name, ctx)
	if err != nil {
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrAuthUserFailed)
		return
	}

	ginutil.Success(ctx, loginSuccess)
}

// Logout
//
//	@Summary		登出接口
//	@Description	登出接口
//	@Tags			用户-auth
//	@Accept			json
//	@Produce		json
//	@Success		200
//	@Router			/auth/logout [post]
func (s *RouteService) Logout(ctx *gin.Context) {
	token, _ := ctx.Cookie(jwt.AccessToken)

	if strutil.IsEmpty(token) {
		http.Errf(ctx, errcode.ErrInvalidParam, "token can't empty")
		return
	}

	userID := ginutil.GetUserID(ctx)
	if userID == 0 {
		http.Errf(ctx, errcode.ErrInvalidParam, "userID can't empty")
		return
	}

	jwt.CleanWhiteListByToken(token)

	ginutil.Success(ctx, nil)
}

// PingLdap
//
//	@Summary		连接测试-ldap
//	@Description	连接测试-ldap
//	@Tags			用户-auth
//	@Accept			json
//	@Produce		json
//	@Param			ldapServer	query	string	true	"ldapServer地址"	default("")
//	@Success		200
//	@Router			/ping/ldap [get]
func (s *RouteService) PingLdap(ctx *gin.Context) {
	ldapServer := ctx.Query("ldapServer")
	if strutil.IsEmpty(ldapServer) {
		http.Errf(ctx, errcode.ErrInvalidParam, "ldapServer can't empty")
		return
	}

	conn, err := s.AuthService.GetConn(ldapServer)
	if err != nil {
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrLDAPConnectFailed)
		return
	}
	defer conn.Close()

}

// OnlineList
//
//	@Summary		登录用户列表
//	@Description	登录用户列表
//	@Tags			用户-auth
//	@Accept			json
//	@Produce		json
//	@Param			param	body		dto.OnlineListRequest		true	"入参"
//	@Success		200		{object}	dto.OnlineUserListResponse	"正常返回表示成功，否则返回错误码和错误信息"
//	@Router			/auth/onlineList [post]
func (s *RouteService) OnlineList(ctx *gin.Context) {
	var req dto.OnlineListRequest
	if err := ctx.BindJSON(&req); err != nil {
		http.Errf(ctx, errcode.ErrFileFailToBindRequest, "failed to bind request, err: %v", err)
		return
	}

	list, err := s.AuthService.GetOnlineList(ctx, req)

	if err != nil {
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrOnlineListFailed)
		return
	}

	ginutil.Success(ctx, list)

}

// OnlineListByUser
//
//	@Summary		登录用户列表-具体用户
//	@Description	登录用户列表-具体用户
//	@Tags			用户-auth
//	@Accept			json
//	@Produce		json
//	@Param			param	body		dto.OnlineListByUserRequest		true	"入参"
//	@Success		200		{object}	dto.OnlineUserInfoListResponse	"正常返回表示成功，否则返回错误码和错误信息"
//	@Router			/auth/onlineListByUser [post]
func (s *RouteService) OnlineListByUser(ctx *gin.Context) {

	var req dto.OnlineListByUserRequest
	if err := ctx.BindJSON(&req); err != nil {
		http.Errf(ctx, errcode.ErrFileFailToBindRequest, "failed to bind request, err: %v", err)
		return
	}

	list, err := s.AuthService.GetOnlineListByUser(ctx, req)

	if err != nil {
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrOnlineListFailed)
		return
	}

	ginutil.Success(ctx, list)
}

// OfflineByUserName
//
//	@Summary		离线接口
//	@Description	离线接口
//	@Tags			用户-auth
//	@Accept			json
//	@Produce		json
//	@Param			param	body	dto.OfflineByUserNameRequest	true	"入参"
//	@Success		200
//	@Router			/auth/offlineByUserName [post]
func (s *RouteService) OfflineByUserName(ctx *gin.Context) {

	var req dto.OfflineByUserNameRequest
	if err := ctx.BindJSON(&req); err != nil {
		http.Errf(ctx, errcode.ErrFileFailToBindRequest, "failed to bind request, err: %v", err)
		return
	}

	if len(req.UserNameList) == 0 {
		http.Errf(ctx, errcode.ErrInvalidParam, "userNameList can't empty")
		return
	}

	for _, name := range req.UserNameList {
		jwt.CleanWhiteListByUserName(name)
	}

	ginutil.Success(ctx, nil)
}

// OfflineByJti
//
//	@Summary		离线接口
//	@Description	离线接口
//	@Tags			用户-auth
//	@Accept			json
//	@Produce		json
//	@Param			param	body	dto.OfflineByJtiRequest	true	"入参"
//	@Success		200
//	@Router			/auth/offlineByJti [post]
func (s *RouteService) OfflineByJti(ctx *gin.Context) {

	var req dto.OfflineByJtiRequest
	if err := ctx.BindJSON(&req); err != nil {
		http.Errf(ctx, errcode.ErrFileFailToBindRequest, "failed to bind request, err: %v", err)
		return
	}

	if len(req.UserName) == 0 {
		http.Errf(ctx, errcode.ErrInvalidParam, "userName can't empty")
		return
	}

	if len(req.JtiList) == 0 {
		http.Errf(ctx, errcode.ErrInvalidParam, "userNameList can't empty")
		return
	}

	for _, jti := range req.JtiList {
		jwt.CleanWhiteListByJti(jti, req.UserName)
	}

	ginutil.Success(ctx, nil)
}
