all:
	docker-compose-release

docker-compose-release:
	docker build --no-cache -t server .
	docker-compose up

logs:
	docker logs server