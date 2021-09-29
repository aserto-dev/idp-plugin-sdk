package plugin

import (
	"io"
	"log"

	"github.com/aserto-dev/aserto-idp/pkg/proto"
	"github.com/aserto-dev/idp-plugin-sdk/config"
	status "google.golang.org/genproto/googleapis/rpc/status"
)

func (s AsertoPluginServer) Import(srv proto.Plugin_ImportServer) error {
	initialized := false
	cfg := s.PluginHandler.GetConfig()
	errc := make(chan error, 128)
	errDone := make(chan bool, 1)

	go func() {
		for {
			e, more := <-errc
			if !more {
				// channel closed
				errDone <- true
				return
			}
			err := srv.Send(
				&proto.ImportResponse{
					Error: &status.Status{
						Message: e.Error(),
					},
				},
			)
			if err != nil {
				log.Println(err.Error())
			}
		}
	}()

	for {
		req, err := srv.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			errc <- err
		}

		if !initialized {
			err = config.NewConfig(req.GetConfig(), cfg)
			if err != nil {
				errc <- err
			}
			err := s.PluginHandler.Open(cfg)
			if err != nil {
				return err
			}
			initialized = true
		}

		if user := req.GetUser(); user != nil {
			if u := user.GetUser(); u != nil {
				err := s.PluginHandler.Write(u)
				if err != nil {
					errc <- err
				}
			}
		}
	}

	err := s.PluginHandler.Close()
	if err != nil {
		errc <- err
	}

	close(errc)
	<-errDone
	return nil
}
