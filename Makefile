APP_NAME = url-crawler-backend

run:
	go run cmd/main.go

build:
	go build -o bin/$(APP_NAME) cmd/main.go

tidy:
	go mod tidy

test:
	go test ./...

test-unit:
	go test ./internal/api -v
	go test ./internal/crawler -v

test-integration:
	go test ./tests -v

test-coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

docker-build:
	docker build -t $(APP_NAME):latest .

docker-run:
	docker run -p 8080:8080 --env-file .env $(APP_NAME):latest