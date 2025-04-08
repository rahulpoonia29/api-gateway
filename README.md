# API Gateway

A simple API gateway written in Go.

## Getting Started

### Prerequisites

-   Go 1.20 or higher

### Running the application

```bash
cd cmd/api-gateway
go run main.go
```

## Project Structure

```text
├── cmd/
│   └── api-gateway/
│       └── main.go       # Your application entry point (move from root)
├── internal/             # Private application code
│   ├── handlers/         # HTTP request handlers
│   ├── middleware/       # HTTP middleware
│   └── service/          # Business logic
├── pkg/                  # Public libraries that can be used by other applications
├── api/                  # API documentation, OpenAPI/Swagger specs
├── configs/              # Configuration files
├── scripts/              # Build and deployment scripts
├── test/                 # Test files
├── go.mod                # Module definition
├── go.sum                # Dependencies checksums (will be generated)
└── README.md             # Project documentation
```
