package plugin

import (
	"io"
	"log"

	api "github.com/aserto-dev/go-grpc/aserto/api/v1"
	proto "github.com/aserto-dev/go-grpc/aserto/idpplugin/v1"
	"github.com/aserto-dev/idp-plugin-sdk/config"
	status "google.golang.org/genproto/googleapis/rpc/status"
)

func (s AsertoPluginServer) Import(srv proto.Plugin_ImportServer) error {
	errc := make(chan error, 128)
	errDone := make(chan struct{}, 1)
	pluginClosed := false

	defer func() {
		close(errc)
		<-errDone

		s.cleanup(pluginClosed, "Import")
	}()

	initialized := false
	cfg := s.Handler.GetConfig()

	go func() {
		defer close(errDone)
		for {
			e, more := <-errc
			if !more {
				// channel closed
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
				return
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
			break
		}

		if !initialized {
			err = config.NewConfig(req.GetConfig(), cfg)
			if err != nil {
				return err
			}
			err := s.Handler.Open(cfg, OperationTypeWrite)
			if err != nil {
				return err
			}
			initialized = true
		}

		if user := req.GetUser(); user != nil {
			err := s.Handler.Write(user)
			if err != nil {
				errc <- err
			}
		}
	}

	stats, err := s.Handler.Close()
	if err != nil {
		errc <- err
	}

	pluginClosed = true
	if stats != nil {
		err = srv.Send(constructResponse(stats))
		if err != nil {
			return err
		}
	}

	return nil
}

func constructResponse(stats *Stats) *proto.ImportResponse {
	return &proto.ImportResponse{
		Stats: &api.UserProcessStats{
			Received: stats.Received,
			Created:  stats.Created,
			Updated:  stats.Updated,
			Deleted:  stats.Deleted,
			Errors:   stats.Errors,
		}}
}
