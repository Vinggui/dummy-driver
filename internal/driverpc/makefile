# For now, protoc-gen-go is obtained from "https://github.com/golang/protobuf"
# This is an old repo that is back compatible with the new one.
# The new repo "google.golang.org/protobuf" has the following command:
# "go install google.golang.org/protobuf/cmd/protoc-gen-go"
# But as for version 1.23.0, this command does not install the required
# "protoc-gen-go-grpc". So for future reference, this may be already working.
.PHONY: all
all:
	protoc --go_out=plugins=grpc,paths=source_relative:. driverpc.proto 