package plugin_test

import (
	"errors"
	"io"
	"testing"

	proto "github.com/aserto-dev/go-grpc/aserto/idpplugin/v1"
	"github.com/aserto-dev/idp-plugin-sdk/mocks"
	"github.com/aserto-dev/idp-plugin-sdk/plugin"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	status "google.golang.org/genproto/googleapis/rpc/status"
)

func TestDeleteNoUsers(t *testing.T) {
	// Arrange
	assert := require.New(t)
	ctrl := gomock.NewController(t)
	handler := mocks.NewMockPluginHandler(ctrl)
	pluginServer := &plugin.AsertoPluginServer{PluginHandler: handler}
	pluginConfig := mocks.NewMockPluginConfig(ctrl)
	deleteServer := mocks.NewMockPlugin_DeleteServer(ctrl)

	pluginServer.PluginHandler.(*mocks.MockPluginHandler).EXPECT().GetConfig().Return(pluginConfig)
	deleteServer.EXPECT().Recv().Return(nil, io.EOF)
	pluginServer.PluginHandler.(*mocks.MockPluginHandler).EXPECT().Close().Return(nil, nil)

	// Act
	err := pluginServer.Delete(deleteServer)

	// Assert
	assert.NoError(err)
}

func TestDeleteNoUsersFailClose(t *testing.T) {
	// Arrange
	assert := require.New(t)
	ctrl := gomock.NewController(t)
	handler := mocks.NewMockPluginHandler(ctrl)
	pluginServer := &plugin.AsertoPluginServer{PluginHandler: handler}
	pluginConfig := mocks.NewMockPluginConfig(ctrl)
	deleteServer := mocks.NewMockPlugin_DeleteServer(ctrl)
	boomErr := errors.New("#boom#")
	deleteResp := &proto.DeleteResponse{Error: &status.Status{Message: boomErr.Error()}}

	pluginServer.PluginHandler.(*mocks.MockPluginHandler).EXPECT().GetConfig().Return(pluginConfig)
	deleteServer.EXPECT().Recv().Return(nil, io.EOF)
	pluginServer.PluginHandler.(*mocks.MockPluginHandler).EXPECT().Close().Return(nil, boomErr)
	deleteServer.EXPECT().Send(deleteResp).Return(nil)

	// Act
	err := pluginServer.Delete(deleteServer)

	// Assert
	assert.NoError(err)
}

func TestDeleteReceiveError(t *testing.T) {
	// Arrange
	assert := require.New(t)
	ctrl := gomock.NewController(t)
	handler := mocks.NewMockPluginHandler(ctrl)
	pluginServer := &plugin.AsertoPluginServer{PluginHandler: handler}
	pluginConfig := mocks.NewMockPluginConfig(ctrl)
	deleteServer := mocks.NewMockPlugin_DeleteServer(ctrl)
	boomErr := errors.New("#boom#")
	deleteResp := &proto.DeleteResponse{Error: &status.Status{Message: boomErr.Error()}}

	pluginServer.PluginHandler.(*mocks.MockPluginHandler).EXPECT().GetConfig().Return(pluginConfig)
	deleteServer.EXPECT().Recv().Return(nil, boomErr)
	pluginServer.PluginHandler.(*mocks.MockPluginHandler).EXPECT().Close().Return(nil, nil)
	deleteServer.EXPECT().Send(deleteResp).Return(nil)

	// Act
	err := pluginServer.Delete(deleteServer)

	// Assert
	assert.NoError(err)
}

func TestDeleteWithUserIdEmpty(t *testing.T) {
	// Arrange
	assert := require.New(t)
	ctrl := gomock.NewController(t)
	handler := mocks.NewMockPluginHandler(ctrl)
	pluginServer := &plugin.AsertoPluginServer{PluginHandler: handler}
	pluginConfig := mocks.NewMockPluginConfig(ctrl)
	deleteServer := mocks.NewMockPlugin_DeleteServer(ctrl)
	deleteResp := &proto.DeleteResponse{Error: nil}

	pluginServer.PluginHandler.(*mocks.MockPluginHandler).EXPECT().GetConfig().Return(pluginConfig)
	deleteServer.EXPECT().Recv().Return(nil, nil)
	pluginServer.PluginHandler.(*mocks.MockPluginHandler).EXPECT().Close().Return(nil, nil)
	deleteServer.EXPECT().Send(deleteResp).Return(nil).AnyTimes()
	pluginServer.PluginHandler.(*mocks.MockPluginHandler).EXPECT().Open(pluginConfig, plugin.OperationTypeDelete).Return(nil).Times(1)
	deleteServer.EXPECT().Recv().Return(nil, io.EOF)

	// Act
	err := pluginServer.Delete(deleteServer)

	// Assert
	assert.NoError(err)
}

func TestDeleteWithUserId(t *testing.T) {
	// Arrange
	assert := require.New(t)
	ctrl := gomock.NewController(t)
	handler := mocks.NewMockPluginHandler(ctrl)
	pluginServer := &plugin.AsertoPluginServer{PluginHandler: handler}
	pluginConfig := mocks.NewMockPluginConfig(ctrl)
	deleteServer := mocks.NewMockPlugin_DeleteServer(ctrl)
	deleteResp := &proto.DeleteResponse{Error: nil}
	deleteReq := &proto.DeleteRequest{Data: &proto.DeleteRequest_UserId{UserId: "testUUID"}}

	pluginServer.PluginHandler.(*mocks.MockPluginHandler).EXPECT().GetConfig().Return(pluginConfig)
	deleteServer.EXPECT().Recv().Return(deleteReq, nil)
	pluginServer.PluginHandler.(*mocks.MockPluginHandler).EXPECT().Close().Return(nil, nil)
	deleteServer.EXPECT().Send(deleteResp).Return(nil).AnyTimes()
	pluginServer.PluginHandler.(*mocks.MockPluginHandler).EXPECT().Open(pluginConfig, plugin.OperationTypeDelete).Return(nil).Times(1)
	pluginServer.PluginHandler.(*mocks.MockPluginHandler).EXPECT().Delete("testUUID").Return(nil).Times(1)
	deleteServer.EXPECT().Recv().Return(nil, io.EOF)

	// Act
	err := pluginServer.Delete(deleteServer)

	// Assert
	assert.NoError(err)
}
