package plugins

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/adamkirk/orca/pkg/plugin"
	"github.com/hashicorp/go-hclog"
	goplugin "github.com/hashicorp/go-plugin"
)

// LoadedPlugin is a running plugin process, along with the commands it provides.
type LoadedPlugin struct {
	Name     string
	Path     string
	Provider plugin.CommandProvider
	Commands []plugin.CommandSpec
}

type Manager struct {
	logger  *slog.Logger
	clients []*goplugin.Client
	loaded  []LoadedPlugin
}

// Load starts each plugin binary and fetches its commands. A plugin that fails to
// start is logged and skipped, so a broken plugin never stops orca from running.
func (m *Manager) Load(paths []string) []LoadedPlugin {
	for _, p := range paths {
		lp, err := m.load(p)

		if err != nil {
			m.logger.Warn("failed to load plugin", "path", p, "err", err)
			continue
		}

		m.loaded = append(m.loaded, lp)
	}

	return m.loaded
}

func (m *Manager) load(path string) (LoadedPlugin, error) {
	client := goplugin.NewClient(&goplugin.ClientConfig{
		HandshakeConfig:  plugin.Handshake,
		Plugins:          plugin.PluginMap,
		Cmd:              exec.Command(path),
		AllowedProtocols: []goplugin.Protocol{goplugin.ProtocolGRPC},
		SyncStdout:       os.Stdout,
		SyncStderr:       os.Stderr,
		Logger: hclog.New(&hclog.LoggerOptions{
			Name:   "plugin",
			Output: os.Stderr,
			Level:  hclogLevel(m.logger),
		}),
	})

	lp, err := m.connect(client, path)

	if err != nil {
		// Make sure a plugin that failed the handshake isn't left running
		client.Kill()

		return LoadedPlugin{}, err
	}

	m.clients = append(m.clients, client)

	return lp, nil
}

func (m *Manager) connect(client *goplugin.Client, path string) (LoadedPlugin, error) {
	rpcClient, err := client.Client()

	if err != nil {
		return LoadedPlugin{}, err
	}

	raw, err := rpcClient.Dispense(plugin.CommandsPluginName)

	if err != nil {
		return LoadedPlugin{}, err
	}

	provider, ok := raw.(plugin.CommandProvider)

	if !ok {
		return LoadedPlugin{}, fmt.Errorf("plugin does not implement CommandProvider")
	}

	commands, err := provider.Commands()

	if err != nil {
		return LoadedPlugin{}, fmt.Errorf("fetching commands: %w", err)
	}

	return LoadedPlugin{
		Name:     strings.TrimPrefix(filepath.Base(path), BinaryPrefix),
		Path:     path,
		Provider: provider,
		Commands: commands,
	}, nil
}

func (m *Manager) Logger() *slog.Logger {
	return m.logger
}

func (m *Manager) Loaded() []LoadedPlugin {
	return m.loaded
}

// Kill stops all plugin processes. It must be called before orca exits.
func (m *Manager) Kill() {
	for _, c := range m.clients {
		c.Kill()
	}
}

// hclogLevel picks the level for go-plugin's own logging. Its debug and info
// output is lifecycle chatter (process started/exited etc.) on every run, so it's
// only shown when orca is logging at debug, otherwise only warnings and errors.
func hclogLevel(l *slog.Logger) hclog.Level {
	ctx := context.Background()

	switch {
	case l.Enabled(ctx, slog.LevelDebug):
		return hclog.Debug
	case l.Enabled(ctx, slog.LevelWarn):
		return hclog.Warn
	case l.Enabled(ctx, slog.LevelError):
		return hclog.Error
	default:
		return hclog.Off
	}
}

func NewManager(logger *slog.Logger) *Manager {
	return &Manager{
		logger: logger,
	}
}
