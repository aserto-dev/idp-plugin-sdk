package srv

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/aserto-dev/aserto-idp/shared/pb"
	api "github.com/aserto-dev/go-grpc/aserto/api/v1"
	"github.com/aserto-dev/idp-plugin-sdk/plugin"
	"golang.org/x/sys/unix"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
)

var jsonOptions = protojson.MarshalOptions{
	Multiline:       false,
	Indent:          "  ",
	AllowPartial:    true,
	UseProtoNames:   true,
	UseEnumNumbers:  false,
	EmitUnpopulated: false,
}

type JsonPlugin struct {
	Config  *JsonPluginConfig
	decoder *json.Decoder
	users   bytes.Buffer
	count   int
}

type JsonPluginConfig struct {
	File string `description:"Json file path" kind:"attribute" mode:"normal" readonly:"false"`
}

func (c *JsonPluginConfig) Validate() error {
	dir := filepath.Dir(c.File)

	info, err := os.Stat(dir)
	if err != nil {
		return status.Error(codes.NotFound, err.Error())
	}

	if !info.IsDir() {
		return status.Errorf(codes.InvalidArgument, "%s is not a directory", dir)
	}

	if runtime.GOOS == "windows" {
		if info.Mode().Perm()&(1<<(uint(7))) == 0 {
			return status.Errorf(codes.PermissionDenied, "cannot access %s", dir)
		}
	} else {
		err = unix.Access(dir, unix.W_OK)
		if err != nil {
			return status.Errorf(codes.PermissionDenied, "cannot access %s: %s", dir, err.Error())
		}
	}

	return nil
}

func (c *JsonPluginConfig) Description() string {
	return "JSON plugin"
}

func NewJsonPlugin() *JsonPlugin {
	return &JsonPlugin{
		Config: &JsonPluginConfig{},
	}
}

func (s *JsonPlugin) GetConfig() plugin.PluginConfig {
	return &JsonPluginConfig{}
}

func (s *JsonPlugin) Open(cfg plugin.PluginConfig) error {
	config, ok := cfg.(*JsonPluginConfig)
	if !ok {
		return errors.New("invalid config")
	}
	s.Config = config
	s.count = 0
	return nil
}

func (s *JsonPlugin) Read() ([]*api.User, error) {
	if s.decoder == nil {
		r, err := os.Open(s.Config.File)
		if err != nil {
			return nil, err
		}

		s.decoder = json.NewDecoder(r)

		if _, err = s.decoder.Token(); err != nil {
			log.Println(err)
			return nil, err
		}
	}

	if s.decoder.More() {
		u := api.User{}
		if err := pb.UnmarshalNext(s.decoder, &u); err != nil {
			return nil, err
		}

		return []*api.User{&u}, nil
	} else {
		if _, err := s.decoder.Token(); err != nil {
			return nil, err
		}

		return nil, io.EOF
	}
}

func (s *JsonPlugin) Write(user *api.User) error {
	if s.count == 0 {
		s.users.Write([]byte("[\n"))
	} else {
		_, _ = s.users.Write([]byte(",\n"))
	}
	b, err := jsonOptions.Marshal(user)
	if err != nil {
		return err
	}

	if _, err := s.users.Write(b); err != nil {
		return err
	}
	s.count++

	return nil
}

func (s *JsonPlugin) Close() error {
	_, err := s.users.Write([]byte("\n]\n"))
	if err != nil {
		return err
	}
	f, err := os.Create(s.Config.File)
	if err != nil {
		return err
	}
	w := bufio.NewWriter(f)
	_, err = s.users.WriteTo(w)
	if err != nil {
		return err
	}
	w.Flush()
	return nil
}
