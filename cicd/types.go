// Copyright 2017-2021 Tensigma Ltd. All rights reserved.
// Use of this source code is governed by Microsoft Reference Source
// License (MS-RSL) that can be found in the LICENSE file.

package cicd

import (
	"github.com/aws/aws-sdk-go/aws/session"
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
	// Kind of the pipeline: Helm
	Kind PipelineKind `json:"kind"`
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
}
