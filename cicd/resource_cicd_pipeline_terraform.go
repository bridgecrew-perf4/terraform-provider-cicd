// Copyright 2017-2020 Tensigma Ltd. All rights reserved.
// Use of this source code is governed by Microsoft Reference Source
// License (MS-RSL) that can be found in the LICENSE file.

package cicd

import (
	"github.com/AtlantPlatform/infra/terraform-provider-cicd/internal/helpers"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourcePipelineTerraform() *schema.Resource {
	return &schema.Resource{
		Create: onPipelineTerraformCreate,
		Read:   onPipelineTerraformRead,
		Update: onPipelineTerraformUpdate,
		Delete: onPipelineTerraformDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		SchemaVersion: 1,

		Schema: map[string]*schema.Schema{
			"archive": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "terraform plan archive location (stored on s3)",
			},
			"values": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Values for terraform plan and apply",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
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

func onPipelineTerraformCreate(d *schema.ResourceData, m interface{}) error {
	ID := helpers.NewRandSeq(8)
	d.SetId(ID)
	d.Set("secret", helpers.NewRandSeq(32))
	return nil
}

func onPipelineTerraformRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func onPipelineTerraformUpdate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func onPipelineTerraformDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
