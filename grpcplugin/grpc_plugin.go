package grpcplugin

import (
	"context"

	"github.com/aserto-dev/aserto-idp/pkg/proto"
	plugin "github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
)

// PluginServer is the server API for the Resource service.
type PluginServer interface {
	proto.PluginServer
}

// PluginClient is the client API for the Resource service.
type PluginClient interface {
	proto.PluginClient
}

// This is the implementation of plugin.GRPCPlugin so we can serve/consume this.
type PluginGRPC struct {
	// GRPCPlugin must still implement the Plugin interface
	plugin.NetRPCUnsupportedPlugin
	plugin.GRPCPlugin
	// Concrete implementation, written in Go. This is only used for plugins
	// that are written in Go.
	Impl PluginServer
}

func (p *PluginGRPC) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	proto.RegisterPluginServer(s, &pluginGRPCServer{
		server: p.Impl,
	})
	return nil
}

func (p *PluginGRPC) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &pluginGRPCClient{client: proto.NewPluginClient(c)}, nil
}

type pluginGRPCServer struct {
	server proto.PluginServer
}

func (s *pluginGRPCServer) Info(ctx context.Context, req *proto.InfoRequest) (*proto.InfoResponse, error) {
	return s.server.Info(ctx, req)
}

func (s *pluginGRPCServer) Import(srv proto.Plugin_ImportServer) error {
	return s.server.Import(srv)
}

func (s *pluginGRPCServer) Export(req *proto.ExportRequest, srv proto.Plugin_ExportServer) error {
	return s.server.Export(req, srv)
}

// func (s *pluginGRPCServer) Delete(srv proto.Plugin_DeleteServer) error {
// 	return s.Impl.Delete(srv)
// }

func (s *pluginGRPCServer) Validate(ctx context.Context, req *proto.ValidateRequest) (*proto.ValidateResponse, error) {
	return s.server.Validate(ctx, req)
}

type pluginGRPCClient struct {
	client proto.PluginClient
}

func (c *pluginGRPCClient) Info(ctx context.Context, in *proto.InfoRequest, opts ...grpc.CallOption) (*proto.InfoResponse, error) {
	return c.client.Info(ctx, in, opts...)
}

func (c *pluginGRPCClient) Import(ctx context.Context, opts ...grpc.CallOption) (proto.Plugin_ImportClient, error) {
	return c.client.Import(ctx, opts...)
}

func (c *pluginGRPCClient) Export(ctx context.Context, in *proto.ExportRequest, opts ...grpc.CallOption) (proto.Plugin_ExportClient, error) {
	return c.client.Export(ctx, in, opts...)
}

// func (c *pluginGRPCClient) Delete(ctx context.Context, opts ...grpc.CallOption) (proto.Plugin_DeleteClient, error) {
// 	return c.client.Delete(ctx, opts...)
// }

func (c *pluginGRPCClient) Validate(ctx context.Context, in *proto.ValidateRequest, opts ...grpc.CallOption) (*proto.ValidateResponse, error) {
	return c.client.Validate(ctx, in, opts...)
}

var _ PluginServer = &pluginGRPCServer{}
var _ PluginClient = &pluginGRPCClient{}
