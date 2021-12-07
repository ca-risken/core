TARGETS = alert finding iam project report
BUILD_TARGETS = $(TARGETS:=.build)
BUILD_CI_TARGETS = $(TARGETS:=.build-ci)
IMAGE_PUSH_TARGETS = $(TARGETS:=.push-image)
MANIFEST_CREATE_TARGETS = $(TARGETS:=.create-manifest)
MANIFEST_PUSH_TARGETS = $(TARGETS:=.push-manifest)
TEST_TARGETS = $(TARGETS:=.go-test)
LINT_TARGETS = $(TARGETS:=.lint)
MOCK_TARGETS = $(TARGETS:=.mock)
BUILD_OPT=""
IMAGE_TAG=latest
MANIFEST_TAG=latest
IMAGE_PREFIX=core
IMAGE_REGISTRY=local

.PHONY: all
all: run

.PHONY: install
install:
	go get \
		google.golang.org/grpc \
		github.com/golang/protobuf/protoc-gen-go \
		github.com/envoyproxy/protoc-gen-validate \
		github.com/grpc-ecosystem/go-grpc-middleware

.PHONY: clean
clean:
	rm -f proto/*/*.pb.go
	rm -f doc/*.md

.PHONY: fmt
fmt: proto/**/*.proto
	@clang-format -i proto/**/*.proto

.PHONY: doc
doc: fmt
	protoc \
		--proto_path=proto \
		--error_format=gcc \
		-I $(GOPATH)/src/github.com/envoyproxy/protoc-gen-validate \
		--doc_out=markdown,README.md:doc \
		proto/**/*.proto;

# build without protoc-gen-validate
.PHONY: proto-without-validation
proto-without-validate: fmt
	for svc in "alert" "finding" "iam"; do \
		protoc \
			--proto_path=proto \
			--error_format=gcc \
			--go_out=plugins=grpc,paths=source_relative:proto \
			proto/$$svc/*.proto; \
	done

# build with protoc-gen-validate
.PHONY: proto-validate
proto-validate: fmt
	for svc in "project" "report"; do \
		protoc \
			--proto_path=proto \
			--error_format=gcc \
			-I $(GOPATH)/src/github.com/envoyproxy/protoc-gen-validate \
			--go_out=plugins=grpc,paths=source_relative:proto \
			--validate_out="lang=go,paths=source_relative:proto" \
			proto/$$svc/*.proto; \
	done

.PHONY: proto
proto : proto-without-validate proto-validate proto-mock

PHONY: build $(BUILD_TARGETS)
build: $(BUILD_TARGETS)
%.build: %.go-test
	. env.sh && TARGET=$(*) IMAGE_TAG=$(IMAGE_TAG) IMAGE_PREFIX=$(IMAGE_PREFIX) BUILD_OPT="$(BUILD_OPT)" . hack/docker-build.sh

PHONY: build-ci $(BUILD_CI_TARGETS)
build-ci: $(BUILD_CI_TARGETS)
%.build-ci:
	TARGET=$(*) IMAGE_TAG=$(IMAGE_TAG) IMAGE_PREFIX=$(IMAGE_PREFIX) BUILD_OPT="$(BUILD_OPT)" . hack/docker-build.sh
	docker tag $(IMAGE_PREFIX)/$(*):$(IMAGE_TAG) $(IMAGE_REGISTRY)/$(IMAGE_PREFIX)/$(*):$(IMAGE_TAG)

PHONY: push-image $(IMAGE_PUSH_TARGETS)
push-image: $(IMAGE_PUSH_TARGETS)
%.push-image:
	docker push $(IMAGE_REGISTRY)/$(IMAGE_PREFIX)/$(*):$(IMAGE_TAG)

PHONY: create-manifest $(MANIFEST_CREATE_TARGETS)
create-manifest: $(MANIFEST_CREATE_TARGETS)
%.create-manifest:
	docker manifest create $(IMAGE_REGISTRY)/$(IMAGE_PREFIX)/$(*):$(MANIFEST_TAG) $(IMAGE_REGISTRY)/$(IMAGE_PREFIX)/$(*):$(IMAGE_TAG_BASE)_linux_amd64 $(IMAGE_REGISTRY)/$(IMAGE_PREFIX)/$(*):$(IMAGE_TAG_BASE)_linux_arm64
	docker manifest annotate --arch amd64 $(IMAGE_REGISTRY)/$(IMAGE_PREFIX)/$(*):$(MANIFEST_TAG) $(IMAGE_REGISTRY)/$(IMAGE_PREFIX)/$(*):$(IMAGE_TAG_BASE)_linux_amd64
	docker manifest annotate --arch arm64 $(IMAGE_REGISTRY)/$(IMAGE_PREFIX)/$(*):$(MANIFEST_TAG) $(IMAGE_REGISTRY)/$(IMAGE_PREFIX)/$(*):$(IMAGE_TAG_BASE)_linux_arm64

PHONY: push-manifest $(MANIFEST_PUSH_TARGETS)
push-manifest: $(MANIFEST_PUSH_TARGETS)
%.push-manifest:
	docker manifest push $(IMAGE_REGISTRY)/$(IMAGE_PREFIX)/$(*):$(MANIFEST_TAG)
	docker manifest inspect $(IMAGE_REGISTRY)/$(IMAGE_PREFIX)/$(*):$(MANIFEST_TAG)

PHONY: go-test $(TEST_TARGETS)
go-test: $(TEST_TARGETS)
%.go-test:
	cd src/$(*) && AWS_XRAY_SDK_DISABLED=TRUE go test ./...
	cd proto/$(*) && AWS_XRAY_SDK_DISABLED=TRUE go test ./...


.PHONY: go-mod-tidy
go-mod-tidy: proto
	source env.sh && cd proto/finding && go mod tidy
	source env.sh && cd proto/iam     && go mod tidy
	source env.sh && cd proto/project && go mod tidy
	source env.sh && cd proto/alert   && go mod tidy
	source env.sh && cd proto/report  && go mod tidy
	source env.sh && cd src/finding   && go mod tidy
	source env.sh && cd src/iam       && go mod tidy
	source env.sh && cd src/project   && go mod tidy
	source env.sh && cd src/alert     && go mod tidy
	source env.sh && cd src/report    && go mod tidy

.PHONY: go-mod-update
go-mod-update:
	source env.sh \
		&& cd src/finding \
		&& go get -u github.com/ca-risken/core/...
	source env.sh \
		cd src/iam \
		&& go get -u github.com/ca-risken/core/...
	source env.sh \
		cd src/project \
		&& go get -u github.com/ca-risken/core/...
	source env.sh \
		cd src/alert \
		&& go get -u github.com/ca-risken/core/...
	source env.sh \
		cd src/report \
		&& go get -u github.com/ca-risken/core/...

.PHONY: lint proto-lint pkg-lint
lint: $(LINT_TARGETS) proto-lint pkg-lint
%.lint: FAKE
	sh hack/golinter.sh src/$(*)
proto-lint:
	sh hack/golinter.sh proto/alert
	sh hack/golinter.sh proto/finding
	sh hack/golinter.sh proto/iam
	sh hack/golinter.sh proto/project
	sh hack/golinter.sh proto/report
pkg-lint:
	sh hack/golinter.sh pkg/model

.PHONY: generate-mock
generate-mock: proto-mock
proto-mock: $(MOCK_TARGETS)
%.mock: FAKE
	sh hack/generate-mock.sh proto/$(*)

FAKE: