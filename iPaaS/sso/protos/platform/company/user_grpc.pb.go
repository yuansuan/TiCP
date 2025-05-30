// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package company

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// UserServiceClient is the client API for UserService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type UserServiceClient interface {
	// 企业邀请用户
	InviteUser(ctx context.Context, in *InviteUserRequest, opts ...grpc.CallOption) (*InviteUserResponse, error)
	// 用户确认企业邀请
	ConfirmInvite(ctx context.Context, in *ConfirmInviteRequest, opts ...grpc.CallOption) (*ConfirmInviteResponse, error)
	// 用户邀请信息
	GetUserInviteInfo(ctx context.Context, in *GetUserInviteInfoRequest, opts ...grpc.CallOption) (*GetUserInviteInfoResponse, error)
	// 用户信息修改
	UserModify(ctx context.Context, in *UserModifyRequest, opts ...grpc.CallOption) (*UserModifyResponse, error)
	// 用户初始化
	UserInit(ctx context.Context, in *UserInitRequest, opts ...grpc.CallOption) (*UserInitResponse, error)
	// 获取用户企业角色及权限
	GetUserRoleAndPermisson(ctx context.Context, in *GetUserRoleAndPermissonRequest, opts ...grpc.CallOption) (*GetUserRoleAndPermissonResponse, error)
	// 验证用户权限
	CheckUserPermisson(ctx context.Context, in *CheckUserPermissonRequest, opts ...grpc.CallOption) (*CheckUserPermissonResponse, error)
	// 获取用户企业信息
	GetUserCompanyInfo(ctx context.Context, in *GetUserCompanyInfoRequest, opts ...grpc.CallOption) (*GetUserCompanyInfoResponse, error)
	// 获取用户信息
	GetUserInfo(ctx context.Context, in *GetUserInfoRequest, opts ...grpc.CallOption) (*GetUserInfoResponse, error)
	// 用户列表查询
	UserListQuery(ctx context.Context, in *UserListQueryRequest, opts ...grpc.CallOption) (*UserListQueryResponse, error)
	// 添加用户备注 (当前实现不是最好实现，未来需要独立的服务来支持)
	AddUserRemark(ctx context.Context, in *AddUserRemarkRequest, opts ...grpc.CallOption) (*AddUserRemarkResponse, error)
	// 批量邀请用户（未注册手机号自动注册）
	BatchInviteUser(ctx context.Context, in *BatchInviteUserRequest, opts ...grpc.CallOption) (*BatchInviteUserResponse, error)
	// 更新用户自定义码
	UpdateUserFeCode(ctx context.Context, in *UpdateUserFeCodeRequest, opts ...grpc.CallOption) (*UpdateUserFeCodeResponse, error)
	// 获取个人用户产品开通信息
	GetUserProductList(ctx context.Context, in *GetUserProductListRequest, opts ...grpc.CallOption) (*GetUserProductListResponse, error)
	// 为个人用户添加产品
	AddProductToUser(ctx context.Context, in *AddProductToUserRequest, opts ...grpc.CallOption) (*AddProductToUserResponse, error)
	// 移除个人用户产品
	RemoveProductFromUser(ctx context.Context, in *RemoveProductFromUserRequest, opts ...grpc.CallOption) (*RemoveProductFromUserResponse, error)
	// 移除企业产品
	CheckUserProduct(ctx context.Context, in *CheckUserProductRequest, opts ...grpc.CallOption) (*CheckUserProductResponse, error)
}

type userServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewUserServiceClient(cc grpc.ClientConnInterface) UserServiceClient {
	return &userServiceClient{cc}
}

func (c *userServiceClient) InviteUser(ctx context.Context, in *InviteUserRequest, opts ...grpc.CallOption) (*InviteUserResponse, error) {
	out := new(InviteUserResponse)
	err := c.cc.Invoke(ctx, "/company.UserService/InviteUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) ConfirmInvite(ctx context.Context, in *ConfirmInviteRequest, opts ...grpc.CallOption) (*ConfirmInviteResponse, error) {
	out := new(ConfirmInviteResponse)
	err := c.cc.Invoke(ctx, "/company.UserService/ConfirmInvite", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) GetUserInviteInfo(ctx context.Context, in *GetUserInviteInfoRequest, opts ...grpc.CallOption) (*GetUserInviteInfoResponse, error) {
	out := new(GetUserInviteInfoResponse)
	err := c.cc.Invoke(ctx, "/company.UserService/GetUserInviteInfo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) UserModify(ctx context.Context, in *UserModifyRequest, opts ...grpc.CallOption) (*UserModifyResponse, error) {
	out := new(UserModifyResponse)
	err := c.cc.Invoke(ctx, "/company.UserService/UserModify", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) UserInit(ctx context.Context, in *UserInitRequest, opts ...grpc.CallOption) (*UserInitResponse, error) {
	out := new(UserInitResponse)
	err := c.cc.Invoke(ctx, "/company.UserService/UserInit", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) GetUserRoleAndPermisson(ctx context.Context, in *GetUserRoleAndPermissonRequest, opts ...grpc.CallOption) (*GetUserRoleAndPermissonResponse, error) {
	out := new(GetUserRoleAndPermissonResponse)
	err := c.cc.Invoke(ctx, "/company.UserService/GetUserRoleAndPermisson", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) CheckUserPermisson(ctx context.Context, in *CheckUserPermissonRequest, opts ...grpc.CallOption) (*CheckUserPermissonResponse, error) {
	out := new(CheckUserPermissonResponse)
	err := c.cc.Invoke(ctx, "/company.UserService/CheckUserPermisson", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) GetUserCompanyInfo(ctx context.Context, in *GetUserCompanyInfoRequest, opts ...grpc.CallOption) (*GetUserCompanyInfoResponse, error) {
	out := new(GetUserCompanyInfoResponse)
	err := c.cc.Invoke(ctx, "/company.UserService/GetUserCompanyInfo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) GetUserInfo(ctx context.Context, in *GetUserInfoRequest, opts ...grpc.CallOption) (*GetUserInfoResponse, error) {
	out := new(GetUserInfoResponse)
	err := c.cc.Invoke(ctx, "/company.UserService/GetUserInfo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) UserListQuery(ctx context.Context, in *UserListQueryRequest, opts ...grpc.CallOption) (*UserListQueryResponse, error) {
	out := new(UserListQueryResponse)
	err := c.cc.Invoke(ctx, "/company.UserService/UserListQuery", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) AddUserRemark(ctx context.Context, in *AddUserRemarkRequest, opts ...grpc.CallOption) (*AddUserRemarkResponse, error) {
	out := new(AddUserRemarkResponse)
	err := c.cc.Invoke(ctx, "/company.UserService/AddUserRemark", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) BatchInviteUser(ctx context.Context, in *BatchInviteUserRequest, opts ...grpc.CallOption) (*BatchInviteUserResponse, error) {
	out := new(BatchInviteUserResponse)
	err := c.cc.Invoke(ctx, "/company.UserService/BatchInviteUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) UpdateUserFeCode(ctx context.Context, in *UpdateUserFeCodeRequest, opts ...grpc.CallOption) (*UpdateUserFeCodeResponse, error) {
	out := new(UpdateUserFeCodeResponse)
	err := c.cc.Invoke(ctx, "/company.UserService/UpdateUserFeCode", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) GetUserProductList(ctx context.Context, in *GetUserProductListRequest, opts ...grpc.CallOption) (*GetUserProductListResponse, error) {
	out := new(GetUserProductListResponse)
	err := c.cc.Invoke(ctx, "/company.UserService/GetUserProductList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) AddProductToUser(ctx context.Context, in *AddProductToUserRequest, opts ...grpc.CallOption) (*AddProductToUserResponse, error) {
	out := new(AddProductToUserResponse)
	err := c.cc.Invoke(ctx, "/company.UserService/AddProductToUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) RemoveProductFromUser(ctx context.Context, in *RemoveProductFromUserRequest, opts ...grpc.CallOption) (*RemoveProductFromUserResponse, error) {
	out := new(RemoveProductFromUserResponse)
	err := c.cc.Invoke(ctx, "/company.UserService/RemoveProductFromUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) CheckUserProduct(ctx context.Context, in *CheckUserProductRequest, opts ...grpc.CallOption) (*CheckUserProductResponse, error) {
	out := new(CheckUserProductResponse)
	err := c.cc.Invoke(ctx, "/company.UserService/CheckUserProduct", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UserServiceServer is the server API for UserService service.
// All implementations must embed UnimplementedUserServiceServer
// for forward compatibility
type UserServiceServer interface {
	// 企业邀请用户
	InviteUser(context.Context, *InviteUserRequest) (*InviteUserResponse, error)
	// 用户确认企业邀请
	ConfirmInvite(context.Context, *ConfirmInviteRequest) (*ConfirmInviteResponse, error)
	// 用户邀请信息
	GetUserInviteInfo(context.Context, *GetUserInviteInfoRequest) (*GetUserInviteInfoResponse, error)
	// 用户信息修改
	UserModify(context.Context, *UserModifyRequest) (*UserModifyResponse, error)
	// 用户初始化
	UserInit(context.Context, *UserInitRequest) (*UserInitResponse, error)
	// 获取用户企业角色及权限
	GetUserRoleAndPermisson(context.Context, *GetUserRoleAndPermissonRequest) (*GetUserRoleAndPermissonResponse, error)
	// 验证用户权限
	CheckUserPermisson(context.Context, *CheckUserPermissonRequest) (*CheckUserPermissonResponse, error)
	// 获取用户企业信息
	GetUserCompanyInfo(context.Context, *GetUserCompanyInfoRequest) (*GetUserCompanyInfoResponse, error)
	// 获取用户信息
	GetUserInfo(context.Context, *GetUserInfoRequest) (*GetUserInfoResponse, error)
	// 用户列表查询
	UserListQuery(context.Context, *UserListQueryRequest) (*UserListQueryResponse, error)
	// 添加用户备注 (当前实现不是最好实现，未来需要独立的服务来支持)
	AddUserRemark(context.Context, *AddUserRemarkRequest) (*AddUserRemarkResponse, error)
	// 批量邀请用户（未注册手机号自动注册）
	BatchInviteUser(context.Context, *BatchInviteUserRequest) (*BatchInviteUserResponse, error)
	// 更新用户自定义码
	UpdateUserFeCode(context.Context, *UpdateUserFeCodeRequest) (*UpdateUserFeCodeResponse, error)
	// 获取个人用户产品开通信息
	GetUserProductList(context.Context, *GetUserProductListRequest) (*GetUserProductListResponse, error)
	// 为个人用户添加产品
	AddProductToUser(context.Context, *AddProductToUserRequest) (*AddProductToUserResponse, error)
	// 移除个人用户产品
	RemoveProductFromUser(context.Context, *RemoveProductFromUserRequest) (*RemoveProductFromUserResponse, error)
	// 移除企业产品
	CheckUserProduct(context.Context, *CheckUserProductRequest) (*CheckUserProductResponse, error)
	mustEmbedUnimplementedUserServiceServer()
}

// UnimplementedUserServiceServer must be embedded to have forward compatible implementations.
type UnimplementedUserServiceServer struct {
}

func (UnimplementedUserServiceServer) InviteUser(context.Context, *InviteUserRequest) (*InviteUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method InviteUser not implemented")
}
func (UnimplementedUserServiceServer) ConfirmInvite(context.Context, *ConfirmInviteRequest) (*ConfirmInviteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ConfirmInvite not implemented")
}
func (UnimplementedUserServiceServer) GetUserInviteInfo(context.Context, *GetUserInviteInfoRequest) (*GetUserInviteInfoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserInviteInfo not implemented")
}
func (UnimplementedUserServiceServer) UserModify(context.Context, *UserModifyRequest) (*UserModifyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UserModify not implemented")
}
func (UnimplementedUserServiceServer) UserInit(context.Context, *UserInitRequest) (*UserInitResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UserInit not implemented")
}
func (UnimplementedUserServiceServer) GetUserRoleAndPermisson(context.Context, *GetUserRoleAndPermissonRequest) (*GetUserRoleAndPermissonResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserRoleAndPermisson not implemented")
}
func (UnimplementedUserServiceServer) CheckUserPermisson(context.Context, *CheckUserPermissonRequest) (*CheckUserPermissonResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CheckUserPermisson not implemented")
}
func (UnimplementedUserServiceServer) GetUserCompanyInfo(context.Context, *GetUserCompanyInfoRequest) (*GetUserCompanyInfoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserCompanyInfo not implemented")
}
func (UnimplementedUserServiceServer) GetUserInfo(context.Context, *GetUserInfoRequest) (*GetUserInfoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserInfo not implemented")
}
func (UnimplementedUserServiceServer) UserListQuery(context.Context, *UserListQueryRequest) (*UserListQueryResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UserListQuery not implemented")
}
func (UnimplementedUserServiceServer) AddUserRemark(context.Context, *AddUserRemarkRequest) (*AddUserRemarkResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddUserRemark not implemented")
}
func (UnimplementedUserServiceServer) BatchInviteUser(context.Context, *BatchInviteUserRequest) (*BatchInviteUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BatchInviteUser not implemented")
}
func (UnimplementedUserServiceServer) UpdateUserFeCode(context.Context, *UpdateUserFeCodeRequest) (*UpdateUserFeCodeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateUserFeCode not implemented")
}
func (UnimplementedUserServiceServer) GetUserProductList(context.Context, *GetUserProductListRequest) (*GetUserProductListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserProductList not implemented")
}
func (UnimplementedUserServiceServer) AddProductToUser(context.Context, *AddProductToUserRequest) (*AddProductToUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddProductToUser not implemented")
}
func (UnimplementedUserServiceServer) RemoveProductFromUser(context.Context, *RemoveProductFromUserRequest) (*RemoveProductFromUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveProductFromUser not implemented")
}
func (UnimplementedUserServiceServer) CheckUserProduct(context.Context, *CheckUserProductRequest) (*CheckUserProductResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CheckUserProduct not implemented")
}
func (UnimplementedUserServiceServer) mustEmbedUnimplementedUserServiceServer() {}

// UnsafeUserServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to UserServiceServer will
// result in compilation errors.
type UnsafeUserServiceServer interface {
	mustEmbedUnimplementedUserServiceServer()
}

func RegisterUserServiceServer(s grpc.ServiceRegistrar, srv UserServiceServer) {
	s.RegisterService(&UserService_ServiceDesc, srv)
}

func _UserService_InviteUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InviteUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).InviteUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/company.UserService/InviteUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).InviteUser(ctx, req.(*InviteUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_ConfirmInvite_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ConfirmInviteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).ConfirmInvite(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/company.UserService/ConfirmInvite",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).ConfirmInvite(ctx, req.(*ConfirmInviteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_GetUserInviteInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserInviteInfoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).GetUserInviteInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/company.UserService/GetUserInviteInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).GetUserInviteInfo(ctx, req.(*GetUserInviteInfoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_UserModify_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserModifyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).UserModify(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/company.UserService/UserModify",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).UserModify(ctx, req.(*UserModifyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_UserInit_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserInitRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).UserInit(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/company.UserService/UserInit",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).UserInit(ctx, req.(*UserInitRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_GetUserRoleAndPermisson_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserRoleAndPermissonRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).GetUserRoleAndPermisson(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/company.UserService/GetUserRoleAndPermisson",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).GetUserRoleAndPermisson(ctx, req.(*GetUserRoleAndPermissonRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_CheckUserPermisson_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CheckUserPermissonRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).CheckUserPermisson(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/company.UserService/CheckUserPermisson",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).CheckUserPermisson(ctx, req.(*CheckUserPermissonRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_GetUserCompanyInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserCompanyInfoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).GetUserCompanyInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/company.UserService/GetUserCompanyInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).GetUserCompanyInfo(ctx, req.(*GetUserCompanyInfoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_GetUserInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserInfoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).GetUserInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/company.UserService/GetUserInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).GetUserInfo(ctx, req.(*GetUserInfoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_UserListQuery_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserListQueryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).UserListQuery(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/company.UserService/UserListQuery",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).UserListQuery(ctx, req.(*UserListQueryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_AddUserRemark_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddUserRemarkRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).AddUserRemark(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/company.UserService/AddUserRemark",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).AddUserRemark(ctx, req.(*AddUserRemarkRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_BatchInviteUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BatchInviteUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).BatchInviteUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/company.UserService/BatchInviteUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).BatchInviteUser(ctx, req.(*BatchInviteUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_UpdateUserFeCode_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateUserFeCodeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).UpdateUserFeCode(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/company.UserService/UpdateUserFeCode",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).UpdateUserFeCode(ctx, req.(*UpdateUserFeCodeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_GetUserProductList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserProductListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).GetUserProductList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/company.UserService/GetUserProductList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).GetUserProductList(ctx, req.(*GetUserProductListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_AddProductToUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddProductToUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).AddProductToUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/company.UserService/AddProductToUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).AddProductToUser(ctx, req.(*AddProductToUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_RemoveProductFromUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RemoveProductFromUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).RemoveProductFromUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/company.UserService/RemoveProductFromUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).RemoveProductFromUser(ctx, req.(*RemoveProductFromUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_CheckUserProduct_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CheckUserProductRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).CheckUserProduct(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/company.UserService/CheckUserProduct",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).CheckUserProduct(ctx, req.(*CheckUserProductRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// UserService_ServiceDesc is the grpc.ServiceDesc for UserService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var UserService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "company.UserService",
	HandlerType: (*UserServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "InviteUser",
			Handler:    _UserService_InviteUser_Handler,
		},
		{
			MethodName: "ConfirmInvite",
			Handler:    _UserService_ConfirmInvite_Handler,
		},
		{
			MethodName: "GetUserInviteInfo",
			Handler:    _UserService_GetUserInviteInfo_Handler,
		},
		{
			MethodName: "UserModify",
			Handler:    _UserService_UserModify_Handler,
		},
		{
			MethodName: "UserInit",
			Handler:    _UserService_UserInit_Handler,
		},
		{
			MethodName: "GetUserRoleAndPermisson",
			Handler:    _UserService_GetUserRoleAndPermisson_Handler,
		},
		{
			MethodName: "CheckUserPermisson",
			Handler:    _UserService_CheckUserPermisson_Handler,
		},
		{
			MethodName: "GetUserCompanyInfo",
			Handler:    _UserService_GetUserCompanyInfo_Handler,
		},
		{
			MethodName: "GetUserInfo",
			Handler:    _UserService_GetUserInfo_Handler,
		},
		{
			MethodName: "UserListQuery",
			Handler:    _UserService_UserListQuery_Handler,
		},
		{
			MethodName: "AddUserRemark",
			Handler:    _UserService_AddUserRemark_Handler,
		},
		{
			MethodName: "BatchInviteUser",
			Handler:    _UserService_BatchInviteUser_Handler,
		},
		{
			MethodName: "UpdateUserFeCode",
			Handler:    _UserService_UpdateUserFeCode_Handler,
		},
		{
			MethodName: "GetUserProductList",
			Handler:    _UserService_GetUserProductList_Handler,
		},
		{
			MethodName: "AddProductToUser",
			Handler:    _UserService_AddProductToUser_Handler,
		},
		{
			MethodName: "RemoveProductFromUser",
			Handler:    _UserService_RemoveProductFromUser_Handler,
		},
		{
			MethodName: "CheckUserProduct",
			Handler:    _UserService_CheckUserProduct_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "protos/platform/company/user.proto",
}
