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
	status "google.golang.org/genproto/googleapis/rpc/status"
)

func TestImportNoUsers(t *testing.T) {
	// Arrange
	assert := require.New(t)
	ctrl := gomock.NewController(t)
	handler := mocks.NewMockPluginHandler(ctrl)
	pluginServer := &plugin.AsertoPluginServer{PluginHandler: handler}
	pluginConfig := mocks.NewMockPluginConfig(ctrl)
	importServer := mocks.NewMockPlugin_ImportServer(ctrl)

	pluginServer.PluginHandler.(*mocks.MockPluginHandler).EXPECT().GetConfig().Return(pluginConfig)
	importServer.EXPECT().Recv().Return(nil, io.EOF)
	pluginServer.PluginHandler.(*mocks.MockPluginHandler).EXPECT().Close().Return(nil, nil)

	//Act
	err := pluginServer.Import(importServer)

	//Assert
	assert.NoError(err)
}

func TestImportNoUsersFailClose(t *testing.T) {
	// Arrange
	assert := require.New(t)
	ctrl := gomock.NewController(t)
	handler := mocks.NewMockPluginHandler(ctrl)
	pluginServer := &plugin.AsertoPluginServer{PluginHandler: handler}
	pluginConfig := mocks.NewMockPluginConfig(ctrl)
	importServer := mocks.NewMockPlugin_ImportServer(ctrl)
	boomErr := errors.New("Boom!")
	importResp := &proto.ImportResponse{Error: &status.Status{Message: boomErr.Error()}}

	pluginServer.PluginHandler.(*mocks.MockPluginHandler).EXPECT().GetConfig().Return(pluginConfig)
	importServer.EXPECT().Recv().Return(nil, io.EOF)
	pluginServer.PluginHandler.(*mocks.MockPluginHandler).EXPECT().Close().Return(nil, boomErr)
	importServer.EXPECT().Send(importResp).Return(nil)

	//Act
	err := pluginServer.Import(importServer)

	//Assert
	assert.NoError(err)
}

func TestImportReceiveError(t *testing.T) {
	// Arrange
	assert := require.New(t)
	ctrl := gomock.NewController(t)
	handler := mocks.NewMockPluginHandler(ctrl)
	pluginServer := &plugin.AsertoPluginServer{PluginHandler: handler}
	pluginConfig := mocks.NewMockPluginConfig(ctrl)
	importServer := mocks.NewMockPlugin_ImportServer(ctrl)
	boomErr := errors.New("Boom!")
	importResp := &proto.ImportResponse{Error: &status.Status{Message: boomErr.Error()}}

	pluginServer.PluginHandler.(*mocks.MockPluginHandler).EXPECT().GetConfig().Return(pluginConfig)
	importServer.EXPECT().Recv().Return(nil, boomErr)
	pluginServer.PluginHandler.(*mocks.MockPluginHandler).EXPECT().Close().Return(nil, nil)
	importServer.EXPECT().Send(importResp).Return(nil)

	//Act
	err := pluginServer.Import(importServer)

	//Assert
	assert.NoError(err)
}

func TestImportWithNilUser(t *testing.T) {
	// Arrange
	assert := require.New(t)
	ctrl := gomock.NewController(t)
	handler := mocks.NewMockPluginHandler(ctrl)
	pluginServer := &plugin.AsertoPluginServer{PluginHandler: handler}
	pluginConfig := mocks.NewMockPluginConfig(ctrl)
	importServer := mocks.NewMockPlugin_ImportServer(ctrl)
	importResp := &proto.ImportResponse{Error: nil}

	pluginServer.PluginHandler.(*mocks.MockPluginHandler).EXPECT().GetConfig().Return(pluginConfig)
	importServer.EXPECT().Recv().Return(nil, nil)
	pluginServer.PluginHandler.(*mocks.MockPluginHandler).EXPECT().Close().Return(nil, nil)
	importServer.EXPECT().Send(importResp).Return(nil).AnyTimes()
	pluginServer.PluginHandler.(*mocks.MockPluginHandler).EXPECT().Open(pluginConfig, plugin.OperationTypeWrite).Return(nil).Times(1)
	importServer.EXPECT().Recv().Return(nil, io.EOF)

	//Act
	err := pluginServer.Import(importServer)

	//Assert
	assert.NoError(err)
}

func TestImportWithUser(t *testing.T) {
	// Arrange
	assert := require.New(t)
	ctrl := gomock.NewController(t)
	handler := mocks.NewMockPluginHandler(ctrl)
	pluginServer := &plugin.AsertoPluginServer{PluginHandler: handler}
	pluginConfig := mocks.NewMockPluginConfig(ctrl)
	importServer := mocks.NewMockPlugin_ImportServer(ctrl)
	importResp := &proto.ImportResponse{Error: nil}
	user := &api.User{Id: "testID"}
	importReq := &proto.ImportRequest{Data: &proto.ImportRequest_User{User: user}}
	stats := &plugin.Stats{Received: 1, Created: 1}

	pluginServer.PluginHandler.(*mocks.MockPluginHandler).EXPECT().GetConfig().Return(pluginConfig)
	importServer.EXPECT().Recv().Return(importReq, nil)
	importServer.EXPECT().Send(importResp).Return(nil).AnyTimes()
	pluginServer.PluginHandler.(*mocks.MockPluginHandler).EXPECT().Open(pluginConfig, plugin.OperationTypeWrite).Return(nil).Times(1)
	pluginServer.PluginHandler.(*mocks.MockPluginHandler).EXPECT().Write(user).Return(nil).Times(1)
	importServer.EXPECT().Recv().Return(nil, io.EOF)
	pluginServer.PluginHandler.(*mocks.MockPluginHandler).EXPECT().Close().Return(stats, nil)
	importServer.EXPECT().Send(gomock.Any()).Return(nil).AnyTimes()

	//Act
	err := pluginServer.Import(importServer)

	//Assert
	assert.NoError(err)
}
