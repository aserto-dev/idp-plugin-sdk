package main

import (
	"log"

	"github.com/aserto-dev/idp-plugin-sdk/examples/dummy/pkg/srv"
	"github.com/aserto-dev/idp-plugin-sdk/plugin"
)

func main() {

	dummyPlugin := srv.NewDummyPlugin()

	options := &plugin.Options{
		Handler: dummyPlugin,
	}

	err := plugin.Serve(options)
	if err != nil {
		log.Println(err.Error())
	}
}
