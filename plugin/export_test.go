package plugin_test

import (
	"errors"
	"io"
	"testing"

	"github.com/aserto-dev/go-grpc/aserto/api/v1"
	proto "github.com/aserto-dev/go-grpc/aserto/idpplugin/v1"
	"github.com/aserto-dev/idp-plugin-sdk/mocks"
	"github.com/aserto-dev/idp-plugin-sdk/plugin"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestExportNoConfig(t *testing.T) {
	// Arrange
	assert := require.New(t)
	ctrl := gomock.NewController(t)
	handler := mocks.NewMockPluginHandler(ctrl)
	pluginServer := &plugin.AsertoPluginServer{PluginHandler: handler}
	exportServer := mocks.NewMockPlugin_ExportServer(ctrl)
	exportReq := &proto.ExportRequest{Config: nil}

	pluginServer.PluginHandler.(*mocks.MockPluginHandler).EXPECT().GetConfig().Return(nil)
	pluginServer.PluginHandler.(*mocks.MockPluginHandler).EXPECT().Close().Return(nil, nil)

	// Act
	err := pluginServer.Export(exportReq, exportServer)

	// Assert
	assert.NotNil(err)
}

func TestExportOpenErrors(t *testing.T) {
	// Arrange
	assert := require.New(t)
	ctrl := gomock.NewController(t)
	handler := mocks.NewMockPluginHandler(ctrl)
	pluginServer := &plugin.AsertoPluginServer{PluginHandler: handler}
	pluginConfig := mocks.NewMockPluginConfig(ctrl)
	exportServer := mocks.NewMockPlugin_ExportServer(ctrl)
	boomErr := errors.New("#boom#")
	exportReq := &proto.ExportRequest{Config: nil}

	pluginServer.PluginHandler.(*mocks.MockPluginHandler).EXPECT().GetConfig().Return(pluginConfig)
	pluginServer.PluginHandler.(*mocks.MockPluginHandler).EXPECT().Open(gomock.Any(), plugin.OperationTypeRead).Return(boomErr)
	pluginServer.PluginHandler.(*mocks.MockPluginHandler).EXPECT().Close().Return(nil, nil)

	// Act
	err := pluginServer.Export(exportReq, exportServer)

	// Assert
	assert.NotNil(err)
}

func TestExportNoUsers(t *testing.T) {
	// Arrange
	assert := require.New(t)
	ctrl := gomock.NewController(t)
	handler := mocks.NewMockPluginHandler(ctrl)
	pluginServer := &plugin.AsertoPluginServer{PluginHandler: handler}
	pluginConfig := mocks.NewMockPluginConfig(ctrl)
	exportServer := mocks.NewMockPlugin_ExportServer(ctrl)
	exportReq := &proto.ExportRequest{Config: nil}

	pluginServer.PluginHandler.(*mocks.MockPluginHandler).EXPECT().GetConfig().Return(pluginConfig)
	pluginServer.PluginHandler.(*mocks.MockPluginHandler).EXPECT().Open(gomock.Any(), plugin.OperationTypeRead).Return(nil)
	pluginServer.PluginHandler.(*mocks.MockPluginHandler).EXPECT().Read().Return(nil, io.EOF)
	pluginServer.PluginHandler.(*mocks.MockPluginHandler).EXPECT().Close().Return(nil, nil)

	// Act
	err := pluginServer.Export(exportReq, exportServer)

	// Assert
	assert.NoError(err)
}

func TestExportWhenReadFails(t *testing.T) {
	// Arrange
	assert := require.New(t)
	ctrl := gomock.NewController(t)
	handler := mocks.NewMockPluginHandler(ctrl)
	pluginServer := &plugin.AsertoPluginServer{PluginHandler: handler}
	pluginConfig := mocks.NewMockPluginConfig(ctrl)
	exportServer := mocks.NewMockPlugin_ExportServer(ctrl)
	boomErr := errors.New("#boom#")
	exportReq := &proto.ExportRequest{Config: nil}

	pluginServer.PluginHandler.(*mocks.MockPluginHandler).EXPECT().GetConfig().Return(pluginConfig)
	pluginServer.PluginHandler.(*mocks.MockPluginHandler).EXPECT().Open(gomock.Any(), plugin.OperationTypeRead).Return(nil)
	pluginServer.PluginHandler.(*mocks.MockPluginHandler).EXPECT().Read().Return(nil, boomErr)
	exportServer.EXPECT().Send(gomock.Any()).Return(boomErr)
	pluginServer.PluginHandler.(*mocks.MockPluginHandler).EXPECT().Close().Return(nil, nil)

	// Act
	err := pluginServer.Export(exportReq, exportServer)

	// Assert
	assert.NoError(err)
}

func TestExportOneUser(t *testing.T) {
	// Arrange
	assert := require.New(t)
	ctrl := gomock.NewController(t)
	handler := mocks.NewMockPluginHandler(ctrl)
	pluginServer := &plugin.AsertoPluginServer{PluginHandler: handler}
	pluginConfig := mocks.NewMockPluginConfig(ctrl)
	exportServer := mocks.NewMockPlugin_ExportServer(ctrl)
	exportReq := &proto.ExportRequest{Config: nil}
	user := &api.User{Id: "testID"}
	var users []*api.User
	users = append(users, user)
	expResp := &proto.ExportResponse{
		Data: &proto.ExportResponse_User{
			User: user,
		},
	}

	pluginServer.PluginHandler.(*mocks.MockPluginHandler).EXPECT().GetConfig().Return(pluginConfig)
	pluginServer.PluginHandler.(*mocks.MockPluginHandler).EXPECT().Open(gomock.Any(), plugin.OperationTypeRead).Return(nil)
	pluginServer.PluginHandler.(*mocks.MockPluginHandler).EXPECT().Read().Return(users, nil)
	exportServer.EXPECT().Send(expResp).Return(nil)
	pluginServer.PluginHandler.(*mocks.MockPluginHandler).EXPECT().Read().Return(nil, io.EOF)
	pluginServer.PluginHandler.(*mocks.MockPluginHandler).EXPECT().Close().Return(nil, nil)

	// Act
	err := pluginServer.Export(exportReq, exportServer)

	// Assert
	assert.NoError(err)
}
