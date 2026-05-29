package tencent

import (
	"context"
	"fmt"
	"cmdb-agent/plugins"
)

type TencentCVMPlugin struct {
	secretID  string
	secretKey string
	region    string
}

func NewTencentCVMPlugin() *TencentCVMPlugin {
	return &TencentCVMPlugin{}
}

func (p *TencentCVMPlugin) Name() string { return "tencent_cvm" }
func (p *TencentCVMPlugin) Type() string { return "cloud" }

func (p *TencentCVMPlugin) Init(config map[string]interface{}) error {
	p.secretID, _ = config["secret_id"].(string)
	p.secretKey, _ = config["secret_key"].(string)
	p.region, _ = config["region"].(string)
	if p.secretID == "" || p.secretKey == "" {
		return fmt.Errorf("secret_id and secret_key are required")
	}
	if p.region == "" {
		p.region = "ap-guangzhou"
	}
	return nil
}

func (p *TencentCVMPlugin) Discover(ctx context.Context) ([]plugins.Resource, error) {
	// Simplified demo
	return []plugins.Resource{
		{
			CITypeName: "TencentCVM",
			UniqueKey:  "ins-demo-12345",
			Attributes: map[string]interface{}{
				"instance_name": "demo-cvm-01",
				"cpu":           4,
				"memory":        8,
				"private_ip":    "10.0.1.10",
				"public_ip":     "119.29.1.1",
				"status":        "RUNNING",
				"zone":          "ap-guangzhou-2",
				"region":        p.region,
			},
		},
	}, nil
}

func init() {
	plugins.Register("tencent_cvm", NewTencentCVMPlugin())
}
