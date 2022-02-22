package plugin

import (
	"context"
	"fmt"
	"log"

	api "github.com/aserto-dev/go-grpc/aserto/api/v1"
	proto "github.com/aserto-dev/go-grpc/aserto/idpplugin/v1"
	"github.com/aserto-dev/idp-plugin-sdk/config"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OperationType int

const (
	OperationTypeWrite OperationType = iota
	OperationTypeRead
	OperationTypeDelete
)

type Handler interface {
	GetConfig() Config
	Open(Config, OperationType) error
	Read() ([]*api.User, error)
	Write(*api.User) error
	Delete(string) error
	Close() (*Stats, error)
	GetVersion() (string, string, string)
}

type Config interface {
	Validate(OperationType) error
	Description() string
}

type AsertoPluginServer struct {
	Handler Handler
}

func (s AsertoPluginServer) cleanup(pluginClosed bool, methodName string) {
	if !pluginClosed {
		_, err := s.Handler.Close()
		if err != nil {
			log.Println(err.Error())
		}
	}

	if r := recover(); r != nil {
		log.Println(fmt.Errorf("recovering from panic in %s error is: %v", methodName, r))
	}
}

func (s AsertoPluginServer) Validate(ctx context.Context, req *proto.ValidateRequest) (*proto.ValidateResponse, error) {
	response := &proto.ValidateResponse{}

	cfg := s.Handler.GetConfig()
	err := config.NewConfig(req.Config, cfg)
	if err != nil {
		return response, status.Error(codes.InvalidArgument, "failed to parse config")
	}
	opType := req.OpType
	if opType == proto.OperationType_OPERATION_TYPE_UNKNOWN {
		return response, status.Error(codes.InvalidArgument, "unknown operation type provided")
	}
	opTypePlugin := OperationType(opType - 1)

	return response, cfg.Validate(opTypePlugin)
}
