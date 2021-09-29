package main

import (
	"github.com/aserto-dev/idp-plugin-sdk/examples/json/pkg/srv"
	"github.com/aserto-dev/idp-plugin-sdk/plugin"
)

func main() {

	options := &plugin.PluginOptions{
		PluginHandler: &srv.JsonPlugin{},
	}

	plugin.Serve(options)
}
