# Nginx Manager Service

A Go-based service for managing Nginx configuration and reverse proxy setups.

## Project Structure

```
nginx-manager-service/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── api/
│   ├── nginx/
│   └── config/
├── templates/
│   └── reverse-proxy.conf.tmpl
├── go.mod
└── README.md
```

## Getting Started

### Prerequisites
- Go 1.21 or higher
- Nginx

### Build and Run

```bash
go run cmd/server/main.go
```

## Development

Add your API handlers in `internal/api/`, Nginx management logic in `internal/nginx/`, and configuration handling in `internal/config/`.
