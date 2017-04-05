deps:
	@docker-compose --project-name mystack-logger up -d

stop-deps:
	@docker-compose --project-name mystack-logger down

run:
	@go run main.go start
