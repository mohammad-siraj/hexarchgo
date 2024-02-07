
PROTODIR = ./internal/domain

generate-protos: $(PROTODIR)/*
		for file in $^ ; do \
			protoc -I ./$$file/driving/adapters/proto \
			--go_out=./$$file/driving/adapters/proto --go_opt=paths=source_relative \
			--go-grpc_out=./$$file/driving/adapters/proto --go-grpc_opt=paths=source_relative \
			--grpc-gateway_out=./$$file/driving/adapters/proto --grpc-gateway_opt=paths=source_relative \
			./$$file/driving/adapters/proto/service/*.proto; \
		done
