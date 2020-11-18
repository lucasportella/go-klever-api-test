generate:
	@protoc --proto_path=proto/ --proto_path=third_party --go_out=plugins=grpc:third_party/gen post.proto

run:
	@echo "---- Running Server ----"
	@go run server/main.go

run_client:
	@echo "---- Running Client ----"
	@go run client/main.go