APP_NAME = url-crawler-backend

run:
	go run cmd/main.go

build:
	go build -o bin/$(APP_NAME) cmd/main.go

tidy:
	go mod tidy

test:
	go test ./...

docker-build:
	docker build -t $(APP_NAME):latest .

docker-run:
	docker run -p 8080:8080 --env-file .env $(APP_NAME):latest