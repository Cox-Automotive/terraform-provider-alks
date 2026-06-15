GOFMT_FILES?=$$(find . -name '*.go' | grep -v vendor)

format:
	gofmt -w $(GOFMT_FILES)

build:
	go fmt
	go build -v -o examples/terraform-provider-alks -mod=vendor .

test:
	go test -timeout 1200s -v .

testacc:
	@echo "set TESTARGS=\"-run TestAccXXX\" to run individual tests"
	TF_ACC=1 go test -timeout 1200s -v . $(TESTARGS)

plan:
	@terraform plan

install:
	go get -t -v ./...
