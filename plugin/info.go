package plugin

import (
	"context"

	proto "github.com/aserto-dev/go-grpc/aserto/idpplugin/v1"
	"github.com/aserto-dev/idp-plugin-sdk/config"
	"github.com/aserto-dev/idp-plugin-sdk/version"
)

func (s AsertoPluginServer) Info(ctx context.Context, req *proto.InfoRequest) (*proto.InfoResponse, error) {
	cfg := s.Handler.GetConfig()

	conf, err := config.ParseAPIConfig(cfg)
	if err != nil {
		return nil, err
	}

	response := proto.InfoResponse{
		Build:       version.GetBuildInfo(s.Handler.GetVersion),
		Description: cfg.Description(),
		Configs:     conf,
	}

	return &response, nil
}
