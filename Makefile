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

PLATFORMS = linux/amd64 linux/386 linux/arm/v6 linux/arm/v7 linux/arm64/v8

build: clean $(BIN_DIR)/$(BINARY_NAME)

$(BIN_DIR)/$(BINARY_NAME):
	$(GOBUILD) -ldflags "$(LDFLAGS)" -o $@ $(CMD_DIR)/naaprs/main.go

push: docker-build docker-push docker-push-multi-arch

push-latest: docker-build docker-push docker-push-multi-arch-latest

docker-build:
	$(foreach platform,$(PLATFORMS),\
	docker buildx build \
	--build-arg LDFLAGS="$(LDFLAGS)" \
	--build-arg GOOS=$(word 1,$(subst /, ,$(platform))) \
	--build-arg GOARCH=$(word 2,$(subst /, ,$(platform))) \
	--build-arg GOARM=$(if $(findstring arm/,$(platform)),$(subst v,,$(word 3,$(subst /, ,$(platform)))),) \
	--platform $(platform) -t $(DOCKER_IMAGE)-$(subst /,-,$(platform)) --load .;)

docker-push:
	$(foreach platform,$(PLATFORMS),docker push $(DOCKER_IMAGE)-$(subst /,-,$(platform));)

docker-push-multi-arch:
	docker manifest create $(DOCKER_IMAGE) $(foreach platform,$(PLATFORMS),--amend $(DOCKER_IMAGE)-$(subst /,-,$(platform)))
	docker manifest push $(DOCKER_IMAGE)

docker-push-multi-arch-latest:
	docker manifest create $(DOCKER_REPO):latest $(foreach platform,$(PLATFORMS),--amend $(DOCKER_IMAGE)-$(subst /,-,$(platform)))
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

.PHONY: build clean test deps run push push-latest docker-build docker-push docker-push-multi-arch docker-push-multi-arch-latest

