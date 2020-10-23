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
	tar -cvzf release/terraform-provider-alks_$(TRAVIS_TAG)_darwin_amd64.tar.gz -C release/ terraform-provider-alks_v$(TRAVIS_TAG)

	GOOS=freebsd GOARCH=386 go build -ldflags "-X main.versionNumber=$(TRAVIS_TAG)" -o release/terraform-provider-alks_v$(TRAVIS_TAG) -mod=vendor $(package)
	chmod +x release/terraform-provider-alks_v$(TRAVIS_TAG)
	tar -cvzf release/terraform-provider-alks_$(TRAVIS_TAG)_freebsd_386.tar.gz -C release/ terraform-provider-alks_v$(TRAVIS_TAG)

	GOOS=freebsd GOARCH=amd64 go build -ldflags "-X main.versionNumber=$(TRAVIS_TAG)" -o release/terraform-provider-alks_v$(TRAVIS_TAG) -mod=vendor $(package)
	chmod +x release/terraform-provider-alks_v$(TRAVIS_TAG)
	tar -cvzf release/terraform-provider-alks_$(TRAVIS_TAG)_freebsd_amd64.tar.gz -C release/ terraform-provider-alks_v$(TRAVIS_TAG)

	GOOS=linux GOARCH=386 go build -ldflags "-X main.versionNumber=$(TRAVIS_TAG)" -o release/terraform-provider-alks_v$(TRAVIS_TAG) -mod=vendor $(package)
	chmod +x release/terraform-provider-alks_v$(TRAVIS_TAG)
	tar -cvzf release/terraform-provider-alks_$(TRAVIS_TAG)_linux_386.tar.gz -C release/ terraform-provider-alks_v$(TRAVIS_TAG)

	GOOS=linux GOARCH=amd64 go build -ldflags "-X main.versionNumber=$(TRAVIS_TAG)" -o release/terraform-provider-alks_v$(TRAVIS_TAG) -mod=vendor $(package)
	chmod +x release/terraform-provider-alks_v$(TRAVIS_TAG)
	tar -cvzf release/terraform-provider-alks_$(TRAVIS_TAG)_linux_amd64.tar.gz -C release/ terraform-provider-alks_v$(TRAVIS_TAG)

	GOOS=solaris GOARCH=amd64 go build -ldflags "-X main.versionNumber=$(TRAVIS_TAG)" -o release/terraform-provider-alks_v$(TRAVIS_TAG) -mod=vendor $(package)
	chmod +x release/terraform-provider-alks_v$(TRAVIS_TAG)
	tar -cvzf release/terraform-provider-alks_$(TRAVIS_TAG)_solaris_amd64.tar.gz -C release/ terraform-provider-alks_v$(TRAVIS_TAG)

	GOOS=windows GOARCH=386 go build -ldflags "-X main.versionNumber=$(TRAVIS_TAG)" -o release/terraform-provider-alks_v$(TRAVIS_TAG).exe -mod=vendor $(package)
	zip release/terraform-provider-alks_$(TRAVIS_TAG)_windows_386.zip release/terraform-provider-alks_v$(TRAVIS_TAG).exe

	GOOS=windows GOARCH=amd64 go build -ldflags "-X main.versionNumber=$(TRAVIS_TAG)" -o release/terraform-provider-alks_v$(TRAVIS_TAG).exe -mod=vendor $(package)
	zip release/terraform-provider-alks_$(TRAVIS_TAG)_windows_amd64.zip release/terraform-provider-alks_v$(TRAVIS_TAG).exe

	shasum -a 256 release/*.tar.gz release/*.zip > release/terraform-provider-alks_$(TRAVIS_TAG)_SHA256SUMS

	echo $(GPG_KEY) | base64 --decode | gpg --batch --no-tty --yes --import

	@gpg --pinentry-mode loopback --passphrase $(GPG_PASSPHRASE) -u C182B91A3A62B0D5 --detach-sign release/terraform-provider-alks_$(TRAVIS_TAG)_SHA256SUMS

	rm release/terraform-provider-alks_v$(TRAVIS_TAG).exe

