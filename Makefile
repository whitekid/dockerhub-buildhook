TARGET=bin/dockerhub-buildhook
SRC=$(shell find . -type f -name '*.go' -not -path "./vendor/*" -not -path "*_test.go")
BUILD_FLAGS?=-v
DOCKER_TAG?=dockerhub-buildhook

.PHONY: clean test get tidy

all: build
build: $(TARGET)

$(TARGET): $(SRC)
	go build -o $(TARGET) ${BUILD_FLAGS} .

clean:
	rm -f ${TARGET}

test:
	go test

# update modules & tidy
get:
	go get -u
	@$(MAKE) tidy

tidy:
	go mod tidy

# build docker image
docker:
	docker build -t $(DOCKER_TAG) -f Dockerfile .

freebsd:
	@GOOS=freebsd GOARCH=386 CGO_ENABLED=0 ${MAKE} build
