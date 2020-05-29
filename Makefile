.PHONY: all clean fmt build doc run log stop
all: run

clean:
	rm -f proto/*/*.pb.go
	rm -f doc/*.md
	. env.sh && docker-compose down --volumes --rmi all

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

run: test
	. env.sh && docker-compose up -d --build

log:
	. env.sh && docker-compose logs -f

stop:
	. env.sh && docker-compose down
