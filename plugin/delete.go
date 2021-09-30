package plugin

import (
	"fmt"

	proto "github.com/aserto-dev/go-grpc/aserto/idpplugin/v1"
)

func (s AsertoPluginServer) Delete(srv proto.Plugin_DeleteServer) error {
	return fmt.Errorf("not implemented")
}
