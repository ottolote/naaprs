VERSION := $(shell git describe --tags --always --dirty)
LDFLAGS := -X main.version=$(VERSION)

GOCMD = go
GOBUILD = $(GOCMD) build
GOCLEAN = $(GOCMD) clean
GOTEST = $(GOCMD) test
GOGET = $(GOCMD) get

CMD_DIR = cmd
BIN_DIR = bin
BINARY_NAME = naaprs
DOCKER_REGISTRY = ottokl
DOCKER_REPO = $(DOCKER_REGISTRY)/$(BINARY_NAME)
DOCKER_IMAGE = $(DOCKER_REPO):$(VERSION)

build: clean $(BIN_DIR)/$(BINARY_NAME)

$(BIN_DIR)/$(BINARY_NAME):
	$(GOBUILD) \
	     -ldflags "$(LDFLAGS)" \
	     -o $(BIN_DIR)/$(BINARY_NAME) \
	     $(CMD_DIR)/naaprs/main.go

push: docker-build docker-push docker-push-multi-arch

push-latest: push docker-push-multi-arch-latest

docker-build: docker-build-amd64 docker-build-386 docker-build-arm32v7 docker-build-arm64v8

docker-build-amd64:
	docker buildx build --platform linux/amd64 -t $(DOCKER_IMAGE)-amd64 --load .

docker-build-386:
	docker buildx build --platform linux/386 -t $(DOCKER_IMAGE)-386 --load .

docker-build-arm32v7:
	docker buildx build --platform linux/arm/v7 -t $(DOCKER_IMAGE)-arm32v7 --load .

docker-build-arm64v8:
	docker buildx build --platform linux/arm64/v8 -t $(DOCKER_IMAGE)-arm64v8 --load .

docker-push: docker-push-amd64 docker-push-386 docker-push-arm32v7 docker-push-arm64v8

docker-push-amd64:
	docker push $(DOCKER_IMAGE)-amd64

docker-push-386:
	docker push $(DOCKER_IMAGE)-386

docker-push-arm32v7:
	docker push $(DOCKER_IMAGE)-arm32v7

docker-push-arm64v8:
	docker push $(DOCKER_IMAGE)-arm64v8

docker-push-multi-arch: docker-manifest-create docker-manifest-push

docker-push-multi-arch-latest: docker-manifest-create-latest docker-manifest-push-latest

docker-manifest-create:
	docker manifest create $(DOCKER_REPO):$(VERSION) \
		--amend $(DOCKER_IMAGE)-amd64 \
		--amend $(DOCKER_IMAGE)-386 \
		--amend $(DOCKER_IMAGE)-arm32v7 \
		--amend $(DOCKER_IMAGE)-arm64v8

docker-manifest-create-latest:
	docker manifest create $(DOCKER_REPO):latest \
		--amend $(DOCKER_IMAGE)-amd64 \
		--amend $(DOCKER_IMAGE)-386 \
		--amend $(DOCKER_IMAGE)-arm32v7 \
		--amend $(DOCKER_IMAGE)-arm64v8

docker-manifest-push:
	docker manifest push $(DOCKER_REPO):$(VERSION)

docker-manifest-push-latest:
	docker manifest push $(DOCKER_REPO):latest

clean:
	$(GOCLEAN)
	rm -f $(BIN_DIR)/*

test:
	$(GOTEST) ./...

deps:
	$(GOGET) ./...

run:
	./$(BIN_DIR)/$(BINARY_NAME)

.PHONY: build clean test deps run docker-build docker-push docker-build-amd64 docker-build-386 docker-build-arm32v7 docker-build-arm64v8 docker-push-amd64 docker-push-386 docker-push-arm32v7 docker-push-arm64v8 docker-push-multi-arch docker-manifest-create docker-manifest-push

