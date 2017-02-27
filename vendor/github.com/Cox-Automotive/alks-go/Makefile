package = github.com/Cox-Automotive/alks-go

build:
	go fmt
	go build -v .

test:
	go test -v .

get-deps:
	go get github.com/hashicorp/go-cleanhttp
	go install github.com/hashicorp/go-cleanhttp
	go get github.com/motain/gocheck
	go install github.com/motain/gocheck

format:
	go fmt