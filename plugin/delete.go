package plugin

import (
	"fmt"
	"io"
	"log"

	proto "github.com/aserto-dev/go-grpc/aserto/idpplugin/v1"
	"github.com/aserto-dev/idp-plugin-sdk/config"
	status "google.golang.org/genproto/googleapis/rpc/status"
)

func (s AsertoPluginServer) Delete(srv proto.Plugin_DeleteServer) error {
	errc := make(chan error, 128)
	errDone := make(chan bool, 1)

	defer func() {
		if r := recover(); r != nil {
			errc <- fmt.Errorf("recovering from panic in Delete. error is: %v", r)
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
				&proto.DeleteResponse{
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
			err := s.PluginHandler.Open(cfg, OperationTypeDelete)
			if err != nil {
				return err
			}
			initialized = true
		}

		if userId := req.GetUserId(); userId != "" {
			err := s.PluginHandler.Delete(userId)
			if err != nil {
				errc <- err
			}
		}
	}

	_, err := s.PluginHandler.Close()
	if err != nil {
		errc <- err
	}

	close(errc)
	<-errDone
	return nil
}
