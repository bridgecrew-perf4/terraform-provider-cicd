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
	log "github.com/sirupsen/logrus"
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
				Required:    true,
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
			"approvals_required": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "Number of approvals required for the pipeline to be finished",
			},
			"approvers": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Comma-separated list of approvers",
			},
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
		Origin:           SafeString(d, "origin"),
		Branches:         make([]string, 0),
		RegistryURL:      SafeString(d, "registry_url"),
		RegistryProvider: SafeString(d, "registry_provider"),
		// helm-specific
		Type:              PipelineKindHelm,
		Archive:           SafeString(d, "archive"),
		Release:           SafeString(d, "release"),
		Namespace:         SafeString(d, "namespace"),
		ApprovalsRequired: SafeNum(d, "approvals_required"),
	}
	if d.Get("approvers") != nil {
		fmt.Printf("approvers=%v\n", d.Get("approvers"))
		for _, v := range d.Get("approvers").([]interface{}) {
			payload.Approvers = append(payload.Approvers, v.(string))
		}
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
		return fmt.Errorf("activation error: %v", err.Error())
	}
	defer resp.Body.Close()

	var out PipelineActivateResponse
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode >= 300 {
		return fmt.Errorf("API Init responded with status %v (%v)",
			resp.StatusCode, string(buf))
	}

	if err := json.Unmarshal(buf, &out); err != nil {
		return err
	}
	if payload.ID != out.ID {
		return fmt.Errorf("IDs don't match, found %s, expected %s (%v)",
			out.ID, payload.ID, string(buf))
	}
	d.SetId(out.ID)
	// secret comes back from pipelines server
	d.Set("secret", out.Secret)
	return nil
}

func onPipelineHelmRead(d *schema.ResourceData, meta interface{}) error {
	// nothing here so far. We fully trust local store
	return nil
}

func onPipelineHelmUpdate(d *schema.ResourceData, meta interface{}) error {
	apiRoot := meta.(*providerConfig).APIRoot
	if len(d.Id()) == 0 {
		return nil
	}
	payload := PipelineHelmCreate{
		ID:               d.Id(),
		Secret:           SafeString(d, "secret"),
		Origin:           SafeString(d, "origin"),
		Branches:         make([]string, 0),
		RegistryURL:      SafeString(d, "registry_url"),
		RegistryProvider: SafeString(d, "registry_provider"),
		// helm-specific
		Type:              PipelineKindHelm,
		Archive:           SafeString(d, "archive"),
		Release:           SafeString(d, "release"),
		Namespace:         SafeString(d, "namespace"),
		ApprovalsRequired: SafeNum(d, "approvals_required"),
	}
	if d.Get("approvers") != nil {
		fmt.Printf("approvers=%v\n", d.Get("approvers"))
		for _, v := range d.Get("approvers").([]interface{}) {
			payload.Approvers = append(payload.Approvers, v.(string))
		}
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
		return fmt.Errorf("activation error: %v", err.Error())
	}
	defer resp.Body.Close()
	var out PipelineActivateResponse
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode >= 300 {
		return fmt.Errorf("API Init responded with status %v (%v)",
			resp.StatusCode, string(body))
	}

	if err := json.Unmarshal(buf, &out); err != nil {
		return err
	}
	if payload.ID != out.ID {
		return fmt.Errorf("IDs don't match, found %s, expected %s", out.ID, payload.ID)
	}
	return nil
}

// all errors of deactivation are silenced
func onPipelineHelmDelete(d *schema.ResourceData, meta interface{}) error {
	apiRoot := meta.(*providerConfig).APIRoot
	if len(d.Id()) == 0 {
		return nil
	}
	Secret := SafeString(d, "secret")
	payload := PipelineRef{ID: d.Id(), Secret: Secret}

	fmt.Println("onPipelineHelmDelete: ", payload)
	body, _ := json.Marshal(&payload)
	resp, err := http.Post(apiRoot+"/api/pipelines/deactivate", "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.WithError(err).WithField("payload", payload).Warn("deactivation error (silenced)")
		return nil
	}
	defer resp.Body.Close()
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode >= 300 {
		log.WithError(err).
			WithField("payload", payload).
			WithField("body", string(buf)).
			WithField("status", resp.StatusCode).Warn("deactivation bad status code (silenced)")
		return nil
	}
	fmt.Println("onPipelineHelmDelete: done")
	return nil
}
