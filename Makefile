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
	./metricstest -test.v -test.run=^TestIteration$(iter)$$ -binary-path=cmd/server/server -agent-binary-path=cmd/agent/agent -server-port=8080 -source-path=.
