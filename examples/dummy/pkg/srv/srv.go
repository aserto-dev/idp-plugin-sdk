package srv

import (
	"io"
	"log"
	"time"

	api "github.com/aserto-dev/go-grpc/aserto/api/v1"
	"github.com/aserto-dev/idp-plugin-sdk/examples/dummy/pkg/config"
	"github.com/aserto-dev/idp-plugin-sdk/plugin"
)

type DummyPlugin struct {
	Config    config.DummyPluginConfig
	Operation plugin.OperationType
}

func NewDummyPlugin() *DummyPlugin {
	return &DummyPlugin{}
}

func (s DummyPlugin) GetConfig() plugin.Config {
	return &config.DummyPluginConfig{}
}

func (s DummyPlugin) GetVersion() (string, string, string) {
	return "0.0.1", time.Now().UTC().Format(time.RFC3339), ""
}

func (s DummyPlugin) Open(config plugin.Config, operation plugin.OperationType) error {
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
