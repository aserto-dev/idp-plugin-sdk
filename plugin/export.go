package plugin

import (
	"fmt"
	"io"
	"log"

	proto "github.com/aserto-dev/go-grpc/aserto/idpplugin/v1"
	"github.com/aserto-dev/idp-plugin-sdk/config"
	multierror "github.com/hashicorp/go-multierror"
	status "google.golang.org/genproto/googleapis/rpc/status"
)

func (s AsertoPluginServer) Export(req *proto.ExportRequest, srv proto.Plugin_ExportServer) error {
	errc := make(chan error, 128)
	errDone := make(chan bool, 1)
	defer func() {
		if r := recover(); r != nil {
			errc <- fmt.Errorf("recovering from panic in Import error is: %v", r)
		}
	}()

	go func() {
		for {
			e, more := <-errc
			if !more {
				// channel closed
				errDone <- true
				return
			}
			err := srv.Send(
				&proto.ExportResponse{
					Data: &proto.ExportResponse_Error{
						Error: &status.Status{
							Message: e.Error(),
						},
					},
				},
			)
			if err != nil {
				log.Println(err.Error())
			}
		}
	}()

	cfg := s.PluginHandler.GetConfig()
	err := config.NewConfig(req.GetConfig(), cfg)
	if err != nil {
		return err
	}

	err = s.PluginHandler.Open(cfg, OperationTypeRead)
	if err != nil {
		return err
	}

	for {
		users, err := s.PluginHandler.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			if merr, ok := err.(*multierror.Error); ok {
				for _, e := range merr.Errors {
					errc <- e
				}
			}
		}
		for _, u := range users {
			res := &proto.ExportResponse{
				Data: &proto.ExportResponse_User{
					User: u,
				},
			}
			if err = srv.Send(res); err != nil {
				log.Println(err)
				return err
			}
		}
	}

	_, err = s.PluginHandler.Close()
	if err != nil {
		errc <- err
	}

	close(errc)
	<-errDone
	return nil
}
