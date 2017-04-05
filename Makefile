build:
	@mkdir -p bin && go build -o ./bin/mystack-logger main.go

build-docker: cross-build-linux-amd64
	@docker build -t mystack-logger .

deps:
	@docker-compose --project-name mystack-logger up -d

cross-build-linux-amd64:
	@env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ./bin/mystack-logger-linux-amd64
	@chmod a+x ./bin/mystack-logger-linux-amd64

run:
	@go run main.go start

stop-deps:
	@docker-compose --project-name mystack-logger down
