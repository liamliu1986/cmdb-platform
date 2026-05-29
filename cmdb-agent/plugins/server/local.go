package server

import (
	"context"
	"os"
	"runtime"

	"cmdb-agent/plugins"
)

// ServerLocalPlugin discovers local server resources
type ServerLocalPlugin struct{}

func NewServerLocalPlugin() *ServerLocalPlugin {
	return &ServerLocalPlugin{}
}

func (p *ServerLocalPlugin) Name() string { return "server_local" }
func (p *ServerLocalPlugin) Type() string { return "server" }

func (p *ServerLocalPlugin) Init(config map[string]interface{}) error {
	return nil
}

func (p *ServerLocalPlugin) Discover(ctx context.Context) ([]plugins.Resource, error) {
	hostname, _ := os.Hostname()
	resources := []plugins.Resource{
		{
			CITypeName: "Server",
			UniqueKey:  hostname,
			Attributes: map[string]interface{}{
				"hostname":   hostname,
				"os":         runtime.GOOS,
				"arch":       runtime.GOARCH,
				"cpu_count":  runtime.NumCPU(),
				"go_version": runtime.Version(),
			},
		},
	}
	return resources, nil
}

func init() {
	plugins.Register("server_local", NewServerLocalPlugin())
}
