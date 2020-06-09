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
	clang-format -i proto/**/*.proto

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
	cd src/gateway && go test ./...
	cd src/finding && go test ./...
	cd proto/finding && go test ./...

go-mod-update:
	cd src/gateway && go get -u github.com/CyberAgent/mimosa-core/proto/iam
	cd src/finding && go get -u github.com/CyberAgent/mimosa-core/proto/finding && go get -u github.com/CyberAgent/mimosa-core/pkg/model

go-mod-tidy: build
	cd src/gateway && go mod tidy
	cd src/finding && go mod tidy

run: go-test network
	. env.sh && docker-compose up -d --build

log:
	. env.sh && docker-compose logs -f

stop:
	. env.sh && docker-compose down

ssh:
	. env.sh && docker-compose exec gateway sh

####################################################
## grpcurl example
####################################################
grpcurl-list:
	grpcurl -plaintext localhost:8081 list

grpcurl-fiding-list:
	grpcurl \
		-plaintext \
		-d '{"project_id":["1001","1002"], "from_score":0.0, "to_score":1.0, "from_at": 1560000000, "to_at": 1660000000}' \
		localhost:8081 core.finding.FindingService.ListFinding

grpcurl-finding-put:
	grpcurl \
		-plaintext \
		-d '{"finding":{"description":"desc", "data_source":"ds", "data_source_id":"ds-001", "resource_name":"rn", "project_id":1001, "original_score":55.51, "original_max_score":100.0, "data":"{\"key\":\"value\"}"}}' \
		localhost:8081 core.finding.FindingService.PutFinding

grpcurl-finding-delete:
	grpcurl \
		-plaintext \
		-d '{"finding_id": 1004}' \
		localhost:8081 core.finding.FindingService.DeleteFinding
