# Copyright 2017-2021 Digital Asset Exchange Limited. All rights reserved.
# Use of this source code is governed by Microsoft Reference Source
# License (MS-RSL) that can be found in the LICENSE file.

all: build install

test:
	TF_LOG=WARN TF_ACC=true go test -v ./...

build:
	mkdir -p ./bin
	go build -o ./bin/terraform-provider-cicd .

upload-local: build
	docker cp bin/terraform-provider-cicd tf-shell-atlantrentals-dev2:/root/.terraform.d/plugins/linux_amd64/terraform-provider-cicd_v0.0.7

install:
	mkdir -p ~/.terraform.d/plugins
	cp -f ./bin/terraform-provider-cicd ~/.terraform.d/plugins/terraform-provider-cicd
	# cp -f ./bin/terraform-provider-cicd ../tf_shell/bin/terraform-provider-cicd
	# cp -f ./bin/terraform-provider-cicd ../pipelines-server/bin/terraform-provider-cicd