PROTO_SRC="paws.proto"

# Golang 
protoc $PROTO_SRC \
--go-grpc_out=../backend/pkg \
--go_out=../backend/pkg

# TS Electron
protoc $PROTO_SRC \
--plugin=protoc-gen-ts_proto=..\\frontend\\node_modules\\.bin\\protoc-gen-ts_proto.cmd \
--ts_proto_out=../frontend/src/main/proto \
--ts_proto_opt=outputServices=grpc-js \
--ts_proto_opt=esModuleInterop=true

# TS React
protoc $PROTO_SRC \
--plugin=protoc-gen-ts_proto=..\\frontend\\node_modules\\.bin\\protoc-gen-ts_proto.cmd \
--ts_proto_out=../frontend/src/common/proto \
--ts_proto_opt=outputEncodeMethods=false,outputJsonMethods=false,outputClientImpl=false

if [ $? -ne 0 ]; then
  read -p "Press any key to continue" -u 1
fi