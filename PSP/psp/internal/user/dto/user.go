package dto

import (
	"github.com/yuansuan/ticp/PSP/psp/pkg/xtype"
)

type UserIDRequest struct {
	Id string `json:"id" form:"id"` // 用户id
}

type UserRequest struct {
	Id       string `json:"id" form:"id"`             // 用户id
	Name     string `json:"name" form:"name"`         // 用户名
	Password string `json:"password" form:"password"` // 密码
}

type LoginSuccessResponse struct {
	User *UserInfo `json:"user"`
}

type UpdatePassRequest struct {
	Name        string `json:"name"`        // 用户名
	Password    string `json:"password"`    // 现密码
	NewPassword string `json:"newPassword"` // 新密码
}

type LoginUserResponse struct {
	UserId     string `json:"user_id"`   // 用户id
	UserName   string `json:"user_name"` // 用户名
	NumSession int64  `json:"num_session"`
	CreatedAt  int64  `json:"created_at"` // 创建时间
	UpdatedAt  int64  `json:"updated_at"` // 修改时间
	Mail       string `json:"mail"`       // 邮箱
}

type UserAddRequest struct {
	Name          string  `json:"name"`           // 用户名
	Password      string  `json:"password"`       // 密码
	Email         string  `json:"email"`          // 邮箱
	Mobile        string  `json:"mobile"`         //手机号
	RealName      string  `json:"real_name"`      // 真实姓名
	Roles         []int64 `json:"roles"`          // 所属角色信息
	EnableOpenapi bool    `json:"enable_openapi"` // 是否启动openapi
}

type QueryByCondRequest struct {
	Enabled bool        `json:"enabled"` // 是否启用
	Query   string      `json:"query"`   // 查询条件
	Order   string      `json:"order"`   // 排序字段
	Desc    bool        `json:"desc"`    //排序方式
	Page    *xtype.Page `json:"page"`    // 分页参数
}

type UserInfo struct {
	Id            string  `json:"id"`             // 用户id
	Name          string  `json:"name"`           // 用户名
	Email         string  `json:"email"`          // 邮箱
	Mobile        string  `json:"mobile"`         // 手机号
	Enabled       bool    `json:"enabled"`        // 是否启用
	IsInternal    bool    `json:"is_internal"`    // 是否内部用户
	CreatedAt     int64   `json:"created_at"`     // 创建时间
	RealName      string  `json:"real_name"`      // 真实姓名
	ApproveStatus int     `json:"approve_status"` // 审批状态
	Roles         []int64 `json:"roles"`          // 所属角色id
	EnableOpenapi bool    `json:"enable_openapi"` // 启用openapi
}

type UserListResponse struct {
	Success bool        `json:"success"`
	Total   int64       `json:"total"`    // 总数
	UserObj []*UserInfo `json:"user_obj"` // 用户信息
}

type UserOptionResponse struct {
	Key   string `json:"key"`
	Title string `json:"title"`
}

type UserDetailResponse struct {
	UserInfo           *UserInfo `json:"user_info"`           // 用户信息
	RoleInfo           []*Role   `json:"role"`                // 角色信息
	Perm               *Perm     `json:"perm"`                // 权限
	Conf               *Conf     `json:"conf"`                // 配置信息
	OpenapiCertificate string    `json:"openapi_certificate"` // openapi凭证
}

type Conf struct {
	LdapEnable    bool `json:"ldap_enable"`
	OpenapiSwitch bool `json:"openapi_switch"`
}

type Perm struct {
	LocalApp       []*CustomPerm `json:"local_app"`
	CloudApp       []*CustomPerm `json:"cloud_app"`
	VisualSoftware []*CustomPerm `json:"visual_software"`
	System         []*CustomPerm `json:"system"`
}

type CustomPerm struct {
	Id         int64  `json:"id"`
	Name       string `json:"name"`
	Key        string `json:"key"`
	ExternalId string `json:"external_id"`
	Has        bool   `json:"has"`
}

type Role struct {
	Id   int64  `json:"id"`   // 用户id
	Name string `json:"name"` // 用户名
}

type UserDataConfigResponse struct {
	HomeDir string `json:"homeDir"` // 家目录
}

type UserUpdateRequest struct {
	Id            string  `json:"id"`             // 用户id
	Email         string  `json:"email"`          // 邮箱
	Mobile        string  `json:"mobile"`         // 手机号
	Roles         []int64 `json:"roles"`          // 角色
	EnableOpenapi bool    `json:"enable_openapi"` // 启用openapi
}

type OnlineListRequest struct {
	FilterName string      `json:"filter_name"` // 用户名筛选
	SortByAsc  bool        `json:"sort_by"`     // 是否升序
	OrderBy    string      `json:"order_by"`    // 排序字段 name
	Page       *xtype.Page `json:"page"`        // 分页参数
}

type OnlineListByUserRequest struct {
	Name string      `json:"user_name"` // 用户名
	Page *xtype.Page `json:"page"`      // 分页参数
}

type OfflineByUserNameRequest struct {
	UserNameList []string `json:"user_name_list"`
}

type OfflineByJtiRequest struct {
	JtiList  []string `json:"jti_list"`
	UserName string   `json:"user_name"`
}

type LdapUserDefRoleRequest struct {
	RoleID int64 `json:"role_id"` // 角色id
}

type ResetPassword struct {
	UserID string `json:"user_id"`
}

type GenOpenapiCertificateRequest struct {
	UserID string `json:"user_id"`
}
