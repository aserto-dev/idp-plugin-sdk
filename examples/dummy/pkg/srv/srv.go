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

// func Send(user *api.User) error {
// 	return nil
// }

// func Receive() ([]*api.User, error) {
// 	return nil, nil
// }

// func (pl *DummyPluginConfig) Validate() error {
// 	return nil
// }

// func Write(config interface{}, send func(user *api.User) error) error {
// 	// init auth0 client
// 	return err
// 	//
// 	err := send(user)
// }

// func Read(config interface{}, recv func() (user *api.User, error)) error {
// 	// init auth0 client
// 	return err
// 	//
// 	u := recv(user)
// }

// func (s DummyPluginServer) Import(srv proto.Plugin_ImportServer) error {
// 	log.Println("not implemented")
// 	return nil
// }

// // func (s DummyPluginServer) Delete(srv proto.Plugin_DeleteServer) error {
// // 	return fmt.Errorf("not implemented")
// // }

// func (DummyPluginServer) Validate(ctx context.Context, req *proto.ValidateRequest) (*proto.ValidateResponse, error) {
// 	response := &proto.ValidateResponse{}
// 	return response, nil
// }

// func (s DummyPluginServer) Export(config interface{}, send func(*api.User), errf func(error)) error {
// 	cfg := config.(JsonConfig)

// 	r, err := os.Open(cfg.File)
// 	if err != nil {
// 		log.Println(err)
// 		return err
// 	}

// 	dec := json.NewDecoder(r)

// 	if _, err = dec.Token(); err != nil {
// 		log.Println(err)
// 		return err
// 	}

// 	for dec.More() {
// 		u := api.User{}
// 		if err := pb.UnmarshalNext(dec, &u); err != nil {
// 			log.Println(err)
// 			return err
// 		}
// 		send(u)
// 	}

// 	if _, err = dec.Token(); err != nil {
// 		log.Println(err)
// 		return err
// 	}
// 	return nil
// }

// func (s JsonPluginServer) Import(config interface{}, recv func() *api.User, errf func(error)) error {
// 	config := config.(JsonConfig)
// 	done := make(chan bool, 1)
// 	count := int32(0)

// 	go grpcerr.SendImportErrors(srv, errc, errDone)

// 	var users bytes.Buffer
// 	users.Write([]byte("[\n"))

// 	options := protojson.MarshalOptions{
// 		Multiline:       false,
// 		Indent:          "  ",
// 		AllowPartial:    true,
// 		UseProtoNames:   true,
// 		UseEnumNumbers:  false,
// 		EmitUnpopulated: false,
// 	}

// 	go func() {
// 		for {
// 			if user := recv(); user != nil {
// 				if count > 0 {
// 					_, _ = users.Write([]byte(",\n"))
// 				}
// 				b, err := options.Marshal(u)
// 				if err != nil {
// 					errc <- err
// 				}

// 				if _, err := users.Write(b); err != nil {
// 					errc <- err
// 				}
// 				count++
// 			}
// 			else {
// 				done <- true
// 				return
// 			}
// 		}
// 	}()
// 	// Wait for done
// 	<-done

// 	_, err := users.Write([]byte("\n]\n"))
// 	if err != nil {
// 		errf(err)
// 	}
// 	f, err := os.Create(config.File)
// 	if err != nil {
// 		errf(err)
// 	}
// 	w := bufio.NewWriter(f)
// 	_, err = users.WriteTo(w)
// 	if err != nil {
// 		errf(err)
// 	}
// 	w.Flush()

// 	return nil
// }
