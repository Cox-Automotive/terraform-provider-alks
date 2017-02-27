GOFMT_FILES?=$$(find . -name '*.go' | grep -v vendor)

format:
	gofmt -w $(GOFMT_FILES)

build:
	go fmt
	go build -v -o terraform-provider-alks

test:
	go test -v .

plan:
	@terraform plan

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
	GOOS=windows GOARCH=386   go build -o release/terraform-provider-alks-windows-386.exe $(package)
	GOOS=windows GOARCH=amd64 go build -o release/terraform-provider-alks-windows-amd64.exe $(package)