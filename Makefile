all: fmt test build

deps:
	go get .

fmt:
	gofmt -w .

test:
	exit 0
	#go test

build:
	go build -o bin/`basename ${PWD}`
