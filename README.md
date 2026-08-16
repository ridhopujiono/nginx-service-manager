# Nginx Manager Service

A modern, production-ready Go service for managing Nginx configurations, reverse proxies, and SSL certificates with a RESTful API.

## ✨ Features

- 🔄 **Dynamic Proxy Management** - Create, update, and delete reverse proxies on the fly
- 🔐 **SSL/TLS Certificate Management** - Automated certificate handling and renewal
- 📊 **Health Monitoring** - Real-time service and Nginx status endpoints
- 🛡️ **Authentication & Authorization** - API key-based security for protected endpoints
- ⚙️ **Nginx Configuration** - Template-based Nginx config generation
- 🔄 **Hot Reload** - Apply changes without downtime
- 📝 **Comprehensive Logging** - Structured logging for debugging and monitoring
- 🐳 **Docker Ready** - Containerized deployment support

## 📋 Project Structure

```
nginx-manager-service/
├── cmd/
│   └── server/
│       └── main.go                 # Application entry point
├── internal/
│   ├── action/                     # Action handlers
│   ├── api/                        # API route handlers
│   ├── certificate/                # SSL certificate management
│   ├── config/                     # Configuration management
│   ├── middleware/                 # HTTP middleware (auth, logging)
│   ├── nginx/                      # Nginx configuration & control
│   ├── proxy/                      # Proxy management logic
│   └── status/                     # Service status handlers
├── templates/
│   └── reverse-proxy.conf.tmpl    # Nginx configuration template
├── go.mod                          # Go module definition
├── go.sum                          # Dependency checksums
└── README.md                       # This file
```

## 🚀 Getting Started

### Prerequisites

- **Go** 1.21 or higher
- **Nginx** 1.19 or higher
- **curl** or **Postman** (for testing API endpoints)

### Installation

1. Clone the repository:
```bash
git clone https://github.com/yourusername/nginx-manager-service.git
cd nginx-manager-service
```

2. Install dependencies:
```bash
go mod download
go mod tidy
```

3. Build the application:
```bash
go build -o nginx-manager ./cmd/server
```

### Configuration

Set the following environment variables:

```bash
# Nginx configuration directory
export NGINX_CONFIG_DIR="/etc/nginx/conf.d"

# Service port
export SERVICE_PORT="8080"

# API Key for authentication
export API_KEY="your-secret-api-key"
```

### Running the Service

```bash
# Development mode
go run cmd/server/main.go

# Production mode
./nginx-manager
```

The service will start on `http://localhost:8080`

## 🔌 API Endpoints

### Public Endpoints

- `GET /health` - Service health check
  ```bash
  curl http://localhost:8080/health
  ```

### Protected Endpoints (require API Key)

- `POST /api/v1/proxies` - Create a new reverse proxy
- `GET /api/v1/proxies` - List all proxies
- `GET /api/v1/proxies/:id` - Get proxy details
- `PUT /api/v1/proxies/:id` - Update a proxy
- `DELETE /api/v1/proxies/:id` - Delete a proxy

- `GET /api/v1/status` - Get service status
- `POST /api/v1/reload` - Reload Nginx configuration

- `POST /api/v1/certificates` - Upload SSL certificate
- `GET /api/v1/certificates` - List certificates

### Example Request

```bash
curl -X POST http://localhost:8080/api/v1/proxies \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "domain": "example.com",
    "backend": "http://localhost:3000",
    "port": 80
  }'
```

## 📚 Development

### Project Layout

- **`cmd/server/`** - Application entry point with HTTP server setup
- **`internal/nginx/`** - Nginx configuration management and validation
- **`internal/proxy/`** - Reverse proxy creation and management logic
- **`internal/certificate/`** - SSL/TLS certificate handling
- **`internal/status/`** - Service status monitoring
- **`internal/middleware/`** - Authentication and logging middleware
- **`internal/action/`** - Business logic for various operations
- **`templates/`** - Nginx configuration templates

### Adding New Features

1. Create handlers in appropriate `internal/` package
2. Register routes in `cmd/server/main.go`
3. Add middleware if authentication required
4. Write tests in `*_test.go` files
5. Update documentation

### Testing

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific test
go test -run TestHandler ./internal/status
```

## 🔒 Security

- API key authentication on protected endpoints
- Input validation on all API requests
- HTTPS/TLS support for secure communication
- Nginx config validation before reload
- Audit logging for all operations

## 📖 License

MIT License - See LICENSE file for details

## 🤝 Contributing

Contributions are welcome! Please follow these steps:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📞 Support

For issues, questions, or suggestions, please open an issue on GitHub or contact the maintainers.

---

**Last Updated:** 2026-08-16
