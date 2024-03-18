SET api_name=template
SET swagger=".\www\swagger.json"

if EXIST %swagger% (
    DEL %swagger%
)

protoc --proto_path=./proto ./proto/*.proto ^
    --proto_path=./proto/libs ^
    --go_out=./server/pb --go_opt paths=source_relative ^
    --govalidators_out=./server

protoc --proto_path=./proto ./proto/%api_name%_api.proto ^
    --proto_path=./proto/libs ^
    --go-grpc_out=./server/pb --go-grpc_opt paths=source_relative ^
    --grpc-gateway_out ./server/pb ^
    --grpc-gateway_opt allow_delete_body=true,logtostderr=true,paths=source_relative,repeated_path_param_separator=ssv ^
    --openapiv2_out ./proto ^
    --openapiv2_opt logtostderr=true,repeated_path_param_separator=ssv

COPY "%~dp0proto\%api_name%_api.swagger.json" "%~dp0www\swagger.json"
DEL "%~dp0proto\%api_name%_api.swagger.json"

protoc-go-inject-tag -input="./server/pb/*.pb.go"