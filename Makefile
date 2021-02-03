.PHONY: all install clean network fmt build doc go-test go-mod-update go-mod-tidy run log stop
all: run

install:
	go get \
		google.golang.org/grpc \
		github.com/golang/protobuf/protoc-gen-go \
		github.com/grpc-ecosystem/go-grpc-middleware

clean:
	rm -f proto/*/*.pb.go
	rm -f doc/*.md
	docker-compose down --volumes --rmi all
	docker network rm local-shared

# @see https://github.com/CyberAgent/mimosa-common/tree/master/local
network:
	@if [ -z "`docker network ls | grep local-shared`" ]; then docker network create local-shared; fi

fmt: proto/**/*.proto
	@clang-format -i proto/**/*.proto

doc: fmt
	protoc \
		--proto_path=proto \
		--error_format=gcc \
		--doc_out=markdown,README.md:doc \
		proto/**/*.proto;

build: fmt
	protoc \
		--proto_path=proto \
		--error_format=gcc \
		--go_out=plugins=grpc,paths=source_relative:proto \
		proto/**/*.proto;

build-validate: fmt
	protoc \
		--proto_path=proto \
		--error_format=gcc \
		-I $GOPATH/src/github.com/envoyproxy/protoc-gen-validate \
		--go_out=plugins=grpc,paths=source_relative:proto \
		--validate_out="lang=go,paths=source_relative:proto" \
		proto/**/*.proto;

go-test: build
	cd proto/finding && go test ./...
	cd proto/iam     && go test ./...
	cd proto/project && go test ./...
	cd proto/alert   && go test ./...
	cd src/finding   && go test ./...
	cd src/iam       && go test ./...
	cd src/project   && go test ./...
	cd src/alert     && go test ./...

go-mod-tidy: build
	cd proto/finding && go mod tidy
	cd proto/iam     && go mod tidy
	cd proto/project && go mod tidy
	cd proto/alert   && go mod tidy
	cd src/finding   && go mod tidy
	cd src/iam       && go mod tidy
	cd src/project && go mod tidy
	cd src/alert     && go mod tidy

go-mod-update:
	cd src/finding \
		&& go get -u \
			github.com/CyberAgent/mimosa-core/...
	cd src/iam \
		&& go get -u \
			github.com/CyberAgent/mimosa-core/...
	cd src/project \
		&& go get -u \
			github.com/CyberAgent/mimosa-core/...
	cd src/alert \
		&& go get -u \
			github.com/CyberAgent/mimosa-core/...

run: go-test network
	. env.sh && docker-compose up -d --build

run-finding: go-test network
	. env.sh && docker-compose up -d --build finding

run-iam: go-test network
	. env.sh && docker-compose up -d --build iam

run-alert: go-test network
	. env.sh && docker-compose up -d --build alert

run-project: go-test network
	. env.sh && docker-compose up -d --build project

log:
	. env.sh && docker-compose logs -f

stop:
	. env.sh && docker-compose down
