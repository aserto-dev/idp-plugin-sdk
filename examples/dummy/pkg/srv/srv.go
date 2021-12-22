package srv

import (
	"io"
	"log"
	"time"

	api "github.com/aserto-dev/go-grpc/aserto/api/v1"
	"github.com/aserto-dev/idp-plugin-sdk/plugin"
)

type DummyPlugin struct {
	Config    DummyPluginConfig
	Operation plugin.OperationType
}

type DummyPluginConfig struct {
	IntValue    int    `description:"int value" kind:"attribute" mode:"normal" readonly:"false" name:"int_value"`
	StringValue string `description:"string value" kind:"secret" mode:"masked" readonly:"false" name:"string_value"`
}

func (c *DummyPluginConfig) Validate(opType plugin.OperationType) error {
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

func (s DummyPlugin) GetVersion() (string, string, string) {
	return "0.0.1", time.Now().UTC().Format(time.RFC3339), ""
}

func (s DummyPlugin) Open(config plugin.PluginConfig, operation plugin.OperationType) error {
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

func (s DummyPlugin) Close() (*plugin.Stats, error) {
	return nil, nil
}
