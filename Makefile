SOURCE_DIR=./cmd/showdown/
BINARY_NAME=showdown

generate:
	go generate ${SOURCE_DIR}

build: generate
	go build -o ${BINARY_NAME} ${SOURCE_DIR}

clean:
	go clean
	if [ -d ./cmd/showdown/assets/ ]; then rm -r ./cmd/showdown/assets/; fi
	if [ -f ${BINARY_NAME} ]; then rm ${BINARY_NAME}; fi

vet:
	go vet ${SOURCE_DIR}

lint:
	golangci-lint run --enable-all ${SOURCE_DIR}
