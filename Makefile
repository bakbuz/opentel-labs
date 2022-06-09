.PHONY: pb up down clean v swag

pb:
	@echo $(GOPATH)
	protoc --proto_path=./protos/ ./protos/*.proto --go_out=./common-service/pb --go_opt=paths=source_relative
	protoc --proto_path=./protos/ ./protos/*.proto --go-grpc_out=./common-service/pb --go-grpc_opt=paths=source_relative --plugin=$(GOPATH)/bin/protoc-gen-go-grpc

	protoc --proto_path=./protos/ ./protos/*.proto --go_out=./restapi/pb --go_opt=paths=source_relative
	protoc --proto_path=./protos/ ./protos/*.proto --go-grpc_out=./restapi/pb --go-grpc_opt=paths=source_relative --plugin=$(GOPATH)/bin/protoc-gen-go-grpc

up:
	@echo "up"
	sudo docker-compose up --build

down:
	@echo "down"
	sudo docker-compose down
	sudo docker image rm maydere/common-service --force
	sudo docker image rm maydere/restapi --force

clean:
	go clean -modcache

v:
	protoc --version
	protoc-gen-go --version

swag:
	cd restapi && $(GOPATH)/bin/swag init -g main.go

upgr:
	cd common-service && go get -u all
	cd common-service && go mod tidy
	cd restapi && go get -u all
	cd restapi && go mod tidy