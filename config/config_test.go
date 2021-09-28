package config

import (
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/structpb"
)

type TestConfig struct {
	IntValue    int    `description:"int value" kind:"attribute" mode:"normal" readonly:"false"`
	StringValue string `description:"string value" kind:"secret" mode:"masked" readonly:"false"`
	BoolValue   bool   `description:"string value" kind:"secret" mode:"masked" readonly:"false"`
}

func TestNewConfig(t *testing.T) {
	assert := require.New(t)

	m, err := structpb.NewStruct(map[string]interface{}{
		"IntValue":    1,
		"StringValue": "str",
		"BoolValue":   true,
	})
	if err != nil {
		t.Error(err)
	}
	v := &TestConfig{}
	err = NewConfig(m, v)
	assert.Nil(err)
	assert.NotNil(v)
	assert.Equal(1, v.IntValue)
	assert.Equal("str", v.StringValue)
	assert.Equal(true, v.BoolValue)
}

func TestApiConfig(t *testing.T) {
	assert := require.New(t)

	cfg, err := ParseApiConfig(TestConfig{})
	assert.Nil(err)
	assert.NotNil(cfg)
	assert.Equal(3, len(cfg))
}
