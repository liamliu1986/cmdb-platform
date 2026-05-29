package plugins

import "context"

// Resource represents a discovered resource
type Resource struct {
	CITypeName string                 `json:"ci_type_name"`
	UniqueKey  string                 `json:"unique_key"`
	Attributes map[string]interface{} `json:"attributes"`
	Relations  []Relation             `json:"relations,omitempty"`
}

type Relation struct {
	TargetCIType string `json:"target_ci_type"`
	TargetKey    string `json:"target_key"`
	RelationType string `json:"relation_type"`
}

// Discoverer is the interface for all discovery plugins
type Discoverer interface {
	Name() string
	Type() string // cloud/server/network
	Init(config map[string]interface{}) error
	Discover(ctx context.Context) ([]Resource, error)
}

// Registry holds all available plugins
var Registry = make(map[string]Discoverer)

func Register(name string, d Discoverer) {
	Registry[name] = d
}

func Get(name string) (Discoverer, bool) {
	d, ok := Registry[name]
	return d, ok
}
