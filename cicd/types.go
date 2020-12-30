// Copyright 2017-2020 Tensigma Ltd. All rights reserved.
// Use of this source code is governed by Microsoft Reference Source
// License (MS-RSL) that can be found in the LICENSE file.

package cicd

import (
	"github.com/aws/aws-sdk-go/aws/session"
)

type ProviderConfig struct {
	ApiRoot    string
	Kubeconfig string
	AwsProfile string
	AwsRegion  string
	Session    *session.Session
}

// PipelineKind embeds type of the pipeline
type PipelineKind string

const (
	PipelineKindHelm PipelineKind = "Helm"
	PipelineKindTerraform = "Terraform"
	PipelineKindScript = "Script"
)

type PipelineRef struct {
	ID string `json:"id"`
    Secret string `json:"secret"`
}

type PipelineInitResponse struct {
	ID string `json:"id"`
	Secret string `json:"secret"`
}

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