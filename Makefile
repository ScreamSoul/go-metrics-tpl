# Makefile

statictest:
	golangci-lint run
	go vet -vettool=./statictest ./...

metrictest:
	@echo "Build server"
	@cd cmd/server && go build -buildvcs=false -o server
	@echo "Build agent"
	@cd cmd/agent && go build -buildvcs=false -o agent
	@echo "Start metrictest"
	./metricstest -test.v -test.run=^TestIteration$(iter)$$ -binary-path=cmd/server/server -agent-binary-path=cmd/agent/agent -server-port=8080 -source-path=. -file-storage-path=./metrics-db.json -database-dsn="host=localhost user=db_user password=db_password dbname=db_metric sslmode=disable"
