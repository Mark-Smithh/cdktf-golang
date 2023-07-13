package main

import (
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/hashicorp/terraform-cdk-go/cdktf"

	"github.com/cdktf/cdktf-provider-aws-go/aws/v16/s3bucket"

	"github.com/cdktf/cdktf-provider-aws-go/aws/v16/instance"
	awsprovider "github.com/cdktf/cdktf-provider-aws-go/aws/v16/provider"
)

func NewMyStack(scope constructs.Construct, id string) cdktf.TerraformStack {
	stack := cdktf.NewTerraformStack(scope, &id)

	awsprovider.NewAwsProvider(stack, jsii.String("AWS"), &awsprovider.AwsProviderConfig{
		Region: jsii.String("us-west-1"),
	})

	instance := instance.NewInstance(stack, jsii.String("compute"), &instance.InstanceConfig{
		Ami:          jsii.String("ami-01456a894f71116f2"),
		InstanceType: jsii.String("t2.micro"),
		Tags: &map[string]*string{
			"Name":        jsii.String("cdktf-golang"),
			"CreatedWith": jsii.String("cdktf"), // Cloud Development Kit Terraform
		},
	})

	bucket := s3bucket.NewS3Bucket(stack,
		jsii.String("bucket"),
		&s3bucket.S3BucketConfig{},
	)

	cdktf.NewTerraformOutput(stack, jsii.String("public_ip"), &cdktf.TerraformOutputConfig{
		Value: instance.PublicIp(),
	})

	cdktf.NewTerraformOutput(stack, jsii.String("bucket_id"), &cdktf.TerraformOutputConfig{
		Value: bucket.Id(),
	})

	return stack
}

func main() {
	app := cdktf.NewApp(nil)
	NewMyStack(app, "aws_instance")

	cdktf.NewRemoteBackend(app, &cdktf.RemoteBackendConfig{
		Hostname:     jsii.String("app.terraform.io"),                                      // Terraform Cloud Hostname
		Organization: jsii.String("terraform_cdp_m_smith"),                                 // Terraform Cloud Organization
		Workspaces:   cdktf.NewNamedRemoteWorkspace(jsii.String("terraform_cdktf_golang")), // Terraform Cloud Workspace
	})

	app.Synth()
}
