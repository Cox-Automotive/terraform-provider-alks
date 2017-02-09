package = github.com/Cox-Automotive/terraform-provider-alks

# default: build plan

get-deps:
	go install github.com/hashicorp/terraform

build:
	go build -v -o terraform-provider-alks .

test:
	go test -v .

plan:
	@terraform plan

.PHONY: install release test

install:
	go get -t -v ./...

release:
	mkdir -p release
	GOOS=linux GOARCH=amd64 go build -o release/terraform-provider-alks-linux-amd64 $(package)
	GOOS=linux GOARCH=386 go build -o release/terraform-provider-alks-linux-386 $(package)
	GOOS=linux GOARCH=arm go build -o release/terraform-provider-alks-linux-arm $(package)