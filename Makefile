serv = server
db = database_api

HELP_FUNC = \
	%help; while(<>){push@{$$help{$$2//'options'}},[$$1,$$3] \
	if/^([\w-_]+)\s*:.*\#\#(?:@(\w+))?\s(.*)$$/}; \
    print"$$_:\n", map"  $$_->[0]".(" "x(20-length($$_->[0])))."$$_->[1]\n",\
    @{$$help{$$_}},"\n" for keys %help; \

all: ##@App application in docker container
	docker-compose-api

docker-compose-api: ##@Runs application in docker container
	docker build --no-cache -t $(serv) .
	docker-compose up

clean-pgdata: ##@DB clean a database saved data
	rm -rf db/pgdata

docker-stop-api: ##@Server stops containers
	docker stop $(db)
	docker stop $(serv)

docker-clean-api: docker-stop-api ##@Server delete server and database containers
	docker rm $(db)
	docker rm $(serv)

server-logs: ##@Server show logs from server container
	docker logs $(serv)

database-logs:  ##@DB show logs from database container
	docker logs $(db)

all-logs: database-logs server-logs ##@App show logs from server and db containers together

help: ##@App Show this help
	@echo -e "Usage: make [target] ...\n"
	@perl -e '$(HELP_FUNC)' $(MAKEFILE_LIST)