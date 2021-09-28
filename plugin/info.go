package plugin

import (
	"context"

	"github.com/aserto-dev/aserto-idp/pkg/proto"
	"github.com/aserto-dev/idp-plugin-sdk/config"
	"github.com/aserto-dev/idp-plugin-sdk/version"
)

func (s AsertoPluginServer) Info(ctx context.Context, req *proto.InfoRequest) (*proto.InfoResponse, error) {
	cfg := s.PluginHandler.GetConfig()

	conf, err := config.ParseApiConfig(cfg)
	if err != nil {
		return nil, err
	}

	response := proto.InfoResponse{
		Build:       version.GetBuildInfo(config.GetVersion),
		Description: cfg.Description(),
		Configs:     conf,
	}

	return &response, nil
}
