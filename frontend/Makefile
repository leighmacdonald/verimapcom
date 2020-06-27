all: gen
gen:
	@protoc -I../pb rpc.proto    --js_out=import_style=commonjs:src/pb \
    		--grpc-web_out=import_style=typescript,mode=grpcwebtext:src/pb
