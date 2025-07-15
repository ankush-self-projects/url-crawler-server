# URL Crawler Backend

A Go-based web crawler backend service that analyzes URLs and extracts metadata including HTML version, page titles, headings, links, and login forms.

## Features

- URL crawling and analysis
- HTML version detection
- Page title extraction
- Heading structure analysis
- Internal/external link counting
- Broken link detection
- Login form detection
- JWT authentication
- MySQL database storage

## Prerequisites

- Go 1.24.5 or higher
- MySQL 8.0 or higher
- Git

## Setup

### 1. Clone the repository
```bash
git clone https://github.com/ankush-self-projects/url-crawler-server
cd url-crawler-backend
```

### 2. Install dependencies
```bash
go mod tidy
```

### 3. Database setup
Create a MySQL database:
```sql
CREATE DATABASE url_info;
```

### 4. Environment configuration
Create a `.env` file in the project root:
```env
DB_USER=root
DB_PASS=your_mysql_password
DB_HOST=localhost
DB_PORT=3306
DB_NAME=url_crawler
JWT_SECRET=your_jwt_secret_key_here
```

### 5. Run database migrations and seed data
```bash
go run cmd/main.go
```

### 6. Seed initial user (optional)
```bash
go run cmd/seed/main.go
```
This creates an admin user with:
- Username: `admin`
- Password: `testpassword`

## Running the application

### Development
```bash
make run
```
or
```bash
go run cmd/main.go
```

### Production build
```bash
make build
./bin/url-crawler-backend
```

## Testing

The project includes comprehensive test coverage with unit tests, integration tests, and API tests.

### Run all tests
```bash
make test
```

### Run specific test types
```bash
# Unit tests only
make test-unit

# Integration tests only
make test-integration

# Generate coverage report
make test-coverage
```

### Manual test commands
```bash
# Run all tests with verbose output
go test ./... -v

# Run tests for specific package
go test ./internal/api -v
go test ./internal/crawler -v
go test ./tests -v

# Run tests with coverage
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

### Test Coverage
The test suite covers:
- **Authentication**: Login functionality, JWT token generation
- **API Handlers**: URL management, crawling operations
- **Crawler Logic**: HTML parsing, link analysis, form detection
- **Integration**: Complete workflows from API to database

## API Endpoints

### Authentication
All API endpoints require JWT authentication. Include the JWT token in the Authorization header:
```
Authorization: Bearer <your_jwt_token>
```

### URL Management

#### Add a new URL for crawling
```http
POST /api/urls
Content-Type: application/json

{
  "url": "https://example.com"
}
```

#### Get all URLs
```http
GET /api/urls
```

#### Start crawling a specific URL
```http
POST /api/urls/{id}/start
```

## Response Format

### URL Object
```json
{
  "id": 1,
  "url": "https://example.com",
  "html_version": "HTML5",
  "page_title": "Example Domain",
  "headings": "H1: 1, H2: 2, H3: 3",
  "internal_links": 5,
  "external_links": 3,
  "broken_links": 0,
  "has_login_form": false,
  "status": "done",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

## Docker

### Build image
```bash
make docker-build
```

### Run container
```bash
make docker-run
```

## Project Structure

```
url-crawler-backend/
├── cmd/
│   ├── main.go          # Application entry point
│   └── seed/            # Database seeding
├── internal/
│   ├── api/             # HTTP handlers and routes
│   │   ├── auth.go      # Authentication handlers
│   │   ├── handlers.go  # URL management handlers
│   │   ├── routes.go    # Route definitions
│   │   └── *_test.go    # API tests
│   ├── crawler/         # Web crawling logic
│   │   └── *_test.go    # Crawler tests
│   ├── db/              # Database connection
│   ├── middleware/      # JWT middleware
│   └── model/           # Data models
├── tests/               # Integration tests
├── .env                 # Environment variables
├── go.mod              # Go modules
├── Makefile            # Build commands
└── README.md           # This file
```

## Development

### Available Make commands
- `make run` - Run the application
- `make build` - Build the binary
- `make tidy` - Install dependencies
- `make test` - Run all tests
- `make test-unit` - Run unit tests only
- `make test-integration` - Run integration tests only
- `make test-coverage` - Generate coverage report
- `make docker-build` - Build Docker image
- `make docker-run` - Run Docker container

## Troubleshooting

### Database connection issues
- Ensure MySQL is running
- Verify database credentials in `.env`
- Check if the database exists

### JWT authentication issues
- Ensure `JWT_SECRET` is set in `.env`
- Verify the JWT token is valid and not expired

### Crawling issues
- Check if the target URL is accessible
- Verify network connectivity
- Review server logs for specific error messages

### Test issues
- Ensure all dependencies are installed: `go mod tidy`
- Check that SQLite is available for tests
- Verify test database connections are working