all:
	docker-compose-release

docker-compose-release:
	docker build --no-cache -t server .
	docker-compose up

docker-clear:
	docker rmi -f server
	docker rmi -f postgres:14.3-alpine

logs:
	docker logs server