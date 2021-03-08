.PHONY: all install clean network fmt build doc go-test go-mod-update go-mod-tidy run log stop
all: run

install:
	go get \
		google.golang.org/grpc \
		github.com/golang/protobuf/protoc-gen-go \
		github.com/envoyproxy/protoc-gen-validate \
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
		-I $(GOPATH)/src/github.com/envoyproxy/protoc-gen-validate \
		--doc_out=markdown,README.md:doc \
		proto/**/*.proto;

# build without protoc-gen-validate
build: fmt doc
	for svc in "alert" "finding" "iam"; do \
		protoc \
			--proto_path=proto \
			--error_format=gcc \
			--go_out=plugins=grpc,paths=source_relative:proto \
			proto/$$svc/*.proto; \
	done

# build with protoc-gen-validate
build-validate: fmt doc
	for svc in "project" "report"; do \
		protoc \
			--proto_path=proto \
			--error_format=gcc \
			-I $(GOPATH)/src/github.com/envoyproxy/protoc-gen-validate \
			--go_out=plugins=grpc,paths=source_relative:proto \
			--validate_out="lang=go,paths=source_relative:proto" \
			proto/$$svc/*.proto; \
	done

go-test: build build-validate
	cd proto/finding && go test ./...
	cd proto/iam     && go test ./...
	cd proto/project && go test ./...
	cd proto/alert   && go test ./...
	cd proto/report  && go test ./...
	cd src/finding   && go test ./...
	cd src/iam       && go test ./...
	cd src/project   && go test ./...
	cd src/alert     && go test ./...
	cd src/report    && go test ./...

go-mod-tidy: build
	cd proto/finding && go mod tidy
	cd proto/iam     && go mod tidy
	cd proto/project && go mod tidy
	cd proto/alert   && go mod tidy
	cd proto/report  && go mod tidy
	cd src/finding   && go mod tidy
	cd src/iam       && go mod tidy
	cd src/project   && go mod tidy
	cd src/alert     && go mod tidy
	cd src/report    && go mod tidy

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
	cd src/report \
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

log-finding:
	. env.sh && docker-compose logs -f finding

log-project:
	. env.sh && docker-compose logs -f project

log-report:
	. env.sh && docker-compose logs -f report

stop:
	. env.sh && docker-compose down
