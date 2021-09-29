package plugin

import (
	"context"

	"github.com/aserto-dev/aserto-idp/pkg/proto"
	api "github.com/aserto-dev/go-grpc/aserto/api/v1"
	"github.com/aserto-dev/idp-plugin-sdk/config"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PluginHandler interface {
	GetConfig() PluginConfig
	Open(config PluginConfig) error
	Read() ([]*api.User, error)
	Write(*api.User) error
	Close() error
}

type PluginConfig interface {
	Validate() error
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

	return response, cfg.Validate()
}
