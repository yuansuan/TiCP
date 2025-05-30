// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.21.12
// source: proto/account_bill/billing.proto

package account_bill

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

const (
	BillingService_BillingCreate_FullMethodName             = "/account_bill.BillingService/BillingCreate"
	BillingService_BillingCreateV2_FullMethodName           = "/account_bill.BillingService/BillingCreateV2"
	BillingService_BillingGet_FullMethodName                = "/account_bill.BillingService/BillingGet"
	BillingService_BillingMessage_FullMethodName            = "/account_bill.BillingService/BillingMessage"
	BillingService_BillingOverview_FullMethodName           = "/account_bill.BillingService/BillingOverview"
	BillingService_BillingQuery_FullMethodName              = "/account_bill.BillingService/BillingQuery"
	BillingService_PlatformBillingOverview_FullMethodName   = "/account_bill.BillingService/PlatformBillingOverview"
	BillingService_PlatformBillingQuery_FullMethodName      = "/account_bill.BillingService/PlatformBillingQuery"
	BillingService_PlatformBillingStat_FullMethodName       = "/account_bill.BillingService/PlatformBillingStat"
	BillingService_GetBillingByOutBizs_FullMethodName       = "/account_bill.BillingService/GetBillingByOutBizs"
	BillingService_UserCost_FullMethodName                  = "/account_bill.BillingService/UserCost"
	BillingService_GetProjectCost_FullMethodName            = "/account_bill.BillingService/GetProjectCost"
	BillingService_GetCompanyLastFewDaysCost_FullMethodName = "/account_bill.BillingService/GetCompanyLastFewDaysCost"
	BillingService_ListBillingByIds_FullMethodName          = "/account_bill.BillingService/ListBillingByIds"
)

// BillingServiceClient is the client API for BillingService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type BillingServiceClient interface {
	// 计费单创建
	BillingCreate(ctx context.Context, in *BillingCreateRequest, opts ...grpc.CallOption) (*BillingCreateReply, error)
	// 计费单创建 企业价格不存在则使用个人商品价格
	BillingCreateV2(ctx context.Context, in *BillingCreateRequest, opts ...grpc.CallOption) (*BillingCreateReply, error)
	// 计费单查询
	BillingGet(ctx context.Context, in *BillingGetRequest, opts ...grpc.CallOption) (*BillingGetReply, error)
	// 计费单更新
	BillingMessage(ctx context.Context, in *BillingMessageRequest, opts ...grpc.CallOption) (*BillingMessageReply, error)
	// 账单总览（API调用）
	BillingOverview(ctx context.Context, in *BillingOverviewRequest, opts ...grpc.CallOption) (*BillingOverviewResponse, error)
	// 账单明细（API调用）
	BillingQuery(ctx context.Context, in *BillingQueryRequest, opts ...grpc.CallOption) (*BillingQueryResponse, error)
	// 账单总览（运营后台调用）
	PlatformBillingOverview(ctx context.Context, in *PlatformBillingOverviewRequest, opts ...grpc.CallOption) (*PlatformBillingOverviewResponse, error)
	// 账单明细（运营后台调用）
	PlatformBillingQuery(ctx context.Context, in *PlatformBillingQueryRequest, opts ...grpc.CallOption) (*PlatformBillingQueryResponse, error)
	// 账单统计（运营后台调用）
	PlatformBillingStat(ctx context.Context, in *PlatformBillingStatRequest, opts ...grpc.CallOption) (*PlatformBillingStatResponse, error)
	// 根据外部业务ID获取计费信息
	GetBillingByOutBizs(ctx context.Context, in *GetBillingByOutBizsRequest, opts ...grpc.CallOption) (*GetBillingByOutBizsResponse, error)
	// 查询一段时间内用户消费 包含已冻结资金
	UserCost(ctx context.Context, in *UserCostRequest, opts ...grpc.CallOption) (*UserCostResponse, error)
	// 按ProjectID查询消耗总金额
	GetProjectCost(ctx context.Context, in *ProjectCostRequest, opts ...grpc.CallOption) (*ProjectCostResponse, error)
	// 获取企业最近N天的消费总额和平均消费金额
	GetCompanyLastFewDaysCost(ctx context.Context, in *GetCompanyLastFewDaysCostRequest, opts ...grpc.CallOption) (*GetCompanyLastFewDaysCostResponse, error)
	// 根据账单ID查询计费信息
	ListBillingByIds(ctx context.Context, in *ListBillingByIdsRequest, opts ...grpc.CallOption) (*ListBillingByIdsResponse, error)
}

type billingServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewBillingServiceClient(cc grpc.ClientConnInterface) BillingServiceClient {
	return &billingServiceClient{cc}
}

func (c *billingServiceClient) BillingCreate(ctx context.Context, in *BillingCreateRequest, opts ...grpc.CallOption) (*BillingCreateReply, error) {
	out := new(BillingCreateReply)
	err := c.cc.Invoke(ctx, BillingService_BillingCreate_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *billingServiceClient) BillingCreateV2(ctx context.Context, in *BillingCreateRequest, opts ...grpc.CallOption) (*BillingCreateReply, error) {
	out := new(BillingCreateReply)
	err := c.cc.Invoke(ctx, BillingService_BillingCreateV2_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *billingServiceClient) BillingGet(ctx context.Context, in *BillingGetRequest, opts ...grpc.CallOption) (*BillingGetReply, error) {
	out := new(BillingGetReply)
	err := c.cc.Invoke(ctx, BillingService_BillingGet_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *billingServiceClient) BillingMessage(ctx context.Context, in *BillingMessageRequest, opts ...grpc.CallOption) (*BillingMessageReply, error) {
	out := new(BillingMessageReply)
	err := c.cc.Invoke(ctx, BillingService_BillingMessage_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *billingServiceClient) BillingOverview(ctx context.Context, in *BillingOverviewRequest, opts ...grpc.CallOption) (*BillingOverviewResponse, error) {
	out := new(BillingOverviewResponse)
	err := c.cc.Invoke(ctx, BillingService_BillingOverview_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *billingServiceClient) BillingQuery(ctx context.Context, in *BillingQueryRequest, opts ...grpc.CallOption) (*BillingQueryResponse, error) {
	out := new(BillingQueryResponse)
	err := c.cc.Invoke(ctx, BillingService_BillingQuery_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *billingServiceClient) PlatformBillingOverview(ctx context.Context, in *PlatformBillingOverviewRequest, opts ...grpc.CallOption) (*PlatformBillingOverviewResponse, error) {
	out := new(PlatformBillingOverviewResponse)
	err := c.cc.Invoke(ctx, BillingService_PlatformBillingOverview_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *billingServiceClient) PlatformBillingQuery(ctx context.Context, in *PlatformBillingQueryRequest, opts ...grpc.CallOption) (*PlatformBillingQueryResponse, error) {
	out := new(PlatformBillingQueryResponse)
	err := c.cc.Invoke(ctx, BillingService_PlatformBillingQuery_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *billingServiceClient) PlatformBillingStat(ctx context.Context, in *PlatformBillingStatRequest, opts ...grpc.CallOption) (*PlatformBillingStatResponse, error) {
	out := new(PlatformBillingStatResponse)
	err := c.cc.Invoke(ctx, BillingService_PlatformBillingStat_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *billingServiceClient) GetBillingByOutBizs(ctx context.Context, in *GetBillingByOutBizsRequest, opts ...grpc.CallOption) (*GetBillingByOutBizsResponse, error) {
	out := new(GetBillingByOutBizsResponse)
	err := c.cc.Invoke(ctx, BillingService_GetBillingByOutBizs_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *billingServiceClient) UserCost(ctx context.Context, in *UserCostRequest, opts ...grpc.CallOption) (*UserCostResponse, error) {
	out := new(UserCostResponse)
	err := c.cc.Invoke(ctx, BillingService_UserCost_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *billingServiceClient) GetProjectCost(ctx context.Context, in *ProjectCostRequest, opts ...grpc.CallOption) (*ProjectCostResponse, error) {
	out := new(ProjectCostResponse)
	err := c.cc.Invoke(ctx, BillingService_GetProjectCost_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *billingServiceClient) GetCompanyLastFewDaysCost(ctx context.Context, in *GetCompanyLastFewDaysCostRequest, opts ...grpc.CallOption) (*GetCompanyLastFewDaysCostResponse, error) {
	out := new(GetCompanyLastFewDaysCostResponse)
	err := c.cc.Invoke(ctx, BillingService_GetCompanyLastFewDaysCost_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *billingServiceClient) ListBillingByIds(ctx context.Context, in *ListBillingByIdsRequest, opts ...grpc.CallOption) (*ListBillingByIdsResponse, error) {
	out := new(ListBillingByIdsResponse)
	err := c.cc.Invoke(ctx, BillingService_ListBillingByIds_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BillingServiceServer is the server API for BillingService service.
// All implementations must embed UnimplementedBillingServiceServer
// for forward compatibility
type BillingServiceServer interface {
	// 计费单创建
	BillingCreate(context.Context, *BillingCreateRequest) (*BillingCreateReply, error)
	// 计费单创建 企业价格不存在则使用个人商品价格
	BillingCreateV2(context.Context, *BillingCreateRequest) (*BillingCreateReply, error)
	// 计费单查询
	BillingGet(context.Context, *BillingGetRequest) (*BillingGetReply, error)
	// 计费单更新
	BillingMessage(context.Context, *BillingMessageRequest) (*BillingMessageReply, error)
	// 账单总览（API调用）
	BillingOverview(context.Context, *BillingOverviewRequest) (*BillingOverviewResponse, error)
	// 账单明细（API调用）
	BillingQuery(context.Context, *BillingQueryRequest) (*BillingQueryResponse, error)
	// 账单总览（运营后台调用）
	PlatformBillingOverview(context.Context, *PlatformBillingOverviewRequest) (*PlatformBillingOverviewResponse, error)
	// 账单明细（运营后台调用）
	PlatformBillingQuery(context.Context, *PlatformBillingQueryRequest) (*PlatformBillingQueryResponse, error)
	// 账单统计（运营后台调用）
	PlatformBillingStat(context.Context, *PlatformBillingStatRequest) (*PlatformBillingStatResponse, error)
	// 根据外部业务ID获取计费信息
	GetBillingByOutBizs(context.Context, *GetBillingByOutBizsRequest) (*GetBillingByOutBizsResponse, error)
	// 查询一段时间内用户消费 包含已冻结资金
	UserCost(context.Context, *UserCostRequest) (*UserCostResponse, error)
	// 按ProjectID查询消耗总金额
	GetProjectCost(context.Context, *ProjectCostRequest) (*ProjectCostResponse, error)
	// 获取企业最近N天的消费总额和平均消费金额
	GetCompanyLastFewDaysCost(context.Context, *GetCompanyLastFewDaysCostRequest) (*GetCompanyLastFewDaysCostResponse, error)
	// 根据账单ID查询计费信息
	ListBillingByIds(context.Context, *ListBillingByIdsRequest) (*ListBillingByIdsResponse, error)
	mustEmbedUnimplementedBillingServiceServer()
}

// UnimplementedBillingServiceServer must be embedded to have forward compatible implementations.
type UnimplementedBillingServiceServer struct {
}

func (UnimplementedBillingServiceServer) BillingCreate(context.Context, *BillingCreateRequest) (*BillingCreateReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BillingCreate not implemented")
}
func (UnimplementedBillingServiceServer) BillingCreateV2(context.Context, *BillingCreateRequest) (*BillingCreateReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BillingCreateV2 not implemented")
}
func (UnimplementedBillingServiceServer) BillingGet(context.Context, *BillingGetRequest) (*BillingGetReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BillingGet not implemented")
}
func (UnimplementedBillingServiceServer) BillingMessage(context.Context, *BillingMessageRequest) (*BillingMessageReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BillingMessage not implemented")
}
func (UnimplementedBillingServiceServer) BillingOverview(context.Context, *BillingOverviewRequest) (*BillingOverviewResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BillingOverview not implemented")
}
func (UnimplementedBillingServiceServer) BillingQuery(context.Context, *BillingQueryRequest) (*BillingQueryResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BillingQuery not implemented")
}
func (UnimplementedBillingServiceServer) PlatformBillingOverview(context.Context, *PlatformBillingOverviewRequest) (*PlatformBillingOverviewResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PlatformBillingOverview not implemented")
}
func (UnimplementedBillingServiceServer) PlatformBillingQuery(context.Context, *PlatformBillingQueryRequest) (*PlatformBillingQueryResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PlatformBillingQuery not implemented")
}
func (UnimplementedBillingServiceServer) PlatformBillingStat(context.Context, *PlatformBillingStatRequest) (*PlatformBillingStatResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PlatformBillingStat not implemented")
}
func (UnimplementedBillingServiceServer) GetBillingByOutBizs(context.Context, *GetBillingByOutBizsRequest) (*GetBillingByOutBizsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetBillingByOutBizs not implemented")
}
func (UnimplementedBillingServiceServer) UserCost(context.Context, *UserCostRequest) (*UserCostResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UserCost not implemented")
}
func (UnimplementedBillingServiceServer) GetProjectCost(context.Context, *ProjectCostRequest) (*ProjectCostResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetProjectCost not implemented")
}
func (UnimplementedBillingServiceServer) GetCompanyLastFewDaysCost(context.Context, *GetCompanyLastFewDaysCostRequest) (*GetCompanyLastFewDaysCostResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCompanyLastFewDaysCost not implemented")
}
func (UnimplementedBillingServiceServer) ListBillingByIds(context.Context, *ListBillingByIdsRequest) (*ListBillingByIdsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListBillingByIds not implemented")
}
func (UnimplementedBillingServiceServer) mustEmbedUnimplementedBillingServiceServer() {}

// UnsafeBillingServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to BillingServiceServer will
// result in compilation errors.
type UnsafeBillingServiceServer interface {
	mustEmbedUnimplementedBillingServiceServer()
}

func RegisterBillingServiceServer(s grpc.ServiceRegistrar, srv BillingServiceServer) {
	s.RegisterService(&BillingService_ServiceDesc, srv)
}

func _BillingService_BillingCreate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BillingCreateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BillingServiceServer).BillingCreate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BillingService_BillingCreate_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BillingServiceServer).BillingCreate(ctx, req.(*BillingCreateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BillingService_BillingCreateV2_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BillingCreateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BillingServiceServer).BillingCreateV2(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BillingService_BillingCreateV2_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BillingServiceServer).BillingCreateV2(ctx, req.(*BillingCreateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BillingService_BillingGet_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BillingGetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BillingServiceServer).BillingGet(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BillingService_BillingGet_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BillingServiceServer).BillingGet(ctx, req.(*BillingGetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BillingService_BillingMessage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BillingMessageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BillingServiceServer).BillingMessage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BillingService_BillingMessage_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BillingServiceServer).BillingMessage(ctx, req.(*BillingMessageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BillingService_BillingOverview_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BillingOverviewRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BillingServiceServer).BillingOverview(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BillingService_BillingOverview_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BillingServiceServer).BillingOverview(ctx, req.(*BillingOverviewRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BillingService_BillingQuery_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BillingQueryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BillingServiceServer).BillingQuery(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BillingService_BillingQuery_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BillingServiceServer).BillingQuery(ctx, req.(*BillingQueryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BillingService_PlatformBillingOverview_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PlatformBillingOverviewRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BillingServiceServer).PlatformBillingOverview(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BillingService_PlatformBillingOverview_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BillingServiceServer).PlatformBillingOverview(ctx, req.(*PlatformBillingOverviewRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BillingService_PlatformBillingQuery_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PlatformBillingQueryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BillingServiceServer).PlatformBillingQuery(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BillingService_PlatformBillingQuery_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BillingServiceServer).PlatformBillingQuery(ctx, req.(*PlatformBillingQueryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BillingService_PlatformBillingStat_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PlatformBillingStatRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BillingServiceServer).PlatformBillingStat(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BillingService_PlatformBillingStat_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BillingServiceServer).PlatformBillingStat(ctx, req.(*PlatformBillingStatRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BillingService_GetBillingByOutBizs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetBillingByOutBizsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BillingServiceServer).GetBillingByOutBizs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BillingService_GetBillingByOutBizs_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BillingServiceServer).GetBillingByOutBizs(ctx, req.(*GetBillingByOutBizsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BillingService_UserCost_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserCostRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BillingServiceServer).UserCost(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BillingService_UserCost_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BillingServiceServer).UserCost(ctx, req.(*UserCostRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BillingService_GetProjectCost_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ProjectCostRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BillingServiceServer).GetProjectCost(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BillingService_GetProjectCost_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BillingServiceServer).GetProjectCost(ctx, req.(*ProjectCostRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BillingService_GetCompanyLastFewDaysCost_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetCompanyLastFewDaysCostRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BillingServiceServer).GetCompanyLastFewDaysCost(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BillingService_GetCompanyLastFewDaysCost_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BillingServiceServer).GetCompanyLastFewDaysCost(ctx, req.(*GetCompanyLastFewDaysCostRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BillingService_ListBillingByIds_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListBillingByIdsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BillingServiceServer).ListBillingByIds(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BillingService_ListBillingByIds_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BillingServiceServer).ListBillingByIds(ctx, req.(*ListBillingByIdsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// BillingService_ServiceDesc is the grpc.ServiceDesc for BillingService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var BillingService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "account_bill.BillingService",
	HandlerType: (*BillingServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "BillingCreate",
			Handler:    _BillingService_BillingCreate_Handler,
		},
		{
			MethodName: "BillingCreateV2",
			Handler:    _BillingService_BillingCreateV2_Handler,
		},
		{
			MethodName: "BillingGet",
			Handler:    _BillingService_BillingGet_Handler,
		},
		{
			MethodName: "BillingMessage",
			Handler:    _BillingService_BillingMessage_Handler,
		},
		{
			MethodName: "BillingOverview",
			Handler:    _BillingService_BillingOverview_Handler,
		},
		{
			MethodName: "BillingQuery",
			Handler:    _BillingService_BillingQuery_Handler,
		},
		{
			MethodName: "PlatformBillingOverview",
			Handler:    _BillingService_PlatformBillingOverview_Handler,
		},
		{
			MethodName: "PlatformBillingQuery",
			Handler:    _BillingService_PlatformBillingQuery_Handler,
		},
		{
			MethodName: "PlatformBillingStat",
			Handler:    _BillingService_PlatformBillingStat_Handler,
		},
		{
			MethodName: "GetBillingByOutBizs",
			Handler:    _BillingService_GetBillingByOutBizs_Handler,
		},
		{
			MethodName: "UserCost",
			Handler:    _BillingService_UserCost_Handler,
		},
		{
			MethodName: "GetProjectCost",
			Handler:    _BillingService_GetProjectCost_Handler,
		},
		{
			MethodName: "GetCompanyLastFewDaysCost",
			Handler:    _BillingService_GetCompanyLastFewDaysCost_Handler,
		},
		{
			MethodName: "ListBillingByIds",
			Handler:    _BillingService_ListBillingByIds_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/account_bill/billing.proto",
}
