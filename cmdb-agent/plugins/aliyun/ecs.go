package aliyun

import (
	"context"
	"fmt"
	"cmdb-agent/plugins"
)

// AliyunECSPlugin discovers Alibaba Cloud ECS instances
type AliyunECSPlugin struct {
	accessKeyID     string
	accessKeySecret string
	regionID        string
}

func NewAliyunECSPlugin() *AliyunECSPlugin {
	return &AliyunECSPlugin{}
}

func (p *AliyunECSPlugin) Name() string { return "aliyun_ecs" }
func (p *AliyunECSPlugin) Type() string { return "cloud" }

func (p *AliyunECSPlugin) Init(config map[string]interface{}) error {
	p.accessKeyID, _ = config["access_key_id"].(string)
	p.accessKeySecret, _ = config["access_key_secret"].(string)
	p.regionID, _ = config["region_id"].(string)
	if p.accessKeyID == "" || p.accessKeySecret == "" {
		return fmt.Errorf("access_key_id and access_key_secret are required")
	}
	if p.regionID == "" {
		p.regionID = "cn-hangzhou"
	}
	return nil
}

func (p *AliyunECSPlugin) Discover(ctx context.Context) ([]plugins.Resource, error) {
	// Simplified: return a demo resource
	// In production, use Alibaba Cloud SDK to list ECS instances
	return []plugins.Resource{
		{
			CITypeName: "AliyunECS",
			UniqueKey:  "i-demo-12345",
			Attributes: map[string]interface{}{
				"instance_name": "demo-server-01",
				"cpu":           4,
				"memory":        8,
				"private_ip":    "192.168.1.10",
				"public_ip":     "47.100.1.1",
				"status":        "Running",
				"zone":          "cn-hangzhou-a",
				"region":        p.regionID,
			},
		},
	}, nil
}

func init() {
	plugins.Register("aliyun_ecs", NewAliyunECSPlugin())
}
