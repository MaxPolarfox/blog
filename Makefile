build-proto:
	@echo "Compiling blog.proto"
	protoc blogpb/blog.proto --go_out=plugins=grpc:.
.PHONY: start

start: build-proto
	@echo "Starting Blog"
	APP_ENV=development \
	go run cmd/server/main.go
.PHONY: start

requests::
	@echo "Sending requests to Blog"
	APP_ENV=development \
	go run cmd/client/main.go
.PHONY: start