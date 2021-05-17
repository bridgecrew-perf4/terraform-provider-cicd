// Copyright 2017-2021 Tensigma Ltd. All rights reserved.
// Use of this source code is governed by Microsoft Reference Source
// License (MS-RSL) that can be found in the LICENSE file.

package cicd

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/AtlantPlatform/terraform-provider-cicd/internal/helmchart"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
			"allowed": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "list of allowed parameters to be overwritten",
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

// HashMetaHeader added to s3 as meta-data
const HashMetaHeader = "chart-hash"

func onHelmChartCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("onHelmChartCreate: start %v", d)
	cli := s3.New(meta.(*providerConfig).Session)

	source := d.Get("source").(string)
	args := d.Get("args").(map[string]interface{})
	chart, err := helmchart.New(source, args, SafeStringList(d, "allowed"))
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
		Bucket:      aws.String(SafeString(d, "aws_bucket")),
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
	// cli := s3.New(meta.(*providerConfig).Session)
	log.Printf("onHelmChartRead: start %v", d)
	if d.Get("name") != nil {
		// 1. reading local chart information
		source := d.Get("source").(string)
		args := d.Get("args").(map[string]interface{})
		localChart, err := helmchart.New(source, args, SafeStringList(d, "allowed"))
		if err != nil {
			return err
		}
		d.Set("hash", localChart.Hash)
		d.Set("archive", localChart.GetZipName())
		log.Printf("onHelmChartRead: checksum: %v %v", localChart.Hash, d)

		// 2. reading external information
		// remoteChart := &helmchart.Builder{
		// 	ID:   d.Id(),
		// 	Name: d.Get("name").(string),
		// }
		// // check s3 location
		// _, errRead := cli.HeadObject(&s3.HeadObjectInput{
		// 	Bucket: aws.String(d.Get("aws_bucket").(string)),
		// 	Key:    aws.String(remoteChart.GetZipName()),
		// })
		// if errRead != nil {
		// 	// no error, but nothing is present on s3, need regeneration
		// 	return nil
		// }
	}
	return nil
}

// in theory, we should never get to update
// on chart modification, as ID should be changed,
// old resource is a subject of deletion
func onHelmChartUpdate(d *schema.ResourceData, meta interface{}) error {
	cli := s3.New(meta.(*providerConfig).Session)

	source := d.Get("source").(string)
	args := d.Get("args").(map[string]interface{})
	localChart, err := helmchart.New(source, args, SafeStringList(d, "allowed"))
	if err != nil {
		return err
	}

	// upload it to S3, return its location
	reader, errZip := localChart.ZIP()
	if errZip != nil {
		return errZip
	}
	if _, errUpload := cli.PutObject(&s3.PutObjectInput{
		Body:        reader,
		Bucket:      aws.String(SafeString(d, "aws_bucket")),
		ContentType: aws.String("application/zip"),
		Key:         aws.String(localChart.GetZipName()),
		Metadata: map[string]*string{
			HashMetaHeader: aws.String(localChart.Hash),
		},
	}); errUpload != nil {
		return errUpload
	}

	d.Set("hash", localChart.Hash)
	d.Set("name", localChart.Name)
	d.Set("archive", localChart.GetZipName())
	return nil
}

func onHelmChartDelete(d *schema.ResourceData, meta interface{}) error {
	// remove file from AWS s3 bucket
	if d.Get("name") != nil {
		remoteChart := &helmchart.Builder{
			ID:   d.Id(),
			Name: SafeString(d, "name"),
		}
		session := meta.(*providerConfig).Session
		if _, errDelete := s3.New(session).DeleteObject(&s3.DeleteObjectInput{
			Bucket: aws.String(SafeString(d, "aws_bucket")),
			Key:    aws.String(remoteChart.GetZipName()),
		}); errDelete != nil {
			log.Printf("[WARN] removal failed %v", errDelete)
		}
	}
	return nil
}
