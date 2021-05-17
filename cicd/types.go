// Copyright 2017-2021 Tensigma Ltd. All rights reserved.
// Use of this source code is governed by Microsoft Reference Source
// License (MS-RSL) that can be found in the LICENSE file.

package cicd

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// providerConfig embeds internal terraform provider configuration
type providerConfig struct {
	APIRoot    string
	Kubeconfig string
	AwsProfile string
	AwsRegion  string
	Session    *session.Session
}

// PipelineKind embeds type of the pipeline
type PipelineKind string

const (
	// PipelineKindHelm is HELM pipeline
	PipelineKindHelm PipelineKind = "helm"
	// PipelineKindTerraform is Terraform pipeline
	PipelineKindTerraform = "terraform"
	// PipelineKindScript is Script pipeline
	PipelineKindScript = "script"
)

// PipelineRef is a secure reference to the pipeline
// (requests without secret matching will not work)
type PipelineRef struct {
	ID     string `json:"id"`
	Secret string `json:"secret"`
}

// PipelineActivateResponse is a response to pipeline activation
type PipelineActivateResponse struct {
	ID     string `json:"id"`
	Secret string `json:"secret"`
}

// PipelineHelmCreate is a structure for HELM pipeline creation updat
type PipelineHelmCreate struct {
	ID string `json:"id"`
	// Secret (required for updates)
	Secret string `json:"secret"`
	// Kind of the pipeline: Helm
	Type PipelineKind `json:"type"`
	// GIT origin to be checked
	Origin string `json:"origin,omitempty"`
	// GIT branches to be checked
	Branches []string `json:"branches,omitempty"`
	// Container registry for docker images
	RegistryURL string `json:"registry_url,omitempty"`
	// Container registrry provider
	RegistryProvider string `json:"registry_provider,omitempty"`
	// only for HELM: Chart ZIP archive location on S3 Bucket
	Archive string `json:"archive,omitempty"`
	// only for HELM: release name. should be Taken from chart.yaml if not speficied
	Release string `json:"release,omitempty"`
	// Chart release namespace. If not specified, 'default' will be used
	Namespace string `json:"namespace,omitempty"`
	// Number of approves required for the pipeline
	ApprovalsRequired int `json:"approvesrequired"`
	// List of approvers who can approve pipeline
	Approvers []string `json:"approvers"`
}

func SafeStringList(d *schema.ResourceData, field string) []string {
	res := []string{}
	if d == nil {
		return res
	}
	if d.Get(field) != nil {
		for _, v := range d.Get(field).([]interface{}) {
			res = append(res, v.(string))
		}
	}
	return res
}

func SafeString(d *schema.ResourceData, field string) string {
	if d == nil {
		return ""
	}
	return fmt.Sprintf("%v", d.Get(field))
}

func SafeNum(d *schema.ResourceData, field string) int {
	if d == nil {
		return 0
	}
	return d.Get(field).(int)
}
