GO111MODULE=on
package=github.com/Cox-Automotive/terraform-provider-alks
GOFMT_FILES?=$$(find . -name '*.go' | grep -v vendor)
RELEASE_DIR=release

format:
	gofmt -w $(GOFMT_FILES)

clean:
	go clean
	rm -rf $(RELEASE_DIR)

deps:
	go mod tidy
	go get -u
	go mod vendor

build: clean deps
	go fmt
	go build -v -o examples/terraform-provider-alks .

test:
	go test -v .

plan:
	@terraform plan

install:
	go get -t -v ./...

release: clean deps
	mkdir -p $(RELEASE_DIR)

	GOOS=darwin GOARCH=amd64 go build -o $(RELEASE_DIR)/terraform-provider-alks_v$(TRAVIS_TAG) $(package)
	chmod +x $(RELEASE_DIR)/terraform-provider-alks_v$(TRAVIS_TAG)
	tar -cvzf $(RELEASE_DIR)/terraform-provider-alks-darwin-amd64.tar.gz -C $(RELEASE_DIR)/ terraform-provider-alks_v$(TRAVIS_TAG)

	GOOS=freebsd GOARCH=386 go build -o $(RELEASE_DIR)/terraform-provider-alks_v$(TRAVIS_TAG) $(package)
	chmod +x $(RELEASE_DIR)/terraform-provider-alks_v$(TRAVIS_TAG)
	tar -cvzf $(RELEASE_DIR)/terraform-provider-alks-freebsd-386.tar.gz -C $(RELEASE_DIR)/ terraform-provider-alks_v$(TRAVIS_TAG)

	GOOS=freebsd GOARCH=amd64 go build -o $(RELEASE_DIR)/terraform-provider-alks_v$(TRAVIS_TAG) $(package)
	chmod +x $(RELEASE_DIR)/terraform-provider-alks_v$(TRAVIS_TAG)
	tar -cvzf $(RELEASE_DIR)/terraform-provider-alks-freebsd-amd64.tar.gz -C $(RELEASE_DIR)/ terraform-provider-alks_v$(TRAVIS_TAG)

	GOOS=linux GOARCH=386 go build -o $(RELEASE_DIR)/terraform-provider-alks_v$(TRAVIS_TAG) $(package)
	chmod +x $(RELEASE_DIR)/terraform-provider-alks_v$(TRAVIS_TAG)
	tar -cvzf $(RELEASE_DIR)/terraform-provider-alks-linux-386.tar.gz -C $(RELEASE_DIR)/ terraform-provider-alks_v$(TRAVIS_TAG)

	GOOS=linux GOARCH=amd64 go build -o $(RELEASE_DIR)/terraform-provider-alks_v$(TRAVIS_TAG) $(package)
	chmod +x $(RELEASE_DIR)/terraform-provider-alks_v$(TRAVIS_TAG)
	tar -cvzf $(RELEASE_DIR)/terraform-provider-alks-linux-amd64.tar.gz -C $(RELEASE_DIR)/ terraform-provider-alks_v$(TRAVIS_TAG)

	GOOS=solaris GOARCH=amd64 go build -o $(RELEASE_DIR)/terraform-provider-alks_v$(TRAVIS_TAG) $(package)
	chmod +x $(RELEASE_DIR)/terraform-provider-alks_v$(TRAVIS_TAG)
	tar -cvzf $(RELEASE_DIR)/terraform-provider-alks-solaris-amd64.tar.gz -C $(RELEASE_DIR)/ terraform-provider-alks_v$(TRAVIS_TAG)

	GOOS=windows GOARCH=386 go build -o $(RELEASE_DIR)/terraform-provider-alks_v$(TRAVIS_TAG).exe $(package)
	zip $(RELEASE_DIR)/terraform-provider-alks-windows-386.zip $(RELEASE_DIR)/terraform-provider-alks_v$(TRAVIS_TAG).exe

	GOOS=windows GOARCH=amd64 go build -o $(RELEASE_DIR)/terraform-provider-alks_v$(TRAVIS_TAG).exe $(package)
	zip $(RELEASE_DIR)/terraform-provider-alks-windows-amd64.zip $(RELEASE_DIR)/terraform-provider-alks_v$(TRAVIS_TAG).exe

	rm $(RELEASE_DIR)/terraform-provider-alks_v$(TRAVIS_TAG).exe
	rm $(RELEASE_DIR)/terraform-provider-alks_v$(TRAVIS_TAG)

