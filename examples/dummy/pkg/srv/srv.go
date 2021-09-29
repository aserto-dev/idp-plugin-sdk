package srv

import (
	api "github.com/aserto-dev/go-grpc/aserto/api/v1"
	"github.com/aserto-dev/idp-plugin-sdk/plugin"
)

type DummyPlugin struct {
	Config DummyPluginConfig
}

type DummyPluginConfig struct {
	IntValue    int    `description:"int value" kind:"attribute" mode:"normal" readonly:"false"`
	StringValue string `description:"string value" kind:"secret" mode:"masked" readonly:"false"`
}

func (c *DummyPluginConfig) Validate() error {
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
	return nil
}

func (s DummyPlugin) Read() ([]*api.User, error) {
	return nil, nil
}

func (s DummyPlugin) Write(*api.User) error {
	return nil
}

func (s DummyPlugin) Close() error {
	return nil
}
