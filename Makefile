.PHONY: all install clean network fmt build doc run log stop
all: run

install:
	go get \
		google.golang.org/grpc \
		github.com/golang/protobuf/protoc-gen-go \
		github.com/grpc-ecosystem/go-grpc-middleware \
		github.com/mwitkow/go-proto-validators/protoc-gen-govalidators

clean:
	rm -f proto/*/*.pb.go
	rm -f doc/*.md
	docker-compose down --volumes --rmi all
	docker network rm local-shared

# @see https://github.com/CyberAgent/mimosa-common/tree/master/local
network:
	@if [ -z "`docker network ls | grep local-shared`" ]; then docker network create local-shared; fi

fmt: proto/*/*.proto
	clang-format -i proto/**/*.proto

build: BUILD_TARGETS := finding iam
build: fmt
	for target in $(BUILD_TARGETS); \
	do \
		protoc \
			--proto_path=proto \
			--proto_path=${GOPATH}/src \
			--proto_path=${GOPATH}/src/github.com/gogo/protobuf/protobuf \
			--error_format=gcc \
			--go_out=plugins=grpc,paths=source_relative:proto \
			--govalidators_out=paths=source_relative:proto \
			proto/$$target/*.proto; \
	done

test: build
	cd src/gateway && go test ./...
	cd src/finding && go test ./...

doc: fmt
	protoc \
		--proto_path=proto \
		--proto_path=${GOPATH}/src \
		--error_format=gcc \
		--doc_out=markdown,README.md:doc \
		proto/**/*.proto

run: test network
	. env.sh && docker-compose up -d --build

log:
	. env.sh && docker-compose logs -f

stop:
	. env.sh && docker-compose down

ssh:
	. env.sh && docker-compose exec gateway sh
