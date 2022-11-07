IMAGE_TAG = 0.1
SQL_SCRIPT_PATH = internal/database/scripts/configsdb.sql
SQL_USER = postgres

build:
	go build -o bin/config-controller ./cmd/config-controller/

clean:
	rm -rf ./bin

init-database:
	psql -U $(SQL_USER) -f $(SQL_SCRIPT_PATH)

docker-build:
	sudo docker build --tag config-controller:$(IMAGE_TAG) -f ./docker/config-controller.Dockerfile .


docker-compose-up:
	docker-compose -f ./docker/docker-compose.yaml up

docker-compose-down-all:
	docker-compose -f ./docker/docker-compose.yaml down -v --rmi all

docker-compose-start:
	docker-compose -f ./docker/docker-compose.yaml start

docker-compose-stop:
	docker-compose -f ./docker/docker-compose.yaml stop