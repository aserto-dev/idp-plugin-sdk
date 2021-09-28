package plugin

import (
	"context"

	"github.com/aserto-dev/aserto-idp/pkg/proto"
	api "github.com/aserto-dev/go-grpc/aserto/api/v1"
)

type PluginHandler interface {
	// Info(ctx context.Context) ([]*api.ConfigElement, error)
	// Import(ctx context.Context, errChannel <-chan error, usersChannel <-chan *api.User) error
	// Export(ctx context.Context, errChannel <-chan error, usersChannel chan<- *api.User) error
	// Export(ctx context.Context, done func()) (func() (*api.User, error), error)
	// Validate(ctx context.Context, config map[string]interface{}) error
	// Import(config interface{}, recv func() *api.User, errf func(error)) error
	// Export(config interface{}, send func(*api.User), errf func(error)) error
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

	return response, s.PluginHandler.GetConfig().Validate()
}
