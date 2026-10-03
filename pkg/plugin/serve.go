package plugin

import goplugin "github.com/hashicorp/go-plugin"

// Serve starts the plugin's gRPC server and blocks until orca disconnects. It
// should be the only thing a plugin's main function does.
func Serve(impl CommandProvider) {
	goplugin.Serve(&goplugin.ServeConfig{
		HandshakeConfig: Handshake,
		Plugins: map[string]goplugin.Plugin{
			CommandsPluginName: &CommandProviderPlugin{Impl: impl},
		},
		GRPCServer: goplugin.DefaultGRPCServer,
	})
}
