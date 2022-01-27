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
	handler := mocks.NewMockPluginHandler(gomock.NewController(t))
	pluginServer := &plugin.AsertoPluginServer{PluginHandler: handler}
	req := &proto.ValidateRequest{}
	pluginServer.PluginHandler.(*mocks.MockPluginHandler).EXPECT().GetConfig().Return(nil)

	//Act
	_, err := pluginServer.Validate(context.Background(), req)

	//Assert
	assert.EqualError(err, "rpc error: code = InvalidArgument desc = failed to parse config")
}

func TestValidateUnknownOperationType(t *testing.T) {
	// Arrange
	assert := require.New(t)
	ctrl := gomock.NewController(t)
	handler := mocks.NewMockPluginHandler(ctrl)
	pluginServer := &plugin.AsertoPluginServer{PluginHandler: handler}
	pluginConfig := mocks.NewMockPluginConfig(ctrl)
	req := &proto.ValidateRequest{OpType: proto.OperationType_OPERATION_TYPE_UNKNOWN}
	pluginServer.PluginHandler.(*mocks.MockPluginHandler).EXPECT().GetConfig().Return(pluginConfig)

	//Act
	_, err := pluginServer.Validate(context.Background(), req)

	//Assert
	assert.EqualError(err, "rpc error: code = InvalidArgument desc = unknown operation type provided")
}

func TestValidate(t *testing.T) {
	// Arrange
	assert := require.New(t)
	ctrl := gomock.NewController(t)
	handler := mocks.NewMockPluginHandler(ctrl)
	pluginServer := &plugin.AsertoPluginServer{PluginHandler: handler}
	pluginConfig := mocks.NewMockPluginConfig(ctrl)
	req := &proto.ValidateRequest{OpType: proto.OperationType_OPERATION_TYPE_IMPORT}
	pluginServer.PluginHandler.(*mocks.MockPluginHandler).EXPECT().GetConfig().Return(pluginConfig)
	pluginConfig.EXPECT().Validate(plugin.OperationTypeWrite).Return(nil)

	//Act
	_, err := pluginServer.Validate(context.Background(), req)

	//Assert
	assert.NoError(err)
}

func TestValidateFail(t *testing.T) {
	// Arrange
	assert := require.New(t)
	ctrl := gomock.NewController(t)
	handler := mocks.NewMockPluginHandler(ctrl)
	pluginServer := &plugin.AsertoPluginServer{PluginHandler: handler}
	pluginConfig := mocks.NewMockPluginConfig(ctrl)
	req := &proto.ValidateRequest{OpType: proto.OperationType_OPERATION_TYPE_IMPORT}
	pluginServer.PluginHandler.(*mocks.MockPluginHandler).EXPECT().GetConfig().Return(pluginConfig)
	pluginConfig.EXPECT().Validate(plugin.OperationTypeWrite).Return(errors.New("Boom!"))

	//Act
	_, err := pluginServer.Validate(context.Background(), req)

	//Assert
	assert.EqualError(err, "Boom!")
}
