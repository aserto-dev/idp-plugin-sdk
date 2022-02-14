package plugin

import (
	"context"

	api "github.com/aserto-dev/go-grpc/aserto/api/v1"
	proto "github.com/aserto-dev/go-grpc/aserto/idpplugin/v1"
	"github.com/aserto-dev/idp-plugin-sdk/config"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//go:generate mockgen -source=plugin.go -destination=mock_plugin.go -package=plugin --build_flags=--mod=mod
//go:generate mockgen -destination=mock_servers.go -package=plugin github.com/aserto-dev/go-grpc/aserto/idpplugin/v1 Plugin_DeleteServer,Plugin_ExportServer,Plugin_ImportServer

type OperationType int

const (
	OperationTypeWrite OperationType = iota
	OperationTypeRead
	OperationTypeDelete
)

type PluginHandler interface {
	GetConfig() PluginConfig
	Open(PluginConfig, OperationType) error
	Read() ([]*api.User, error)
	Write(*api.User) error
	Delete(string) error
	Close() (*Stats, error)
	GetVersion() (string, string, string)
}

type PluginConfig interface {
	Validate(OperationType) error
	Description() string
}

type AsertoPluginServer struct {
	PluginHandler PluginHandler
}

func (s AsertoPluginServer) Validate(ctx context.Context, req *proto.ValidateRequest) (*proto.ValidateResponse, error) {
	response := &proto.ValidateResponse{}

	cfg := s.PluginHandler.GetConfig()
	err := config.NewConfig(req.Config, cfg)
	if err != nil {
		return response, status.Error(codes.InvalidArgument, "failed to parse config")
	}
	opType := req.OpType
	if opType == proto.OperationType_OPERATION_TYPE_UNKNOWN {
		return response, status.Error(codes.InvalidArgument, "unknown operation type provided")
	}
	opTypePlugin := (OperationType)(opType - 1)

	return response, cfg.Validate(opTypePlugin)
}
