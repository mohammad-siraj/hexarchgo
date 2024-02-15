
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
	cd infrastructure/docker && docker-compose down;\
	echo "docker image down Successfully";

sql-migarte:
	@echo "Migrating sql database...";\
	migrate create -ext sql -dir data/database/migration/ -seq sqlmigrations;\
	@echo "Migrate sql database completed!";\

migrate-up:
	@echo "Running migrate up...";\
	migrate -path data/database/migration/ -database "postgresql://postgres:postgres@localhost:5432/mainserver?sslmode=disable" -verbose up;\
	@echo "Successfully migrate database schema";\

migrate-down:
	@echo "Running migrate up...";\
	migrate -path data/database/migration/ -database "postgresql://postgres:postgres@localhost:5432/mainserver?sslmode=disable" -verbose down;\
	echo "Successfully migrate database schema";\

bundle-openapi:
	touch api/openapi.yaml;\
	for file in $^ ; do \
		redocly join openapi.yaml ./$$file/driving/adapters/proto/service/*.yaml -o api/openapi.yaml
	done

start-server:
	go mod tidy;\
	make deploy-docker;\
	sleep 5s;\
	make migrate-up;\
	go run cmd/main.go;\

stop-server:
	make down-docker;