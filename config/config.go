package config

import (
	"encoding/json"
	"strconv"

	"reflect"

	api "github.com/aserto-dev/go-grpc/aserto/api/v1"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/structpb"
)

func NewConfig(pbStruct *structpb.Struct, v interface{}) error {
	configBytes, err := protojson.Marshal(pbStruct)
	if err != nil {
		return err
	}

	err = json.Unmarshal(configBytes, v)
	if err != nil {
		return err
	}

	return nil
}

func ParseApiConfig(cfg interface{}) ([]*api.ConfigElement, error) {
	v := reflect.ValueOf(cfg)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	typeOfS := v.Type()

	var configs []*api.ConfigElement

	for i := 0; i < v.NumField(); i++ {
		field := typeOfS.Field(i)
		tag := field.Tag
		readOnly, err := strconv.ParseBool(tag.Get("readonly"))
		if err != nil {
			return nil, err
		}

		c := &api.ConfigElement{
			Id:          int32(i + 1),
			Kind:        getElementKind(tag.Get("kind")),
			Type:        getElementType(field.Type.String()),
			Name:        field.Name,
			Description: tag.Get("description"),
			Mode:        getDisplayMode(tag.Get("mode")),
			ReadOnly:    readOnly,
		}
		configs = append(configs, c)
	}
	return configs, nil
}

// values set by linker using ldflag -X
var (
	ver    string // nolint:gochecknoglobals // set by linker
	date   string // nolint:gochecknoglobals // set by linker
	commit string // nolint:gochecknoglobals // set by linker
)

func GetVersion() (string, string, string) {
	return ver, date, commit
}

func getElementKind(kind string) api.ConfigElementKind {
	switch kind {
	case "attribute":
		return api.ConfigElementKind_CONFIG_ELEMENT_KIND_ATTRIBUTE
	case "secret":
		return api.ConfigElementKind_CONFIG_ELEMENT_KIND_SECRET
	case "certificate":
		return api.ConfigElementKind_CONFIG_ELEMENT_KIND_CERTIFICATE
	default:
		return api.ConfigElementKind_CONFIG_ELEMENT_KIND_UNKNOWN

	}
}

func getElementType(t string) api.ConfigElementType {
	switch t {
	case "string":
		return api.ConfigElementType_CONFIG_ELEMENT_TYPE_STRING
	case "int":
		return api.ConfigElementType_CONFIG_ELEMENT_TYPE_INTEGER
	case "bool":
		return api.ConfigElementType_CONFIG_ELEMENT_TYPE_BOOLEAN
	default:
		return api.ConfigElementType_CONFIG_ELEMENT_TYPE_UNKNOWN

	}
}

func getDisplayMode(m string) api.DisplayMode {
	switch m {
	case "normal":
		return api.DisplayMode_DISPLAY_MODE_NORMAL
	case "masked":
		return api.DisplayMode_DISPLAY_MODE_MASKED
	default:
		return api.DisplayMode_DISPLAY_MODE_UNKNOWN
	}
}
