// Copyright 2017-2020 Tensigma Ltd. All rights reserved.
// Use of this source code is governed by Microsoft Reference Source
// License (MS-RSL) that can be found in the LICENSE file.

package cicd

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/AtlantPlatform/infra/terraform-provider-cicd/internal/helmchart"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	log "github.com/sirupsen/logrus"
)

func resourceHelmChart() *schema.Resource {
	return &schema.Resource{
		Create: onHelmChartCreate,
		Read:   onHelmChartRead,
		Update: onHelmChartUpdate,
		Delete: onHelmChartDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		SchemaVersion: 1,

		Schema: map[string]*schema.Schema{
			"source": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "local folder where HELM chart is originally located",
			},
			"aws_bucket": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "AWS S3 bucket where ZIP of HELM chart will be uploaded",
			},
			"args": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Arguments for values.yaml substitutions - for its installation",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "output value: name of the chart taken from Chart.yaml",
			},
			"archive": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "output value: file name on AWS S3 bucket",
			},
			"hash": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "hash of ZIP file",
			},
		},
	}
}

const HashMetaHeader = "chart-hash"

func onHelmChartCreate(d *schema.ResourceData, meta interface{}) error {
	cli := s3.New(meta.(*ProviderConfig).Session)

	source := d.Get("source").(string)
	args := d.Get("args").(map[string]interface{})
	chart, err := helmchart.New(source, args)
	if err != nil {
		return err
	}

	// upload it to S3, return its location
	reader, errZip := chart.ZIP()
	if errZip != nil {
		return errZip
	}
	if _, errUpload := cli.PutObject(&s3.PutObjectInput{
		Body:        reader,
		Bucket:      aws.String(d.Get("aws_bucket").(string)),
		ContentType: aws.String("application/zip"),
		Key:         aws.String(chart.GetZipName()),
		Metadata: map[string]*string{
			HashMetaHeader: &chart.Hash,
		},
	}); errUpload != nil {
		return errUpload
	}
	
	d.Set("hash", chart.Hash)
	d.Set("name", chart.Name)
	d.Set("archive", chart.GetZipName())
	d.SetId(chart.ID)
	return nil
}

func onHelmChartRead(d *schema.ResourceData, meta interface{}) error {
	cli := s3.New(meta.(*ProviderConfig).Session)
	if d.Get("id") != nil && d.Get("name") != nil {
		chart := &helmchart.Builder{
			ID:   d.Get("id").(string),
			Name: d.Get("name").(string),
		}
		// check s3 location
		sourceObject, errRead := cli.HeadObject(&s3.HeadObjectInput{
			Bucket: aws.String(d.Get("aws_bucket").(string)),
			Key:    aws.String(chart.GetZipName()),
		}); 
		if errRead != nil {
			// no error, but nothing is present, need regeneration
			d.Set("archive", "")
			d.Set("hash", "")
			return nil
		}
		if sourceObject.Metadata != nil {
			d.Set("archive", chart.GetZipName())
			if hash, ok := sourceObject.Metadata[HashMetaHeader]; ok {
				d.Set("hash", hash)
			}
		}
	}
	return nil
}

// in theory, we should never get to update
// on chart modification, as ID should be changed, 
// old resource is a subject of deletion
func onHelmChartUpdate(d *schema.ResourceData, meta interface{}) error {
	cli := s3.New(meta.(*ProviderConfig).Session)

	source := d.Get("source").(string)
	args := d.Get("args").(map[string]interface{})
	chart, err := helmchart.New(source, args)
	if err != nil {
		return err
	}

	// upload it to S3, return its location
	reader, errZip := chart.ZIP()
	if errZip != nil {
		return errZip
	}
	if _, errUpload := cli.PutObject(&s3.PutObjectInput{
		Body:        reader,
		Bucket:      aws.String(d.Get("aws_bucket").(string)),
		ContentType: aws.String("application/zip"),
		Key:         aws.String(chart.GetZipName()),
		Metadata: map[string]*string{
			HashMetaHeader: aws.String(chart.Hash),
		},
	}); errUpload != nil {
		return errUpload
	}
	
	d.Set("hash", chart.Hash)
	d.Set("name", chart.Name)
	d.Set("archive", chart.GetZipName())
	return nil
}

func onHelmChartDelete(d *schema.ResourceData, meta interface{}) error {
	// remove file from AWS s3 bucket
	if d.Get("id") != nil && d.Get("name") != nil {
		chart := &helmchart.Builder{ID: d.Get("id").(string), Name: d.Get("name").(string)}
		session := meta.(*ProviderConfig).Session
		if _, errDelete := s3.New(session).DeleteObject(&s3.DeleteObjectInput{
			Bucket: aws.String(d.Get("aws_bucket").(string)),
			Key:    aws.String(chart.GetZipName()),
		}); errDelete != nil {
			log.WithError(errDelete).Warn("removal failed")
		}
	}
	return nil
}
