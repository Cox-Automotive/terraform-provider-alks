package = github.com/Cox-Automotive/terraform-provider-alks

get-deps:
	go install github.com/hashicorp/terraform
	go get github.com/hashicorp/terraform
	go install github.com/hashicorp/go-cleanhttp
	go get github.com/hashicorp/go-cleanhttp
	go install github.com/Cox-Automotive/alks-go
	go get github.com/Cox-Automotive/alks-go

format:
	go fmt

build:
	go fmt
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
	GOOS=darwin  GOARCH=386   go build -o release/terraform-provider-alks-darwin-386 $(package)
	GOOS=darwin  GOARCH=amd64 go build -o release/terraform-provider-alks-darwin-amd64 $(package)
	GOOS=freebsd GOARCH=386   go build -o release/terraform-provider-alks-freebsd-386 $(package)
	GOOS=freebsd GOARCH=amd64 go build -o release/terraform-provider-alks-freebsd-amd64 $(package)
	GOOS=linux   GOARCH=386   go build -o release/terraform-provider-alks-linux-386 $(package)
	GOOS=linux   GOARCH=amd64 go build -o release/terraform-provider-alks-linux-amd64 $(package)
	GOOS=solaris GOARCH=amd64 go build -o release/terraform-provider-alks-solaris-amd64 $(package)
	GOOS=windows GOARCH=386   go build -o release/terraform-provider-alks-windows-386 $(package)
	GOOS=windows GOARCH=amd64 go build -o release/terraform-provider-alks-windows-amd64 $(package)