// Copyright 2017-2021 Tensigma Ltd. All rights reserved.
// Use of this source code is governed by Microsoft Reference Source
// License (MS-RSL) that can be found in the LICENSE file.

package cicd

import (
	"github.com/AtlantPlatform/terraform-provider-cicd/internal/helpers"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourcePipelineScript() *schema.Resource {
	return &schema.Resource{
		Create: onPipelineScriptCreate,
		Read:   onPipelineScriptRead,
		Update: onPipelineScriptUpdate,
		Delete: onPipelineScriptDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		SchemaVersion: 1,

		Schema: map[string]*schema.Schema{
			"exec": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "command for exectution",
			},
			"plan": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "(optional) command to preview before the execution",
			},
			"env": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "(optional) environment variables for execution",
			},
			// TODO: approvals
			"secret": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "pipeline secret to use this pipeline",
			},
		},
	}
}

func onPipelineScriptCreate(d *schema.ResourceData, m interface{}) error {
	ID := helpers.NewRandSeq(8)
	d.SetId(ID)
	d.Set("secret", helpers.NewRandSeq(32))
	return nil
}

func onPipelineScriptRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func onPipelineScriptUpdate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func onPipelineScriptDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
