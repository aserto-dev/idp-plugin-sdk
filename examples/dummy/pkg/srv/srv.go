package srv

import (
	"io"
	"log"

	api "github.com/aserto-dev/go-grpc/aserto/api/v1"
	"github.com/aserto-dev/idp-plugin-sdk/plugin"
)

type DummyPlugin struct {
	Config DummyPluginConfig
}

type DummyPluginConfig struct {
	IntValue    int    `description:"int value" kind:"attribute" mode:"normal" readonly:"false" name:"int_value"`
	StringValue string `description:"string value" kind:"secret" mode:"masked" readonly:"false" name:"string_value"`
}

func (c *DummyPluginConfig) Validate() error {
	log.Printf("Validating %d", c.IntValue)
	log.Printf("Validating %s", c.StringValue)
	return nil
}

func (c *DummyPluginConfig) Description() string {
	return "dummy plugin"
}

func NewDummyPlugin() *DummyPlugin {
	return &DummyPlugin{}
}

func (s DummyPlugin) GetConfig() plugin.PluginConfig {
	return &DummyPluginConfig{}
}

func (s DummyPlugin) Open(config plugin.PluginConfig) error {
	log.Println("Open()")
	return nil
}

func (s DummyPlugin) Read() ([]*api.User, error) {
	log.Println("Received read()")
	return nil, io.EOF
}

func (s DummyPlugin) Write(user *api.User) error {
	log.Printf("Writing user: %v", user)
	return nil
}

func (s DummyPlugin) Delete(userId string) error {
	log.Printf("Deleting user: %s", userId)
	return nil
}

func (s DummyPlugin) Close() error {
	return nil
}
