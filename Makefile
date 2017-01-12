all: fmt deps build

deps:
	go get .

fmt:
	gofmt -w ./..

build:
	go build -o bin/`basename ${PWD}`
