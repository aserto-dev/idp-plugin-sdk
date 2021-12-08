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
	errDone := make(chan struct{}, 1)
	pluginClosed := false

	defer func() {
		close(errc)
		<-errDone

		if !pluginClosed {
			_, err := s.PluginHandler.Close()
			if err != nil {
				log.Println(err.Error())
			}
		}

		if r := recover(); r != nil {
			log.Println(fmt.Errorf("recovering from panic in Export error is: %v", r))
		}
	}()

	go func() {
		defer close(errDone)
		for {
			e, more := <-errc
			if !more {
				// channel closed
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
				return
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
			} else {
				errc <- err
			}
			break
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

	pluginClosed = true
	return nil
}
