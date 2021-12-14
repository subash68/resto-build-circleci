protoc --proto_path=api/proto/v1 --go_out=plugins=grpc:pkg/api setting-service.proto
protoc --proto_path=api/proto/v1 --grpc-gateway_out=logtostderr=true:pkg/api setting-service.proto
protoc --proto_path=api/proto/v1 --swagger_out=logtostderr=true:api/swagger/v1 setting-service.proto