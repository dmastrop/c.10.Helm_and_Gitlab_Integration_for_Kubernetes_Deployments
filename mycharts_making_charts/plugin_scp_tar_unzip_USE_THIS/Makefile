BINARY_NAME=helmscp

.PHONY: build dep

build: dep
	@go build -o bin/${BINARY_NAME}

all:
	@go mod download
	@GOARCH=amd64 GOOS=darwin go build -o bin/${BINARY_NAME}-darwin main.go
	@GOARCH=amd64 GOOS=linux go build -o bin/${BINARY_NAME}-linux main.go
	@GOARCH=amd64 GOOS=windows go build -o bin/${BINARY_NAME}-windows main.go

clean:
	go clean
	rm ${BINARY_NAME}-darwin
	rm ${BINARY_NAME}-linux
	rm ${BINARY_NAME}-windows

dep:
	@go mod download