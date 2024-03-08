metrictest: 
	go vet -vettool=statictest ./...

	cd cmd/server && go build -buildvcs=false -o server main.go
	cd cmd/agent && go build -buildvcs=false -o agent main.go
	./metricstest -test.v -test.run=^TestIteration1$$ -binary-path=cmd/server/server
