// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.25.1
// source: inventory.proto

package rpc

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

// InventoryClient is the client API for Inventory service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type InventoryClient interface {
	AddItem(ctx context.Context, in *InventoryRequest, opts ...grpc.CallOption) (*SimpleResponse, error)
	GetItems(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*ItemsResponse, error)
	AddQuantity(ctx context.Context, in *AddQuantityRequest, opts ...grpc.CallOption) (*SimpleResponse, error)
	Search(ctx context.Context, in *SearchRequest, opts ...grpc.CallOption) (Inventory_SearchClient, error)
	AddItems(ctx context.Context, opts ...grpc.CallOption) (Inventory_AddItemsClient, error)
}

type inventoryClient struct {
	cc grpc.ClientConnInterface
}

func NewInventoryClient(cc grpc.ClientConnInterface) InventoryClient {
	return &inventoryClient{cc}
}

func (c *inventoryClient) AddItem(ctx context.Context, in *InventoryRequest, opts ...grpc.CallOption) (*SimpleResponse, error) {
	out := new(SimpleResponse)
	err := c.cc.Invoke(ctx, "/inventory.Inventory/AddItem", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *inventoryClient) GetItems(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*ItemsResponse, error) {
	out := new(ItemsResponse)
	err := c.cc.Invoke(ctx, "/inventory.Inventory/GetItems", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *inventoryClient) AddQuantity(ctx context.Context, in *AddQuantityRequest, opts ...grpc.CallOption) (*SimpleResponse, error) {
	out := new(SimpleResponse)
	err := c.cc.Invoke(ctx, "/inventory.Inventory/AddQuantity", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *inventoryClient) Search(ctx context.Context, in *SearchRequest, opts ...grpc.CallOption) (Inventory_SearchClient, error) {
	stream, err := c.cc.NewStream(ctx, &Inventory_ServiceDesc.Streams[0], "/inventory.Inventory/Search", opts...)
	if err != nil {
		return nil, err
	}
	x := &inventorySearchClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Inventory_SearchClient interface {
	Recv() (*Item, error)
	grpc.ClientStream
}

type inventorySearchClient struct {
	grpc.ClientStream
}

func (x *inventorySearchClient) Recv() (*Item, error) {
	m := new(Item)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *inventoryClient) AddItems(ctx context.Context, opts ...grpc.CallOption) (Inventory_AddItemsClient, error) {
	stream, err := c.cc.NewStream(ctx, &Inventory_ServiceDesc.Streams[1], "/inventory.Inventory/AddItems", opts...)
	if err != nil {
		return nil, err
	}
	x := &inventoryAddItemsClient{stream}
	return x, nil
}

type Inventory_AddItemsClient interface {
	Send(*InventoryRequest) error
	CloseAndRecv() (*TotalItemsResponse, error)
	grpc.ClientStream
}

type inventoryAddItemsClient struct {
	grpc.ClientStream
}

func (x *inventoryAddItemsClient) Send(m *InventoryRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *inventoryAddItemsClient) CloseAndRecv() (*TotalItemsResponse, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(TotalItemsResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// InventoryServer is the server API for Inventory service.
// All implementations must embed UnimplementedInventoryServer
// for forward compatibility
type InventoryServer interface {
	AddItem(context.Context, *InventoryRequest) (*SimpleResponse, error)
	GetItems(context.Context, *Empty) (*ItemsResponse, error)
	AddQuantity(context.Context, *AddQuantityRequest) (*SimpleResponse, error)
	Search(*SearchRequest, Inventory_SearchServer) error
	AddItems(Inventory_AddItemsServer) error
	mustEmbedUnimplementedInventoryServer()
}

// UnimplementedInventoryServer must be embedded to have forward compatible implementations.
type UnimplementedInventoryServer struct {
}

func (UnimplementedInventoryServer) AddItem(context.Context, *InventoryRequest) (*SimpleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddItem not implemented")
}
func (UnimplementedInventoryServer) GetItems(context.Context, *Empty) (*ItemsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetItems not implemented")
}
func (UnimplementedInventoryServer) AddQuantity(context.Context, *AddQuantityRequest) (*SimpleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddQuantity not implemented")
}
func (UnimplementedInventoryServer) Search(*SearchRequest, Inventory_SearchServer) error {
	return status.Errorf(codes.Unimplemented, "method Search not implemented")
}
func (UnimplementedInventoryServer) AddItems(Inventory_AddItemsServer) error {
	return status.Errorf(codes.Unimplemented, "method AddItems not implemented")
}
func (UnimplementedInventoryServer) mustEmbedUnimplementedInventoryServer() {}

// UnsafeInventoryServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to InventoryServer will
// result in compilation errors.
type UnsafeInventoryServer interface {
	mustEmbedUnimplementedInventoryServer()
}

func RegisterInventoryServer(s grpc.ServiceRegistrar, srv InventoryServer) {
	s.RegisterService(&Inventory_ServiceDesc, srv)
}

func _Inventory_AddItem_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InventoryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(InventoryServer).AddItem(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/inventory.Inventory/AddItem",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(InventoryServer).AddItem(ctx, req.(*InventoryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Inventory_GetItems_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(InventoryServer).GetItems(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/inventory.Inventory/GetItems",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(InventoryServer).GetItems(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Inventory_AddQuantity_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddQuantityRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(InventoryServer).AddQuantity(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/inventory.Inventory/AddQuantity",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(InventoryServer).AddQuantity(ctx, req.(*AddQuantityRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Inventory_Search_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(SearchRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(InventoryServer).Search(m, &inventorySearchServer{stream})
}

type Inventory_SearchServer interface {
	Send(*Item) error
	grpc.ServerStream
}

type inventorySearchServer struct {
	grpc.ServerStream
}

func (x *inventorySearchServer) Send(m *Item) error {
	return x.ServerStream.SendMsg(m)
}

func _Inventory_AddItems_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(InventoryServer).AddItems(&inventoryAddItemsServer{stream})
}

type Inventory_AddItemsServer interface {
	SendAndClose(*TotalItemsResponse) error
	Recv() (*InventoryRequest, error)
	grpc.ServerStream
}

type inventoryAddItemsServer struct {
	grpc.ServerStream
}

func (x *inventoryAddItemsServer) SendAndClose(m *TotalItemsResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *inventoryAddItemsServer) Recv() (*InventoryRequest, error) {
	m := new(InventoryRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Inventory_ServiceDesc is the grpc.ServiceDesc for Inventory service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Inventory_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "inventory.Inventory",
	HandlerType: (*InventoryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AddItem",
			Handler:    _Inventory_AddItem_Handler,
		},
		{
			MethodName: "GetItems",
			Handler:    _Inventory_GetItems_Handler,
		},
		{
			MethodName: "AddQuantity",
			Handler:    _Inventory_AddQuantity_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Search",
			Handler:       _Inventory_Search_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "AddItems",
			Handler:       _Inventory_AddItems_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "inventory.proto",
}
