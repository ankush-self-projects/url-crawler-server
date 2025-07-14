# URL Crawler Backend

## Setup
- Copy `.env.example` to `.env` and update DB creds
- Run `make tidy` to install modules
- Run `make run` to start server

## Endpoints
- `POST /api/urls`
- `GET /api/urls`
- `POST /api/urls/:id/start`

## Docker
- `make docker-build`
- `make docker-run`