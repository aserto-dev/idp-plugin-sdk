package plugin_test

import (
	"context"
	"errors"
	"testing"

	proto "github.com/aserto-dev/go-grpc/aserto/idpplugin/v1"
	"github.com/aserto-dev/idp-plugin-sdk/mocks"
	"github.com/aserto-dev/idp-plugin-sdk/plugin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestValidateFailParseConfig(t *testing.T) {
	// Arrange
	assert := require.New(t)
	handler := mocks.NewMockHandler(gomock.NewController(t))
	pluginServer := &plugin.AsertoPluginServer{Handler: handler}
	req := &proto.ValidateRequest{}
	pluginServer.Handler.(*mocks.MockHandler).EXPECT().GetConfig().Return(nil)

	// Act
	_, err := pluginServer.Validate(context.Background(), req)

	// Assert
	assert.EqualError(err, "rpc error: code = InvalidArgument desc = failed to parse config")
}

func TestValidateUnknownOperationType(t *testing.T) {
	// Arrange
	assert := require.New(t)
	ctrl := gomock.NewController(t)
	handler := mocks.NewMockHandler(ctrl)
	pluginServer := &plugin.AsertoPluginServer{Handler: handler}
	pluginConfig := mocks.NewMockConfig(ctrl)
	req := &proto.ValidateRequest{OpType: proto.OperationType_OPERATION_TYPE_UNKNOWN}
	pluginServer.Handler.(*mocks.MockHandler).EXPECT().GetConfig().Return(pluginConfig)

	// Act
	_, err := pluginServer.Validate(context.Background(), req)

	// Assert
	assert.EqualError(err, "rpc error: code = InvalidArgument desc = unknown operation type provided")
}

func TestValidate(t *testing.T) {
	// Arrange
	assert := require.New(t)
	ctrl := gomock.NewController(t)
	handler := mocks.NewMockHandler(ctrl)
	pluginServer := &plugin.AsertoPluginServer{Handler: handler}
	pluginConfig := mocks.NewMockConfig(ctrl)
	req := &proto.ValidateRequest{OpType: proto.OperationType_OPERATION_TYPE_IMPORT}
	pluginServer.Handler.(*mocks.MockHandler).EXPECT().GetConfig().Return(pluginConfig)
	pluginConfig.EXPECT().Validate(plugin.OperationTypeWrite).Return(nil)

	// Act
	_, err := pluginServer.Validate(context.Background(), req)

	// Assert
	assert.NoError(err)
}

func TestValidateFail(t *testing.T) {
	// Arrange
	assert := require.New(t)
	ctrl := gomock.NewController(t)
	handler := mocks.NewMockHandler(ctrl)
	pluginServer := &plugin.AsertoPluginServer{Handler: handler}
	pluginConfig := mocks.NewMockConfig(ctrl)
	req := &proto.ValidateRequest{OpType: proto.OperationType_OPERATION_TYPE_IMPORT}
	pluginServer.Handler.(*mocks.MockHandler).EXPECT().GetConfig().Return(pluginConfig)
	pluginConfig.EXPECT().Validate(plugin.OperationTypeWrite).Return(errors.New("#boom#"))

	// Act
	_, err := pluginServer.Validate(context.Background(), req)

	// Assert
	assert.EqualError(err, "#boom#")
}
