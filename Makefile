
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

build-appimage:
	@echo "Building image";\
	CGO_ENABLED=0 GOOS=linux go build -o app cmd/main.go;\
	chmod 777 app;\
	sudo docker build . -t modelapp;\

deploy-docker:
	make build-appimage;\
	@echo "Deploying docker image...";\
	cd infrastructure/docker && docker-compose up -d;\
	docker image prune;\
	rm app;

down-docker:
	@echo "bringing down docker images...";\
	cd infrastructure/docker && docker-compose down;\
	docker volume prune;\
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
	make deploy-docker;

#sleep 5s;\
	#make migrate-up;\
	#go run cmd/main.go;\
	#make down-docker;\

stop-server:
	make down-docker;\
	docker volume prune -a;\


terraform:
	@tflocal -chdir=infrastructure/cloud/terraform init
	@tflocal -chdir=infrastructure/cloud/terraform apply --auto-approve

localstack:
	@ cd infrastructure/cloud  && docker compose up -d && cd ../..

localstack-down:
	@ cd infrastructure/cloud  && docker compose down && cd ../..


zip:
	@CGO_ENABLED=0 GOOS=linux  go build -ldflags="-d -s -w" -o app ./cmd/serverless/main.go
	@chmod 777 app 
	@zip infrastructure/cloud/terraform/app.zip app
	@chmod 777 infrastructure/cloud/terraform/app.zip
	@rm app

start-lambda:
	@make localstack
	#@make zip
	@echo Its here
	@make terraform
	#@make localstack-down

test-invoke-lambda:
	@bash ./infrastructure/cloud/invokelambda.sh


#for minikube users use minikube service --url to get the url to acces the app
deploy-k8:
	@make build-appimage;\
	minikube start;\
	docker save modelapp:latest -o modelapp.tar;\
	chmod 777 modelapp.tar;\
	docker cp modelapp.tar minikube:/var;\
	docker exec -it minikube docker load -i /var/modelapp.tar;\
	cd infrastructure;\
	kubectl apply -f k8;\


down-k8:
	@cd infrastructure;\
	kubectl delete -f k8;\
	minikube stop;\


build-zip:
	#dep ensure -v
	env GOOS=linux go build -ldflags="-d -s -w" -o app ./cmd/serverless/main.go
	@chmod 777 app 
	build-lambda-zip -o infrastructure/cloud/terraform/app.zip app
	@chmod 777 infrastructure/cloud/terraform/app.zip
	@rm app