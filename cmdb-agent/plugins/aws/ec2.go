package aws

import (
	"context"
	"fmt"
	"cmdb-agent/plugins"
)

type AWSEC2Plugin struct {
	accessKeyID     string
	secretAccessKey string
	region          string
}

func NewAWSEC2Plugin() *AWSEC2Plugin {
	return &AWSEC2Plugin{}
}

func (p *AWSEC2Plugin) Name() string { return "aws_ec2" }
func (p *AWSEC2Plugin) Type() string { return "cloud" }

func (p *AWSEC2Plugin) Init(config map[string]interface{}) error {
	p.accessKeyID, _ = config["access_key_id"].(string)
	p.secretAccessKey, _ = config["secret_access_key"].(string)
	p.region, _ = config["region"].(string)
	if p.accessKeyID == "" || p.secretAccessKey == "" {
		return fmt.Errorf("access_key_id and secret_access_key are required")
	}
	if p.region == "" {
		p.region = "us-east-1"
	}
	return nil
}

func (p *AWSEC2Plugin) Discover(ctx context.Context) ([]plugins.Resource, error) {
	return []plugins.Resource{
		{
			CITypeName: "AWSEC2",
			UniqueKey:  "i-demo-aws-12345",
			Attributes: map[string]interface{}{
				"instance_name": "demo-aws-01",
				"cpu":           4,
				"memory":        8,
				"private_ip":    "172.31.1.10",
				"public_ip":     "54.200.1.1",
				"status":        "running",
				"zone":          "us-east-1a",
				"region":        p.region,
			},
		},
	}, nil
}

func init() {
	plugins.Register("aws_ec2", NewAWSEC2Plugin())
}
