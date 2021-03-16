build-proto:
	@echo "Compiling blog.proto"
	protoc blogpb/blog.proto --go_out=plugins=grpc:.
.PHONY: start

start-server: build-proto
	@echo "Starting Blog"
	APP_ENV=development \
	go run server/server.go
.PHONY: start