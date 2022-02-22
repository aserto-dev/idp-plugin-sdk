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
	handler := mocks.NewMockHandler(ctrl)
	pluginServer := &plugin.AsertoPluginServer{Handler: handler}
	pluginConfig := mocks.NewMockConfig(ctrl)
	deleteServer := mocks.NewMockPlugin_DeleteServer(ctrl)

	pluginServer.Handler.(*mocks.MockHandler).EXPECT().GetConfig().Return(pluginConfig)
	deleteServer.EXPECT().Recv().Return(nil, io.EOF)
	pluginServer.Handler.(*mocks.MockHandler).EXPECT().Close().Return(nil, nil)

	// Act
	err := pluginServer.Delete(deleteServer)

	// Assert
	assert.NoError(err)
}

func TestDeleteNoUsersFailClose(t *testing.T) {
	// Arrange
	assert := require.New(t)
	ctrl := gomock.NewController(t)
	handler := mocks.NewMockHandler(ctrl)
	pluginServer := &plugin.AsertoPluginServer{Handler: handler}
	pluginConfig := mocks.NewMockConfig(ctrl)
	deleteServer := mocks.NewMockPlugin_DeleteServer(ctrl)
	boomErr := errors.New("#boom#")
	deleteResp := &proto.DeleteResponse{Error: &status.Status{Message: boomErr.Error()}}

	pluginServer.Handler.(*mocks.MockHandler).EXPECT().GetConfig().Return(pluginConfig)
	deleteServer.EXPECT().Recv().Return(nil, io.EOF)
	pluginServer.Handler.(*mocks.MockHandler).EXPECT().Close().Return(nil, boomErr)
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
	handler := mocks.NewMockHandler(ctrl)
	pluginServer := &plugin.AsertoPluginServer{Handler: handler}
	pluginConfig := mocks.NewMockConfig(ctrl)
	deleteServer := mocks.NewMockPlugin_DeleteServer(ctrl)
	boomErr := errors.New("#boom#")
	deleteResp := &proto.DeleteResponse{Error: &status.Status{Message: boomErr.Error()}}

	pluginServer.Handler.(*mocks.MockHandler).EXPECT().GetConfig().Return(pluginConfig)
	deleteServer.EXPECT().Recv().Return(nil, boomErr)
	pluginServer.Handler.(*mocks.MockHandler).EXPECT().Close().Return(nil, nil)
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
	handler := mocks.NewMockHandler(ctrl)
	pluginServer := &plugin.AsertoPluginServer{Handler: handler}
	pluginConfig := mocks.NewMockConfig(ctrl)
	deleteServer := mocks.NewMockPlugin_DeleteServer(ctrl)
	deleteResp := &proto.DeleteResponse{Error: nil}

	pluginServer.Handler.(*mocks.MockHandler).EXPECT().GetConfig().Return(pluginConfig)
	deleteServer.EXPECT().Recv().Return(nil, nil)
	pluginServer.Handler.(*mocks.MockHandler).EXPECT().Close().Return(nil, nil)
	deleteServer.EXPECT().Send(deleteResp).Return(nil).AnyTimes()
	pluginServer.Handler.(*mocks.MockHandler).EXPECT().Open(pluginConfig, plugin.OperationTypeDelete).Return(nil).Times(1)
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
	handler := mocks.NewMockHandler(ctrl)
	pluginServer := &plugin.AsertoPluginServer{Handler: handler}
	pluginConfig := mocks.NewMockConfig(ctrl)
	deleteServer := mocks.NewMockPlugin_DeleteServer(ctrl)
	deleteResp := &proto.DeleteResponse{Error: nil}
	deleteReq := &proto.DeleteRequest{Data: &proto.DeleteRequest_UserId{UserId: "testUUID"}}

	pluginServer.Handler.(*mocks.MockHandler).EXPECT().GetConfig().Return(pluginConfig)
	deleteServer.EXPECT().Recv().Return(deleteReq, nil)
	pluginServer.Handler.(*mocks.MockHandler).EXPECT().Close().Return(nil, nil)
	deleteServer.EXPECT().Send(deleteResp).Return(nil).AnyTimes()
	pluginServer.Handler.(*mocks.MockHandler).EXPECT().Open(pluginConfig, plugin.OperationTypeDelete).Return(nil).Times(1)
	pluginServer.Handler.(*mocks.MockHandler).EXPECT().Delete("testUUID").Return(nil).Times(1)
	deleteServer.EXPECT().Recv().Return(nil, io.EOF)

	// Act
	err := pluginServer.Delete(deleteServer)

	// Assert
	assert.NoError(err)
}
