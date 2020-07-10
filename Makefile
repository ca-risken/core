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
		--proto_path=${GOPATH}/src \
		--error_format=gcc \
		--doc_out=markdown,README.md:doc \
		proto/**/*.proto;

build: fmt
	protoc \
		--proto_path=proto \
		--proto_path=${GOPATH}/src \
		--error_format=gcc \
		--go_out=plugins=grpc,paths=source_relative:proto \
		proto/**/*.proto;

go-test: build
	cd proto/finding && go test ./...
	cd proto/iam     && go test ./...
	cd src/finding   && go test ./...
	cd src/iam       && go test ./...

go-mod-update:
	cd src/finding \
		&& go get -u \
			github.com/CyberAgent/mimosa-core/proto/finding \
			github.com/CyberAgent/mimosa-core/pkg/model
	cd src/iam \
		&& go get -u \
			github.com/CyberAgent/mimosa-core/proto/iam \
			github.com/CyberAgent/mimosa-core/pkg/model

go-mod-tidy: build
	cd proto/finding && go mod tidy
	cd proto/iam     && go mod tidy
	cd src/finding   && go mod tidy
	cd src/iam       && go mod tidy

run: go-test network
	. env.sh && docker-compose up -d --build

log:
	. env.sh && docker-compose logs -f

stop:
	. env.sh && docker-compose down
