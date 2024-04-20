FORCE:

DOCKER_COMPOSE_DB=docker-compose.yml


up-db:
	docker-compose -f ${DOCKER_COMPOSE_DB} up  -d

down-db:
	docker-compose -f ${DOCKER_COMPOSE_DB} down 

logs-db:
	docker-compose -f ${DOCKER_COMPOSE_DB} logs 

statictest: FORCE
	golangci-lint run
	go vet -vettool=./statictest ./...

buld-server:
	@echo "Build server"
	@cd cmd/server && go build -buildvcs=false -o server

build-agent:
	@echo "Build agent"
	@cd cmd/agent && go build -buildvcs=false -o agent

metrictest: buld-server build-agent
	@echo "Start metrictest"
	./metricstest -test.v -test.run=^TestIteration$(iter)$$ -binary-path=cmd/server/server -agent-binary-path=cmd/agent/agent -server-port=8080 -source-path=. -file-storage-path=./metrics-db.json -database-dsn="host=localhost user=db_user password=db_password dbname=db_metric sslmode=disable"
