package config

import (
	"log"

	"github.com/aserto-dev/idp-plugin-sdk/plugin"
)

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
