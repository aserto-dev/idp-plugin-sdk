package mocks

//go:generate mockgen -destination=mock_plugin.go -package=mocks github.com/aserto-dev/idp-plugin-sdk/plugin PluginHandler,PluginConfig
//go:generate mockgen -destination=mock_servers.go -package=mocks github.com/aserto-dev/go-grpc/aserto/idpplugin/v1 Plugin_DeleteServer,Plugin_ExportServer,Plugin_ImportServer
