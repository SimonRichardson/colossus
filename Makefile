.PHONY: setup, test, unit-test, test-ci

project=colossus
compose=docker-compose -f docker/docker-compose.yml -p $(project)
compose-test=docker-compose -f docker/docker-tests.yml -p $(project) run --rm base-test

registry ?= hub.docker.com

TARGET ?= colossus
TAG ?= dev

setup:
	go get github.com/Masterminds/glide
	glide cc
	glide install
	go get -u github.com/jteeuwen/go-bindata/...
	make build-bindata

test:
	make TARGET=colossus test-target
	make TARGET=colossusw test-target
	make TARGET=colossuss test-target

test-target:
	$(compose-test) internal-${TARGET}-integration-tests
	$(compose-test) internal-${TARGET}-benchmark-tests

unit-test:
	$(compose-test) internal-${TARGET}-tests

test-ci: unit-test
	$(compose-test) internal-${TARGET}-integration-tests
	$(compose-test) internal-${TARGET}-benchmark-tests
	$(compose-test) internal-${TARGET}-perf-tests

.PHONY: build-bindata, build-schemas

build-bindata:
	cd scripts && go-bindata --nocompress -pkg scripts ./../scripts/... && cd ../

build-schemas:
	@rm -rf ./schemas/schema
	flatc -g -o ./schemas/ ./schemas/*fbs

.PHONY: build, build-image, build-images, publish-images

build:
	mkdir -p artifacts
	CGO_ENABLED=0 go build -o bin/colossus github.com/SimonRichardson/colossus/colossus-http
	tar -czf artifacts/colossus.tar.gz bin/colossus -C bin/ .
	CGO_ENABLED=0 go build -o bin/colossusw github.com/SimonRichardson/colossus/colossus-walker
	tar -czf artifacts/colossusw.tar.gz bin/colossusw -C bin/ .
	CGO_ENABLED=0 go build -o bin/colossuss github.com/SimonRichardson/colossus/colossus-shim
	tar -czf artifacts/colossuss.tar.gz bin/colossuss -C bin/ .

build-images: build-builder
	docker run --rm $(registry)/colossus-builder:${TAG} internal-colossus-build | docker build -f colossus-http/Dockerfile -t $(registry)/colossus-http:${TAG} -
	docker run --rm $(registry)/colossus-builder:${TAG} internal-colossusw-build | docker build -f colossus-walker/Dockerfile -t $(registry)/colossus-walker:${TAG} -
	docker run --rm $(registry)/colossus-builder:${TAG} internal-colossuss-build | docker build -f colossus-shim/Dockerfile -t $(registry)/colossus-shim:${TAG} -

build-builder:
	@docker build --build-arg GITHUB_TOKEN=${GITHUB_TOKEN} -f docker/Dockerfile.build --rm -t $(registry)/colossus-builder:${TAG} .

publish-releases:
	docker push $(registry)/colossus-http:${TAG}
	docker push $(registry)/colossus-walker:${TAG}
	docker push $(registry)/colossus-shim:${TAG}

.PHONY: internal-colossus-build, internal-colossusw-build, internal-colossuss-build

internal-colossus-build:
	@make setup -s
	CGO_ENABLED=0 go build -o bin/colossus github.com/SimonRichardson/colossus/colossus-http
	tar -czf - colossus-http/Dockerfile bin/colossus

internal-colossusw-build:
	@make setup -s
	CGO_ENABLED=0 go build -o bin/colossus github.com/SimonRichardson/colossus/colossus-walker
	tar -czf - colossus-walker/Dockerfile bin/colossusw

internal-colossuss-build:
	@make setup -s
	CGO_ENABLED=0 go build -o bin/colossus github.com/SimonRichardson/colossus/colossus-shim
	tar -czf - colossus-shim/Dockerfile bin/colossuss

.PHONY: internal-colossus-tests, internal-colossus-integration-tests, internal-colossus-benchmark-tests, internal-colossus-perf-tests

## colossus

internal-colossus-tests:
	go test -v ./cluster/counter/... -stubs=true
	go test -v ./cluster/counter/...
	go test -v ./cluster/store/... -stubs=true
	go test -v ./cluster/store/...
	go test -v ./selectors/...
	go test -v ./semaphore/...

internal-colossus-integration-tests:
	go test -v ./colossus-http/...

internal-colossus-benchmark-tests:
	go test -bench=. -run=BenchmarkRequest ./colossus-http/...

internal-colossus-perf-tests:
	sh -c 'go run colossus-perf/*.go -ciserver=true -clusterserver=true'

.PHONY: internal-colossusw-tests, internal-colossusw-integration-tests, internal-colossusw-benchmark-tests, internal-colossusw-perf-tests

## Walker

internal-colossusw-tests:
	# Do nothing!

internal-colossusw-integration-tests:
	go test -v ./colossus-walker/...

internal-colossusw-benchmark-tests:
	# Do nothing!

internal-colossusw-perf-tests:
	# Do nothing!

.PHONY: internal-colossuss-tests, internal-colossuss-integration-tests, internal-colossuss-benchmark-tests, internal-colossuss-perf-tests

## Shim

internal-colossuss-tests:
	# Do nothing!

internal-colossuss-integration-tests:
	CGO_ENABLED=0 go build -o bin/colossus github.com/SimonRichardson/colossus/colossus-http
	./bin/colossus &
	sleep 4
	go test -v ./colossus-shim/...

internal-colossuss-benchmark-tests:
	# Do nothing!

internal-colossuss-perf-tests:
	# Do nothing!

.PHONY: setup-env, teardown-env-dev, clean-env-dev

setup-env:
	$(compose) up -d

teardown-env-dev:
	$(compose) kill
	$(compose) rm -f

clean-env-dev:
	docker volume ls -qf "dangling=true" | xargs docker volume rm
	docker images -qf "dangling=true" | xargs docker rmi
