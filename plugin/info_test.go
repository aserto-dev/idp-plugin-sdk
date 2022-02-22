package plugin_test

import (
	"context"
	"testing"

	proto "github.com/aserto-dev/go-grpc/aserto/idpplugin/v1"
	"github.com/aserto-dev/idp-plugin-sdk/mocks"
	"github.com/aserto-dev/idp-plugin-sdk/plugin"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestInfoFail(t *testing.T) {
	// Arrange
	assert := require.New(t)
	ctrl := gomock.NewController(t)
	handler := mocks.NewMockHandler(ctrl)
	pluginServer := &plugin.AsertoPluginServer{Handler: handler}
	pluginConfig := mocks.NewMockConfig(ctrl)
	req := &proto.InfoRequest{}
	pluginServer.Handler.(*mocks.MockHandler).EXPECT().GetConfig().Return(pluginConfig)
	pluginServer.Handler.(*mocks.MockHandler).EXPECT().GetVersion().Return("", "date", "")
	pluginConfig.EXPECT().Description().Return("This is a description")

	// Act
	resp, err := pluginServer.Info(context.Background(), req)

	// Assert
	assert.NoError(err)
	assert.NotNil(resp)
	assert.Equal(resp.GetDescription(), "This is a description")
	assert.Equal(resp.Build.Version, "0.0.0")
	assert.Equal(resp.Build.Date, "date")
	assert.Equal(resp.Build.Commit, "undefined")
}
