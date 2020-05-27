.PHONY: all clean fmt build doc run log stop
all: run

clean:
	rm -f proto/*/*.pb.go
	# docker-compose down --volumes --rmi all

fmt: proto/*/*.proto
	clang-format -i proto/*/*.proto

build: fmt
	protoc \
		--proto_path=proto \
		--error_format=gcc \
		--go_out=plugins=grpc,paths=source_relative:proto \
		proto/*/*.proto
	# docker-compose build

doc: fmt
	protoc \
		--proto_path=proto \
		--error_format=gcc \
		--doc_out=markdown,README.md:doc \
		proto/*/*.proto

# run: build
# 	docker-compose up -d

# log:
# 	docker-compose logs -f

# stop:
# 	docker-compose down
