# idp-plugin-sdk
This SDK enables building aserto idp plugins using Go.

[![License](https://img.shields.io/github/license/aserto-dev/idp-plugin-sdk)](LICENSE)
[![ci](https://github.com/aserto-dev/idp-plugin-sdk/actions/workflows/ci.yaml/badge.svg)](https://github.com/aserto-dev/idp-plugin-sdk/actions/workflows/ci.yaml)

## Tools

- [Git](https://git-scm.com/)
- [Go 1.19](https://golang.org/dl/)
- [Mage](https://magefile.org/)

## Building

[Mage](https://magefile.org/) is used as a tool for building and testing etc.

List available mage targets:

```bash
mage -l
```

Check that your code is compiling:

```bash
mage build
```

## Testing

```bash
mage test
```

## Linting

```bash
mage deps
mage lint
```
## Developing plugins
To start developing your own plugin, you need to implement [PluginHandler](./plugin/plugin.go#L21) and call `plugin.Serve` from your main func.

```go
package main

import (
	"log"

	"github.com/aserto-dev/idp-plugin-sdk/examples/dummy/pkg/srv"
	"github.com/aserto-dev/idp-plugin-sdk/plugin"
)

func main() {

	dummyPlugin := srv.NewDummyPlugin()

	options := &plugin.PluginOptions{
		PluginHandler: dummyPlugin,
	}

	err := plugin.Serve(options)
	if err != nil {
		log.Println(err.Error())
	}
}
```

For a plugin example, check out [aserto-idp-plugin-dummy](./examples/dummy)

## License

[Apache 2.0 License](./LICENSE)