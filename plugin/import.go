package plugin

import (
	"fmt"
	"io"
	"log"

	api "github.com/aserto-dev/go-grpc/aserto/api/v1"
	proto "github.com/aserto-dev/go-grpc/aserto/idpplugin/v1"
	"github.com/aserto-dev/idp-plugin-sdk/config"
	status "google.golang.org/genproto/googleapis/rpc/status"
)

func (s AsertoPluginServer) Import(srv proto.Plugin_ImportServer) error {
	errc := make(chan error, 128)
	errDone := make(chan bool, 1)

	defer func() {
		if r := recover(); r != nil {
			errc <- fmt.Errorf("recovering from panic in Import error is: %v", r)
		}
	}()

	initialized := false
	cfg := s.PluginHandler.GetConfig()

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
			err := s.PluginHandler.Open(cfg, OperationTypeWrite)
			if err != nil {
				return err
			}
			initialized = true
		}

		if user := req.GetUser(); user != nil {
			err := s.PluginHandler.Write(user)
			if err != nil {
				errc <- err
			}
		}
	}

	stats, err := s.PluginHandler.Close()
	if err != nil {
		errc <- err
	}

	if stats != nil {
		err = srv.Send(
			&proto.ImportResponse{
				Stats: &api.UserProcessStats{
					Received: stats.Received,
					Created:  stats.Created,
					Updated:  stats.Updated,
					Deleted:  stats.Deleted,
					Errors:   stats.Errors,
				}},
		)
		if err != nil {
			errc <- err
		}
	}

	close(errc)
	<-errDone
	return nil
}
