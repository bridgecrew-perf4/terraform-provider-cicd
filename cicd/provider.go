// Copyright 2017-2020 Tensigma Ltd. All rights reserved.
// Use of this source code is governed by Microsoft Reference Source
// License (MS-RSL) that can be found in the LICENSE file.

package cicd

import (
	"errors"
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

// Provider creates the Docker provider
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_root": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Pipelines API root",
			},
			"kubernetes_config_path": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("KUBECONFIG", ""),
				Description: "Location of k8s configuration file (used for installing helm charts)",
			},
			"aws_profile": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("AWS_PROFILE", "default"),
				Description: "Name of AWS profile to access S3 configuration bucket (put helm charts)",
			},
			"aws_region": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("AWS_REGION", "eu-central-1"),
				Description: "Name of AWS profile to access S3 configuration bucket (put helm charts)",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"cicd_helm_chart":         resourceHelmChart(),
			"cicd_pipeline_helm":      resourcePipelineHelm(),
			"cicd_pipeline_terraform": resourcePipelineTerraform(),
			"cicd_pipeline_script":    resourcePipelineScript(),
		},
		DataSourcesMap: map[string]*schema.Resource{},
		ConfigureFunc:  providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	root := d.Get("api_root").(string)
	if root == "" {
		return nil, errors.New("api_root is not provided")
	}
	kubeconfig := d.Get("kubernetes_config_path").(string)
	if kubeconfig != "" {
		if _, err := os.Stat(kubeconfig); os.IsNotExist(err) {
			return nil, errors.New("kubernetes_config_path must be a valid file")
		}
	}

	profile := d.Get("aws_profile").(string)
	region := d.Get("aws_region").(string)

	var sess *session.Session
	var err error
	sess, err = session.NewSessionWithOptions(session.Options{
		// Specify profile to load for the session's config
		Profile: profile,
		// Provide SDK Config options, such as Region.
		Config: aws.Config{
			Region: aws.String(region),
		},
		// Force enable Shared Config support
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		return nil, fmt.Errorf("Error creating AWS session: %s", err.Error())
	}
	// validating that we can obtain credentials from aws profile
	if _, errCred := sess.Config.Credentials.Get(); errCred != nil {
		return nil, fmt.Errorf("Error checking AWS profile: %s", errCred.Error())
	}
	return &ProviderConfig{
		ApiRoot:    root,
		Kubeconfig: kubeconfig,
		AwsProfile: profile,
		AwsRegion:  region,
		Session:    sess,
	}, nil
}
