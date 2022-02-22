package config

import (
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/structpb"
)

type TestConfig struct {
	IntValue    int     `description:"int value" kind:"attribute" mode:"normal" readonly:"false" name:"int_value"`
	StringValue string  `description:"string value" kind:"secret" mode:"masked" readonly:"false" name:"string_value"`
	BoolValue   bool    `description:"bool value" kind:"secret" mode:"masked" readonly:"false" name:"bool_value"`
	FloatValue  float32 `description:"float32 value" kind:"something" mode:"supernatural" readonly:"false" name:"float_value"`
}

func TestNewConfig(t *testing.T) {
	assert := require.New(t)

	m, err := structpb.NewStruct(map[string]interface{}{
		"int_value":    1,
		"string_value": "str",
		"bool_value":   true,
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

	cfg, err := ParseAPIConfig(TestConfig{})
	assert.Nil(err)
	assert.NotNil(cfg)
	assert.Equal(4, len(cfg))
}

func TestParameterName(t *testing.T) {
	assert := require.New(t)

	cfg, err := ParseAPIConfig(TestConfig{})
	assert.Nil(err)
	assert.NotNil(cfg)
	assert.Equal(4, len(cfg))
	assert.Equal("int_value", cfg[0].Name)
	assert.Equal("string_value", cfg[1].Name)
	assert.Equal("bool_value", cfg[2].Name)
	assert.Equal("float_value", cfg[3].Name)
}

type ConfigNoName struct {
	IntValue    int    `description:"int value" kind:"attribute" mode:"normal" readonly:"false"`
	StringValue string `description:"string value" kind:"secret" mode:"masked" readonly:"false"`
	BoolValue   bool   `description:"string value" kind:"secret" mode:"masked" readonly:"false"`
}

func TestParameterTagMissing(t *testing.T) {

	assert := require.New(t)

	cfg, err := ParseAPIConfig(ConfigNoName{})
	assert.Nil(err)
	assert.NotNil(cfg)
	assert.Equal(3, len(cfg))
	assert.Equal("intvalue", cfg[0].Name)
	assert.Equal("stringvalue", cfg[1].Name)
	assert.Equal("boolvalue", cfg[2].Name)
}

func TestNewConfigMissingNameTag(t *testing.T) {
	assert := require.New(t)

	m, err := structpb.NewStruct(map[string]interface{}{
		"intvalue":    1,
		"stringvalue": "str",
		"boolvalue":   true,
	})
	if err != nil {
		t.Error(err)
	}
	v := &ConfigNoName{}
	err = NewConfig(m, v)
	assert.Nil(err)
	assert.NotNil(v)
	assert.Equal(1, v.IntValue)
	assert.Equal("str", v.StringValue)
	assert.Equal(true, v.BoolValue)
}
