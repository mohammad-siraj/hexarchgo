PROTODIR = ./internal/domain/adapters/proto/apis

generate-protos: $(PROTODIR)/*
		for file in $^ ; do \
			protoc -I ./internal/domain/adapters/proto \
			--go_out=./internal/domain/adapters/proto --go_opt=paths=source_relative \
			--go-grpc_out=./internal/domain/adapters/proto --go-grpc_opt=paths=source_relative \
			--grpc-gateway_out=./internal/domain/adapters/proto --grpc-gateway_opt=paths=source_relative \
			./$$file/*.proto; \
		done
