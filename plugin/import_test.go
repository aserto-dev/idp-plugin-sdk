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
	handler := mocks.NewMockHandler(ctrl)
	pluginServer := &plugin.AsertoPluginServer{Handler: handler}
	pluginConfig := mocks.NewMockConfig(ctrl)
	importServer := mocks.NewMockPlugin_ImportServer(ctrl)

	pluginServer.Handler.(*mocks.MockHandler).EXPECT().GetConfig().Return(pluginConfig)
	importServer.EXPECT().Recv().Return(nil, io.EOF)
	pluginServer.Handler.(*mocks.MockHandler).EXPECT().Close().Return(nil, nil)

	// Act
	err := pluginServer.Import(importServer)

	// Assert
	assert.NoError(err)
}

func TestImportNoUsersFailClose(t *testing.T) {
	// Arrange
	assert := require.New(t)
	ctrl := gomock.NewController(t)
	handler := mocks.NewMockHandler(ctrl)
	pluginServer := &plugin.AsertoPluginServer{Handler: handler}
	pluginConfig := mocks.NewMockConfig(ctrl)
	importServer := mocks.NewMockPlugin_ImportServer(ctrl)
	boomErr := errors.New("#boom#")
	importResp := &proto.ImportResponse{Error: &status.Status{Message: boomErr.Error()}}

	pluginServer.Handler.(*mocks.MockHandler).EXPECT().GetConfig().Return(pluginConfig)
	importServer.EXPECT().Recv().Return(nil, io.EOF)
	pluginServer.Handler.(*mocks.MockHandler).EXPECT().Close().Return(nil, boomErr)
	importServer.EXPECT().Send(importResp).Return(nil)

	// Act
	err := pluginServer.Import(importServer)

	// Assert
	assert.NoError(err)
}

func TestImportReceiveError(t *testing.T) {
	// Arrange
	assert := require.New(t)
	ctrl := gomock.NewController(t)
	handler := mocks.NewMockHandler(ctrl)
	pluginServer := &plugin.AsertoPluginServer{Handler: handler}
	pluginConfig := mocks.NewMockConfig(ctrl)
	importServer := mocks.NewMockPlugin_ImportServer(ctrl)
	boomErr := errors.New("#boom#")
	importResp := &proto.ImportResponse{Error: &status.Status{Message: boomErr.Error()}}

	pluginServer.Handler.(*mocks.MockHandler).EXPECT().GetConfig().Return(pluginConfig)
	importServer.EXPECT().Recv().Return(nil, boomErr)
	pluginServer.Handler.(*mocks.MockHandler).EXPECT().Close().Return(nil, nil)
	importServer.EXPECT().Send(importResp).Return(nil)

	// Act
	err := pluginServer.Import(importServer)

	// Assert
	assert.NoError(err)
}

func TestImportWithNilUser(t *testing.T) {
	// Arrange
	assert := require.New(t)
	ctrl := gomock.NewController(t)
	handler := mocks.NewMockHandler(ctrl)
	pluginServer := &plugin.AsertoPluginServer{Handler: handler}
	pluginConfig := mocks.NewMockConfig(ctrl)
	importServer := mocks.NewMockPlugin_ImportServer(ctrl)
	importResp := &proto.ImportResponse{Error: nil}

	pluginServer.Handler.(*mocks.MockHandler).EXPECT().GetConfig().Return(pluginConfig)
	importServer.EXPECT().Recv().Return(nil, nil)
	pluginServer.Handler.(*mocks.MockHandler).EXPECT().Close().Return(nil, nil)
	importServer.EXPECT().Send(importResp).Return(nil).AnyTimes()
	pluginServer.Handler.(*mocks.MockHandler).EXPECT().Open(pluginConfig, plugin.OperationTypeWrite).Return(nil).Times(1)
	importServer.EXPECT().Recv().Return(nil, io.EOF)

	// Act
	err := pluginServer.Import(importServer)

	// Assert
	assert.NoError(err)
}

func TestImportWithUser(t *testing.T) {
	// Arrange
	assert := require.New(t)
	ctrl := gomock.NewController(t)
	handler := mocks.NewMockHandler(ctrl)
	pluginServer := &plugin.AsertoPluginServer{Handler: handler}
	pluginConfig := mocks.NewMockConfig(ctrl)
	importServer := mocks.NewMockPlugin_ImportServer(ctrl)
	importResp := &proto.ImportResponse{Error: nil}
	user := &api.User{Id: "testID"}
	importReq := &proto.ImportRequest{Data: &proto.ImportRequest_User{User: user}}
	stats := &plugin.Stats{Received: 1, Created: 1}

	pluginServer.Handler.(*mocks.MockHandler).EXPECT().GetConfig().Return(pluginConfig)
	importServer.EXPECT().Recv().Return(importReq, nil)
	importServer.EXPECT().Send(importResp).Return(nil).AnyTimes()
	pluginServer.Handler.(*mocks.MockHandler).EXPECT().Open(pluginConfig, plugin.OperationTypeWrite).Return(nil).Times(1)
	pluginServer.Handler.(*mocks.MockHandler).EXPECT().Write(user).Return(nil).Times(1)
	importServer.EXPECT().Recv().Return(nil, io.EOF)
	pluginServer.Handler.(*mocks.MockHandler).EXPECT().Close().Return(stats, nil)
	importServer.EXPECT().Send(gomock.Any()).Return(nil).AnyTimes()

	// Act
	err := pluginServer.Import(importServer)

	// Assert
	assert.NoError(err)
}
