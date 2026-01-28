# Load Balancer Go

A high-performance HTTP load balancer built in Go with optional rate limiting and JWT authentication.

## Features

- **Load Balancing Algorithms**
  - Round Robin
  - Least Connections

- **Optional Security**
  - JWT-based authentication (optional)
  - Rate limiting with token bucket algorithm (optional)
  - Per-client IP rate limiting

- **Health Monitoring**
  - Automatic backend health checks
  - Real-time backend status monitoring

- **Performance**
  - HTTP reverse proxy
  - Concurrent request handling
  - Configurable middleware

## Installation

### From Source
```bash
git clone https://github.com/anjy7/Loadbalancer-go.git
cd Loadbalancer-go
go mod tidy
go build -o load-balancer main.go
```

## Quick Start

### 1. Create Configuration

Create a `config.json` file:

```json
{
    "port": "8080",
    "jwt_secret": "",
    "load_balancer": {
        "algorithm": "round_robin"
    },
    "rate_limit": {
        "enable": false,
        "requests_per_second": 100,
        "burst_size": 50
    },
    "backends": [
        {
            "id": "backend1",
            "url": "http://localhost:8081",
            "health": true
        },
        {
            "id": "backend2",
            "url": "http://localhost:8082",
            "health": true
        }
    ]
}
```

### 2. Run Load Balancer

```bash
go run main.go
```

### 3. Test It

```bash
# Basic request
curl http://localhost:8080/

# Health check
curl http://localhost:8080/health
```

## Configuration

### Basic Settings

| Field | Description | Example |
|-------|-------------|---------|
| `port` | Server port | `"8080"` |
| `jwt_secret` | JWT secret (leave empty to disable auth) | `""` |

### Load Balancer

```json
"load_balancer": {
    "algorithm": "round_robin"
}
```

**Available algorithms:**
- `round_robin` - Distributes requests evenly
- `least_connections` - Routes to backend with fewest connections

### Rate Limiting

```json
"rate_limit": {
    "enable": true,
    "requests_per_second": 100,
    "burst_size": 50
}
```

- Set `enable: false` to disable rate limiting
- `requests_per_second` - Sustained request rate
- `burst_size` - Maximum burst capacity

### Backends

```json
"backends": [
    {
        "id": "backend1",
        "url": "http://localhost:8081",
        "health": true
    }
]
```

## Authentication

When `jwt_secret` is set, all requests except `/health` require a JWT token:

```bash
curl -H "Authorization: Bearer YOUR_JWT_TOKEN" http://localhost:8080/
```

## Health Checks

The load balancer sends GET requests to `/health` on each backend every 30 seconds. Backends should return HTTP 200 for healthy status.

Available endpoints:
- `GET /health` - Load balancer health status

## Project Structure

```
load-balancer-go/
├── main.go                 # Application entry point
├── config/
│   └── config.go          # Configuration management
├── loadbalancer/
│   ├── roundrobin.go      # Round robin algorithm
│   ├── leastconn.go       # Least connections algorithm
│   └── balancer.go        # Load balancer interface
├── middleware/
│   ├── auth.go            # JWT authentication
│   └── ratelimiter.go     # Rate limiting
├── proxy/
│   └── http.go            # HTTP reverse proxy
└── health/
    └── checker.go         # Health monitoring
```

## Use Cases

- **API Gateway** - Route requests with optional authentication
- **Web Application Load Balancer** - Distribute traffic across multiple servers
- **Development Environment** - Simple load balancing for local testing

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test your changes
5. Submit a pull request

## License

This project is licensed under the MIT License.

---

<div align="center">

Made in Go

</div>
