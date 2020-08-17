package=github.com/Cox-Automotive/terraform-provider-alks
GOFMT_FILES?=$$(find . -name '*.go' | grep -v vendor)

format:
	gofmt -w $(GOFMT_FILES)

build:
	go fmt
	go build -v -o examples/terraform-provider-alks -mod=vendor .

test:
	go test -v .

plan:
	@terraform plan

install:
	go get -t -v ./...

release:
	mkdir -p release

	GOOS=darwin GOARCH=amd64 go build -ldflags "-X main.versionNumber=$(TRAVIS_TAG)" -o release/terraform-provider-alks_v$(TRAVIS_TAG) -mod=vendor $(package)
	chmod +x release/terraform-provider-alks_v$(TRAVIS_TAG)
	tar -cvzf release/terraform-provider-alks-darwin-amd64.tar.gz -C release/ terraform-provider-alks_v$(TRAVIS_TAG)

	GOOS=freebsd GOARCH=386 go build -ldflags "-X main.versionNumber=$(TRAVIS_TAG)" -o release/terraform-provider-alks_v$(TRAVIS_TAG) -mod=vendor $(package)
	chmod +x release/terraform-provider-alks_v$(TRAVIS_TAG)
	tar -cvzf release/terraform-provider-alks-freebsd-386.tar.gz -C release/ terraform-provider-alks_v$(TRAVIS_TAG)

	GOOS=freebsd GOARCH=amd64 go build -ldflags "-X main.versionNumber=$(TRAVIS_TAG)" -o release/terraform-provider-alks_v$(TRAVIS_TAG) -mod=vendor $(package)
	chmod +x release/terraform-provider-alks_v$(TRAVIS_TAG)
	tar -cvzf release/terraform-provider-alks-freebsd-amd64.tar.gz -C release/ terraform-provider-alks_v$(TRAVIS_TAG)

	GOOS=linux GOARCH=386 go build -ldflags "-X main.versionNumber=$(TRAVIS_TAG)" -o release/terraform-provider-alks_v$(TRAVIS_TAG) -mod=vendor $(package)
	chmod +x release/terraform-provider-alks_v$(TRAVIS_TAG)
	tar -cvzf release/terraform-provider-alks-linux-386.tar.gz -C release/ terraform-provider-alks_v$(TRAVIS_TAG)

	GOOS=linux GOARCH=amd64 go build -ldflags "-X main.versionNumber=$(TRAVIS_TAG)" -o release/terraform-provider-alks_v$(TRAVIS_TAG) -mod=vendor $(package)
	chmod +x release/terraform-provider-alks_v$(TRAVIS_TAG)
	tar -cvzf release/terraform-provider-alks-linux-amd64.tar.gz -C release/ terraform-provider-alks_v$(TRAVIS_TAG)

	GOOS=solaris GOARCH=amd64 go build -ldflags "-X main.versionNumber=$(TRAVIS_TAG)" -o release/terraform-provider-alks_v$(TRAVIS_TAG) -mod=vendor $(package)
	chmod +x release/terraform-provider-alks_v$(TRAVIS_TAG)
	tar -cvzf release/terraform-provider-alks-solaris-amd64.tar.gz -C release/ terraform-provider-alks_v$(TRAVIS_TAG)

	GOOS=windows GOARCH=386 go build -ldflags "-X main.versionNumber=$(TRAVIS_TAG)" -o release/terraform-provider-alks_v$(TRAVIS_TAG).exe -mod=vendor $(package)
	zip release/terraform-provider-alks-windows-386.zip release/terraform-provider-alks_v$(TRAVIS_TAG).exe

	GOOS=windows GOARCH=amd64 go build -ldflags "-X main.versionNumber=$(TRAVIS_TAG)" -o release/terraform-provider-alks_v$(TRAVIS_TAG).exe -mod=vendor $(package)
	zip release/terraform-provider-alks-windows-amd64.zip release/terraform-provider-alks_v$(TRAVIS_TAG).exe

	shasum -a 256 release/*.zip > release/terraform-provider-alks_v$(TRAVIS_TAG)_SHA256SUMS

	rm release/terraform-provider-alks_v$(TRAVIS_TAG).exe
	rm release/terraform-provider-alks_v$(TRAVIS_TAG)

