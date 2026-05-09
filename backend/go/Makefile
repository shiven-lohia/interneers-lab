-include .env.local


.tidy:
	go mod tidy

.setup-env:
	cp .env.sample .env.local

setup: .tidy .setup-env


.PHONY: .build
.build:
	go build -o bin/app cmd/app/main.go

run:
	./bin/app

build-and-run: .build run



mongo-login:
	mongosh mongodb://127.0.0.1:${MONGO_PORT}  
# Catch-all target to handle positional arguments (prevents "No rule to make target" error)
%:
	@: