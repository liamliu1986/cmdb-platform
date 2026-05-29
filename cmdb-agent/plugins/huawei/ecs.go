package huawei

import (
	"context"
	"fmt"
	"cmdb-agent/plugins"
)

type HuaweiECSPlugin struct {
	accessKeyID     string
	accessKeySecret string
	region          string
}

func NewHuaweiECSPlugin() *HuaweiECSPlugin {
	return &HuaweiECSPlugin{}
}

func (p *HuaweiECSPlugin) Name() string { return "huawei_ecs" }
func (p *HuaweiECSPlugin) Type() string { return "cloud" }

func (p *HuaweiECSPlugin) Init(config map[string]interface{}) error {
	p.accessKeyID, _ = config["access_key_id"].(string)
	p.accessKeySecret, _ = config["access_key_secret"].(string)
	p.region, _ = config["region"].(string)
	if p.accessKeyID == "" || p.accessKeySecret == "" {
		return fmt.Errorf("access_key_id and access_key_secret are required")
	}
	if p.region == "" {
		p.region = "cn-north-4"
	}
	return nil
}

func (p *HuaweiECSPlugin) Discover(ctx context.Context) ([]plugins.Resource, error) {
	return []plugins.Resource{
		{
			CITypeName: "HuaweiECS",
			UniqueKey:  "demo-hw-12345",
			Attributes: map[string]interface{}{
				"instance_name": "demo-huawei-01",
				"cpu":           4,
				"memory":        8,
				"private_ip":    "192.168.2.10",
				"public_ip":     "114.115.1.1",
				"status":        "ACTIVE",
				"zone":          "cn-north-4a",
				"region":        p.region,
			},
		},
	}, nil
}

func init() {
	plugins.Register("huawei_ecs", NewHuaweiECSPlugin())
}
