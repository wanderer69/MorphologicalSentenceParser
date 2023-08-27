cd ../api/grpc/
# protoc -I. --go_out . --go_opt paths=source_relative --go-grpc_out=. ./morphological_sentense_parser.proto
protoc -I. --go_out . --go-grpc_out=. ./morphological_sentense_parser.proto
# protoc -I. -I%GOPATH%\src -IC:/Development/Go_pkg/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.16.0/third_party/googleapis -IC:/Development/Go_pkg/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.16.0 --go_out=plugins=grpc:. auth.proto
# protoc -I. -I%GOPATH%\src -IC:/Development/Go_pkg/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.16.0/third_party/googleapis -IC:/Development/Go_pkg/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.16.0 --grpc-gateway_out=logtostderr=true:. bookstore.proto
# protoc -I. -I%GOPATH%\src -IC:/Development/Go_pkg/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.16.0/third_party/googleapis -IC:/Development/Go_pkg/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.16.0 -I%GOPATH%\src\github.com\grpc-ecosystem\grpc-gateway --swagger_out=logtostderr=true:. http_bookstore.proto
