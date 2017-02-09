default: build plan

get-deps:
	go install github.com/hashicorp/terraform

build:
	go build -v -o terraform-provider-alks .

test:
	go test -v .

plan:
	@terraform plan
