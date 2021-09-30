package plugin

import (
	"fmt"

	"github.com/aserto-dev/aserto-idp/pkg/proto"
)

func (s AsertoPluginServer) Delete(srv proto.Plugin_DeleteServer) error {
	return fmt.Errorf("not implemented")
}
