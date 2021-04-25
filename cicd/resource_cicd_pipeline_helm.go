// Copyright 2017-2021 Tensigma Ltd. All rights reserved.
// Use of this source code is governed by Microsoft Reference Source
// License (MS-RSL) that can be found in the LICENSE file.

package cicd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/AtlantPlatform/terraform-provider-cicd/internal/helpers"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourcePipelineHelm() *schema.Resource {
	return &schema.Resource{
		Create: onPipelineHelmCreate,
		Read:   onPipelineHelmRead,
		Update: onPipelineHelmUpdate,
		Delete: onPipelineHelmDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		SchemaVersion: 1,

		Schema: map[string]*schema.Schema{
			"archive": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Chart ZIP archive location, returned from cicd_helm_chart",
			},
			"release": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Chart release name. If not specified, value from Chart.yaml will be used",
			},
			"namespace": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Chart release namespace. If not specified, 'default' will be used",
			},
			"origin": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Git repository (CIRCLE_REPOSITORY_URL). Used for verification of the source",
			},
			"branches": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Git branch that is allowed to be build in this environment (CIRCLE_BRANCH)",
			},
			"registry_url": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Address of the registry for the storage of AWS image (ECR_ACCOUNT_URL)",
			},
			"registry_provider": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Docker provider (aws, ibm, gitlab). AWS by default.",
			},
			// TODO: install on start
			// "install": {
			// 	Type:     schema.TypeBool,
			// 	Optional: true,
			// 	Description: "TODO: Whether to install HELM chart if it is not present",
			// },
			// TODO: approvals: who is allowed to approve
			"secret": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "pipeline secret to use this pipeline",
			},
		},
	}
}

func onPipelineHelmCreate(d *schema.ResourceData, meta interface{}) error {
	apiRoot := meta.(*providerConfig).APIRoot
	payload := PipelineHelmCreate{
		ID:               helpers.NewRandSeq(32),
		Origin:           d.Get("origin").(string),
		Branches:         make([]string, 0),
		RegistryURL:      d.Get("registry_url").(string),
		RegistryProvider: d.Get("registry_provider").(string),
		// helm-specific
		Kind:      PipelineKindHelm,
		Archive:   d.Get("archive").(string),
		Release:   d.Get("release").(string),
		Namespace: d.Get("namespace").(string),
	}
	if d.Get("branches") != nil {
		fmt.Printf("branches=%v\n", d.Get("branches"))
		for _, v := range d.Get("branches").([]interface{}) {
			payload.Branches = append(payload.Branches, v.(string))
		}
	}
	body, _ := json.Marshal(&payload)
	resp, err := http.Post(apiRoot+"/api/pipelines/activate", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("API Init responded with status %v", resp.StatusCode)
	}

	var out PipelineActivateResponse
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(buf, &out); err != nil {
		return err
	}
	if payload.ID != out.ID {
		return fmt.Errorf("IDs don't match, found %s, expected %s", out.ID, payload.ID)
	}
	d.SetId(out.ID)
	d.Set("secret", out.Secret)
	return nil
}

func onPipelineHelmRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func onPipelineHelmUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func onPipelineHelmDelete(d *schema.ResourceData, meta interface{}) error {
	apiRoot := meta.(*providerConfig).APIRoot
	ID := d.Get("id").(string)
	Secret := d.Get("secret").(string)

	payload := PipelineRef{ID: ID, Secret: Secret}
	body, _ := json.Marshal(&payload)
	resp, err := http.Post(apiRoot+"/api/pipelines/deactivate", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("API Pipeline responded with status %v", resp.StatusCode)
	}
	return nil
}
