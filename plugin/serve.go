package plugin

import (
	"github.com/aserto-dev/idp-plugin-sdk/grpcplugin"
	plugin "github.com/hashicorp/go-plugin"
)

type Options struct {
	Handler Handler
}

// Handshake is a common handshake that is shared by plugin and host.
var Handshake = plugin.HandshakeConfig{
	// This isn't required when using VersionedPlugins
	ProtocolVersion:  1,
	MagicCookieKey:   "IDP_PLUGIN",
	MagicCookieValue: "be60172b-2526-432a-865c-12386998a714",
}

// PluginMap is the map of plugins we can dispense.
var PluginMap = plugin.PluginSet{
	"idp-plugin": &grpcplugin.PluginGRPC{},
}

func Serve(options *Options) error {
	pSet := make(plugin.PluginSet)
	pSet["idp-plugin"] = &grpcplugin.PluginGRPC{
		Impl: &AsertoPluginServer{
			Handler: options.Handler,
		},
	}
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: Handshake,
		Plugins:         pSet,

		// A non-nil value here enables gRPC serving for this plugin...
		GRPCServer: plugin.DefaultGRPCServer,
	})
	return nil
}
