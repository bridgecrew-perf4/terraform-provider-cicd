all: build install

test:
	TF_LOG=WARN TF_ACC=true go test -v ./...

build:
	mkdir -p ./bin
	go build -o ./bin/terraform-provider-cicd .

install:
	mkdir -p ~/.terraform.d/plugins
	cp -f ./bin/terraform-provider-cicd ~/.terraform.d/plugins/terraform-provider-cicd
	# cp -f ./bin/terraform-provider-cicd ../tf_shell/bin/terraform-provider-cicd
	# cp -f ./bin/terraform-provider-cicd ../pipelines-server/bin/terraform-provider-cicd