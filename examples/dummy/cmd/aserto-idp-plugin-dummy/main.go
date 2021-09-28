package main

import (
	"github.com/aserto-dev/idp-plugin-sdk/examples/dummy/pkg/srv"
	"github.com/aserto-dev/idp-plugin-sdk/plugin"
)

func main() {

	options := &plugin.PluginOptions{
		PluginHandler: srv.DummyPlugin{},
	}

	plugin.Serve(options)
}
