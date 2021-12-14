# --go-grpc_out=pkg/api
protoc --proto_path=api/proto/v1 --proto_path=third_party --go_out=pkg/api --go-grpc_out=pkg/api dispatcher-service.proto
protoc --proto_path=api/proto/v1 --proto_path=third_party --grpc-gateway_out=logtostderr=true:pkg/api dispatcher-service.proto
protoc --proto_path=api/proto/v1 --proto_path=third_party --swagger_out=logtostderr=true:api/swagger/v1 dispatcher-service.proto