BINARY_NAME=snark
BUILD_TIME=$(shell date)
VERSION=""

build:
	GOARCH=amd64 GOOS=darwin go build -ldflags="-X 'main.Version=${VERSION}' -X 'main.BuildTime=${BUILD_TIME}'" -o ${BINARY_NAME} main.go
	tar -cvf "${BINARY_NAME}-${VERSION}-darwin-amd64.tar" ./${BINARY_NAME}
	gzip ./${BINARY_NAME}-${VERSION}-darwin-amd64.tar
	rm ${BINARY_NAME}

	GOARCH=amd64 GOOS=linux go build -ldflags="-X 'main.Version=${VERSION}' -X 'main.BuildTime=${BUILD_TIME}'" -o ${BINARY_NAME} main.go
	tar -cvf "${BINARY_NAME}-${VERSION}-linux-amd64.tar" ./${BINARY_NAME}
	gzip ./${BINARY_NAME}-${VERSION}-linux-amd64.tar
	rm ${BINARY_NAME}

dep:
	go mod download

vet:
	go vet

lint:
	golangci-lint run --enable-all

clean:
	go clean
	rm ${BINARY_NAME}*