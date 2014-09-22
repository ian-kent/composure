all: deps build

fmt:
	go fmt ./...

deps:
	-go get github.com/smartystreets/goconvey
	-go get github.com/gorilla/pat

test:
	go test ./...

build: fmt
	go install .

.PHONY: all fmt deps test build
