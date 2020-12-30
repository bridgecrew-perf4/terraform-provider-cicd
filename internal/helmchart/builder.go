// Copyright 2017-2020 Tensigma Ltd. All rights reserved.
// Use of this source code is governed by Microsoft Reference Source
// License (MS-RSL) that can be found in the LICENSE file.

package helmchart

import (
	"archive/zip"
	"bytes"

	"fmt"
	"github.com/AtlantPlatform/infra/terraform-provider-cicd/internal/helpers"
	"golang.org/x/mod/sumdb/dirhash"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Declaration struct {
	ApiVersion  string `yaml:"apiVersion"`
	AppVersion  string `yaml:"appVersion"`
	Description string `yaml:"description"`
	Name        string `yaml:"name"`
	Version     string `yaml:"version"`
}

type Builder struct {
	ID   string
	Hash string
	Name string

	source string
	yamlChart  string
	yamlValues string
	txtOverride string
}

func New(source string, args map[string]interface{}) (*Builder, error) {
	if _, err := os.Stat(source + "/values.yaml"); os.IsNotExist(err) {
		return nil, fmt.Errorf("%s/values.yaml file not found", source)
	} else if _, err := os.Stat(source + "/Chart.yaml"); os.IsNotExist(err) {
		return nil, fmt.Errorf("%s/Chart.yaml file not found", source)
	} else if _, err := os.Stat(source + "/templates"); os.IsNotExist(err) {
		return nil, fmt.Errorf("%s/templates folder not found", source)
	}

	yamlChartFile, err := ioutil.ReadFile(fmt.Sprintf("%s/Chart.yaml", source))
	if err != nil {
		return nil, fmt.Errorf("%s/Chart.yaml read failure %v", source, err)
	}
	var decl *Declaration
	if err := yaml.Unmarshal(yamlChartFile, &decl); err != nil {
		return nil, fmt.Errorf("%s/Chart.yaml parse failure %v", source, err)
	}
	
	yamlValuesFile, err := ioutil.ReadFile(fmt.Sprintf("%s/values.yaml", source))
	if err != nil {
		return nil, fmt.Errorf("%s/values.yaml read failure %v", source, err)
	}
	var t interface{}
	if err := yaml.Unmarshal(yamlValuesFile, &t); err != nil {
	 	return nil, fmt.Errorf("%s/values.yaml parse failure %v", source, err)
	}
	arrOverride := make([]string, 0)
	if args != nil {
		for k, v := range args {
			arrOverride = append(arrOverride, fmt.Sprintf("--set '%s'='%s'", k, v))
		}
	}
	hash, err := dirhash.HashDir(source, "", dirhash.Hash1)
	if err != nil {
		return nil, fmt.Errorf("source %v hash failure %v", source, err)
	}
	id := hash[0:12]
	return &Builder{
		Name: decl.Name,
		ID: id,
		Hash: hash,
		source: source,
		yamlChart: string(yamlChartFile),
		yamlValues: string(yamlValuesFile),
		txtOverride: strings.Join(arrOverride, " "),
	}, nil
}

func (s *Builder) GetZipName() string {
	if s.ID == "" || s.Name == "" {
		return ""
	}
	return "helm/" + s.Name + "-" + s.ID + ".zip"
}

type zipFile struct {
	Name string
	Body string
}

func (s *Builder) GetHash() string {
	return s.Hash
}

func (s *Builder) GetID() string {
	return s.Hash
}

func (s *Builder) ZIP() (io.ReadSeeker, error) {
	// Create a buffer to write our archive to.
	buf := new(bytes.Buffer)
	// Create a new zip archive.
	w := zip.NewWriter(buf)
	// Add some files to the archive.
	var files = []zipFile{
		{"values.yaml", s.yamlValues},
		{"override.txt", s.txtOverride},
		{"Chart.yaml", s.yamlChart},
	}
	// read files under templates/ folder
	err := filepath.Walk(s.source + "/templates", func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		fmt.Printf("reading %v\n", info.Name())
		yamlFile, err := ioutil.ReadFile(path)
		if err != nil {
			return fmt.Errorf("%s/%s read failure %v", s.source, path, err)
		}
        files = append(files, zipFile{ 
			Name: "templates/" + info.Name(), 
			Body: string(yamlFile),
		})
        return nil
    })
    if err != nil {
        panic(err)
    }
   
	for _, file := range files {
		f, err := w.Create(file.Name)
		if err != nil {
			return nil, err
		}
		_, err = f.Write([]byte(file.Body))
		if err != nil {
			return nil, err
		}
	}

	// Make sure to check the error on Close.
	if errClose := w.Close(); errClose != nil {
		return nil, err
	}
	return helpers.NewReadSeeker(buf.Bytes()), nil
}
