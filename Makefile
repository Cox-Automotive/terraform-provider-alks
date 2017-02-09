default: build plan

get-deps:
	go install github.com/hashicorp/terraform

build:
	go build -v -o terraform-provider-alks .

test:
	go test -v .

plan:
	@terraform plan

.PHONY: install release test travis

install:
	go get -t -v ./...

travis:
	$(HOME)/gopath/bin/goveralls -service=travis-ci

release:
	mkdir -p release
	GOOS=linux GOARCH=amd64 go build -o release/terraform-provider-alks-linux-amd64 $(package)
	GOOS=linux GOARCH=386 go build -o release/terraform-provider-alks-linux-386 $(package)
	GOOS=linux GOARCH=arm go build -o release/terraform-provider-alks-linux-arm $(package)