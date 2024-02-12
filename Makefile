
PROTODIR = ./internal

generate-protos: $(PROTODIR)/*
		for file in $^ ; do \
			protoc -I ./$$file/driving/adapters/proto \
			--openapiv2_out ./$$file/driving/adapters/proto --openapiv2_opt use_go_templates=true\
			--go_out=./$$file/driving/adapters/proto --go_opt=paths=source_relative \
			--go-grpc_out=./$$file/driving/adapters/proto --go-grpc_opt=paths=source_relative \
			--grpc-gateway_out=./$$file/driving/adapters/proto --grpc-gateway_opt=paths=source_relative \
			./$$file/driving/adapters/proto/service/*.proto; \
		done

deploy-docker:
		@echo "Deploying docker image...";\
		cd infrastructure/docker && docker-compose up -d;

down-docker:
		@echo "bringing down docker images...";\
		cd infrastructure/docker && docker-compose down;