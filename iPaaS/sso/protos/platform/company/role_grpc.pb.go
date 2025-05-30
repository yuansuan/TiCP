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

// RoleServiceClient is the client API for RoleService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RoleServiceClient interface {
	// 创建系统角色 （运营后台使用）
	CreateSysRole(ctx context.Context, in *CreateSysRoleRequest, opts ...grpc.CallOption) (*CreateSysRoleResponse, error)
	// 创建企业角色
	CreateCompanyRole(ctx context.Context, in *CreateCompanyRoleRequest, opts ...grpc.CallOption) (*CreateCompanyRoleResponse, error)
	// 修改角色
	ModifyRole(ctx context.Context, in *ModifyRoleRequest, opts ...grpc.CallOption) (*ModifyRoleResponse, error)
	// 获取企业角色列表
	GetCompanyRoleList(ctx context.Context, in *GetCompanyRoleListRequest, opts ...grpc.CallOption) (*GetCompanyRoleListResponse, error)
	// 获取系统角色列表
	GetSysRoleList(ctx context.Context, in *GetSysRoleListRequest, opts ...grpc.CallOption) (*GetSysRoleListResponse, error)
	// 给角色赋权限
	GrantPermissionToRole(ctx context.Context, in *GrantPermissionToRoleRequest, opts ...grpc.CallOption) (*GrantPermissionToRoleResponse, error)
	// 通过角色ID获取角色所属权限
	GetPermissionByRoleID(ctx context.Context, in *GetPermissionByRoleIDRequest, opts ...grpc.CallOption) (*GetPermissionByRoleIDResponse, error)
}

type roleServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewRoleServiceClient(cc grpc.ClientConnInterface) RoleServiceClient {
	return &roleServiceClient{cc}
}

func (c *roleServiceClient) CreateSysRole(ctx context.Context, in *CreateSysRoleRequest, opts ...grpc.CallOption) (*CreateSysRoleResponse, error) {
	out := new(CreateSysRoleResponse)
	err := c.cc.Invoke(ctx, "/company.RoleService/CreateSysRole", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *roleServiceClient) CreateCompanyRole(ctx context.Context, in *CreateCompanyRoleRequest, opts ...grpc.CallOption) (*CreateCompanyRoleResponse, error) {
	out := new(CreateCompanyRoleResponse)
	err := c.cc.Invoke(ctx, "/company.RoleService/CreateCompanyRole", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *roleServiceClient) ModifyRole(ctx context.Context, in *ModifyRoleRequest, opts ...grpc.CallOption) (*ModifyRoleResponse, error) {
	out := new(ModifyRoleResponse)
	err := c.cc.Invoke(ctx, "/company.RoleService/ModifyRole", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *roleServiceClient) GetCompanyRoleList(ctx context.Context, in *GetCompanyRoleListRequest, opts ...grpc.CallOption) (*GetCompanyRoleListResponse, error) {
	out := new(GetCompanyRoleListResponse)
	err := c.cc.Invoke(ctx, "/company.RoleService/GetCompanyRoleList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *roleServiceClient) GetSysRoleList(ctx context.Context, in *GetSysRoleListRequest, opts ...grpc.CallOption) (*GetSysRoleListResponse, error) {
	out := new(GetSysRoleListResponse)
	err := c.cc.Invoke(ctx, "/company.RoleService/GetSysRoleList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *roleServiceClient) GrantPermissionToRole(ctx context.Context, in *GrantPermissionToRoleRequest, opts ...grpc.CallOption) (*GrantPermissionToRoleResponse, error) {
	out := new(GrantPermissionToRoleResponse)
	err := c.cc.Invoke(ctx, "/company.RoleService/GrantPermissionToRole", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *roleServiceClient) GetPermissionByRoleID(ctx context.Context, in *GetPermissionByRoleIDRequest, opts ...grpc.CallOption) (*GetPermissionByRoleIDResponse, error) {
	out := new(GetPermissionByRoleIDResponse)
	err := c.cc.Invoke(ctx, "/company.RoleService/GetPermissionByRoleID", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RoleServiceServer is the server API for RoleService service.
// All implementations must embed UnimplementedRoleServiceServer
// for forward compatibility
type RoleServiceServer interface {
	// 创建系统角色 （运营后台使用）
	CreateSysRole(context.Context, *CreateSysRoleRequest) (*CreateSysRoleResponse, error)
	// 创建企业角色
	CreateCompanyRole(context.Context, *CreateCompanyRoleRequest) (*CreateCompanyRoleResponse, error)
	// 修改角色
	ModifyRole(context.Context, *ModifyRoleRequest) (*ModifyRoleResponse, error)
	// 获取企业角色列表
	GetCompanyRoleList(context.Context, *GetCompanyRoleListRequest) (*GetCompanyRoleListResponse, error)
	// 获取系统角色列表
	GetSysRoleList(context.Context, *GetSysRoleListRequest) (*GetSysRoleListResponse, error)
	// 给角色赋权限
	GrantPermissionToRole(context.Context, *GrantPermissionToRoleRequest) (*GrantPermissionToRoleResponse, error)
	// 通过角色ID获取角色所属权限
	GetPermissionByRoleID(context.Context, *GetPermissionByRoleIDRequest) (*GetPermissionByRoleIDResponse, error)
	mustEmbedUnimplementedRoleServiceServer()
}

// UnimplementedRoleServiceServer must be embedded to have forward compatible implementations.
type UnimplementedRoleServiceServer struct {
}

func (UnimplementedRoleServiceServer) CreateSysRole(context.Context, *CreateSysRoleRequest) (*CreateSysRoleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateSysRole not implemented")
}
func (UnimplementedRoleServiceServer) CreateCompanyRole(context.Context, *CreateCompanyRoleRequest) (*CreateCompanyRoleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateCompanyRole not implemented")
}
func (UnimplementedRoleServiceServer) ModifyRole(context.Context, *ModifyRoleRequest) (*ModifyRoleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ModifyRole not implemented")
}
func (UnimplementedRoleServiceServer) GetCompanyRoleList(context.Context, *GetCompanyRoleListRequest) (*GetCompanyRoleListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCompanyRoleList not implemented")
}
func (UnimplementedRoleServiceServer) GetSysRoleList(context.Context, *GetSysRoleListRequest) (*GetSysRoleListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSysRoleList not implemented")
}
func (UnimplementedRoleServiceServer) GrantPermissionToRole(context.Context, *GrantPermissionToRoleRequest) (*GrantPermissionToRoleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GrantPermissionToRole not implemented")
}
func (UnimplementedRoleServiceServer) GetPermissionByRoleID(context.Context, *GetPermissionByRoleIDRequest) (*GetPermissionByRoleIDResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPermissionByRoleID not implemented")
}
func (UnimplementedRoleServiceServer) mustEmbedUnimplementedRoleServiceServer() {}

// UnsafeRoleServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RoleServiceServer will
// result in compilation errors.
type UnsafeRoleServiceServer interface {
	mustEmbedUnimplementedRoleServiceServer()
}

func RegisterRoleServiceServer(s grpc.ServiceRegistrar, srv RoleServiceServer) {
	s.RegisterService(&RoleService_ServiceDesc, srv)
}

func _RoleService_CreateSysRole_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateSysRoleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RoleServiceServer).CreateSysRole(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/company.RoleService/CreateSysRole",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RoleServiceServer).CreateSysRole(ctx, req.(*CreateSysRoleRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RoleService_CreateCompanyRole_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateCompanyRoleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RoleServiceServer).CreateCompanyRole(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/company.RoleService/CreateCompanyRole",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RoleServiceServer).CreateCompanyRole(ctx, req.(*CreateCompanyRoleRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RoleService_ModifyRole_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ModifyRoleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RoleServiceServer).ModifyRole(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/company.RoleService/ModifyRole",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RoleServiceServer).ModifyRole(ctx, req.(*ModifyRoleRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RoleService_GetCompanyRoleList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetCompanyRoleListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RoleServiceServer).GetCompanyRoleList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/company.RoleService/GetCompanyRoleList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RoleServiceServer).GetCompanyRoleList(ctx, req.(*GetCompanyRoleListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RoleService_GetSysRoleList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetSysRoleListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RoleServiceServer).GetSysRoleList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/company.RoleService/GetSysRoleList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RoleServiceServer).GetSysRoleList(ctx, req.(*GetSysRoleListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RoleService_GrantPermissionToRole_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GrantPermissionToRoleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RoleServiceServer).GrantPermissionToRole(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/company.RoleService/GrantPermissionToRole",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RoleServiceServer).GrantPermissionToRole(ctx, req.(*GrantPermissionToRoleRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RoleService_GetPermissionByRoleID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetPermissionByRoleIDRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RoleServiceServer).GetPermissionByRoleID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/company.RoleService/GetPermissionByRoleID",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RoleServiceServer).GetPermissionByRoleID(ctx, req.(*GetPermissionByRoleIDRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// RoleService_ServiceDesc is the grpc.ServiceDesc for RoleService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var RoleService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "company.RoleService",
	HandlerType: (*RoleServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateSysRole",
			Handler:    _RoleService_CreateSysRole_Handler,
		},
		{
			MethodName: "CreateCompanyRole",
			Handler:    _RoleService_CreateCompanyRole_Handler,
		},
		{
			MethodName: "ModifyRole",
			Handler:    _RoleService_ModifyRole_Handler,
		},
		{
			MethodName: "GetCompanyRoleList",
			Handler:    _RoleService_GetCompanyRoleList_Handler,
		},
		{
			MethodName: "GetSysRoleList",
			Handler:    _RoleService_GetSysRoleList_Handler,
		},
		{
			MethodName: "GrantPermissionToRole",
			Handler:    _RoleService_GrantPermissionToRole_Handler,
		},
		{
			MethodName: "GetPermissionByRoleID",
			Handler:    _RoleService_GetPermissionByRoleID_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "protos/platform/company/role.proto",
}
