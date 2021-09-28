package plugin

import (
	"io"
	"log"

	"github.com/aserto-dev/aserto-idp/pkg/proto"
	api "github.com/aserto-dev/go-grpc/aserto/api/v1"
	"github.com/aserto-dev/idp-plugin-sdk/config"
	status "google.golang.org/genproto/googleapis/rpc/status"
)

func (s AsertoPluginServer) Import(srv proto.Plugin_ImportServer) error {
	cfg := s.PluginHandler.GetConfig()
	errc := make(chan error, 128)
	done := make(chan bool, 1)
	subDone := make(chan bool, 1)
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

	users := make(chan *api.User, 128)

	go func() {
		u, more := <-users
		if !more {
			// channel closed
			subDone <- true
			return
		}
		err := s.PluginHandler.Write(u)
		if err != nil {
			errc <- err
		}
	}()

	go func() {
		for {
			req, err := srv.Recv()
			if err == io.EOF {
				done <- true
				return
			}
			if err != nil {
				errc <- err
			}

			if cfg == nil {
				err = config.NewConfig(req.GetConfig(), cfg)
				if err != nil {
					errc <- err
				}
				s.PluginHandler.Open(cfg)
			}

			if user := req.GetUser(); user != nil {
				if u := user.GetUser(); u != nil {
					users <- u
				}
			}
		}
	}()

	<-done
	close(users)
	<-subDone
	close(errc)
	<-errDone
	return nil
}
