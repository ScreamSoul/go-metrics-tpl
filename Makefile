metrictest: 
	cd cmd/server && go build -o server *.go
	cd cmd/agent && go build -o agent *.go
	./metricstest -test.v -test.run=^TestIteration5$ -agent-binary-path=./cmd/agent -binary-path=./cmd/server -server-port=8080 -source-path=.
