build:
	go build -o ./bin/app ./main.go
run:
	./bin/app cli
lint:
	gofmt -s -w ./ && golangci-lint run