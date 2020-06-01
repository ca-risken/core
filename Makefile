.PHONY: all clean network fmt build doc run log stop
all: run

clean:
	rm -f proto/*/*.pb.go
	rm -f doc/*.md
	docker-compose down --volumes --rmi all
	docker network rm local-shared

# @see https://github.com/CyberAgent/mimosa-common/tree/master/local
network:
	@if [ -z "`docker network ls | grep local-shared`" ]; then docker network create local-shared; fi

fmt: proto/*/*.proto
	clang-format -i proto/*/*.proto

build: fmt
	protoc \
		--proto_path=proto \
		--error_format=gcc \
		--go_out=plugins=grpc,paths=source_relative:proto \
		proto/*/*.proto

test: build
	cd src/gateway && go test ./...
	cd src/finding && go test ./...

doc: fmt
	protoc \
		--proto_path=proto \
		--error_format=gcc \
		--doc_out=markdown,README.md:doc \
		proto/*/*.proto

run: test network
	. env.sh && docker-compose up -d --build

log:
	. env.sh && docker-compose logs -f

stop:
	. env.sh && docker-compose down

ssh:
	. env.sh && docker-compose exec gateway sh
