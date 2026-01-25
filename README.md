# POS API - Point of Sales Backend

A modern, scalable REST API backend for Point of Sales (POS) systems built with Go and best practices.

![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)
![License](https://img.shields.io/badge/License-MIT-green.svg)

## ✨ Features

- **RESTful API** with versioning (`/api/v1/`)
- **JWT Authentication** with role-based access control (Admin, Manager, Cashier)
- **SQLite Database** (D1 Cloudflare compatible)
- **Rate Limiting** per IP address
- **Standardized Responses** with validation errors
- **Docker Ready** with multi-stage builds
- **Health Checks** for monitoring
- **Sales Reports** (daily, monthly, top products)

## 🏗️ Architecture

```
├── cmd/api/              # Application entry point
├── internal/
│   ├── config/           # Configuration management
│   ├── database/         # Database connection & migrations
│   ├── dto/              # Data Transfer Objects
│   ├── handler/          # HTTP handlers
│   ├── middleware/       # Auth, CORS, Rate Limit, Logger
│   ├── models/           # Domain models
│   ├── repository/       # Data access layer
│   ├── router/           # Route definitions
│   ├── service/          # Business logic
│   └── utils/            # Helpers (JWT, Response, Validation)
├── scripts/              # Database seeder
└── docs/                 # API documentation
```

## 🚀 Quick Start

### Prerequisites

- Go 1.21 or higher
- Docker (optional)

### Local Development

1. **Clone & Install Dependencies**

   ```bash
   git clone <repository-url>
   cd goland-dasar
   go mod tidy
   ```

2. **Configure Environment**

   ```bash
   cp .env.example .env
   # Edit .env as needed
   ```

3. **Seed Database**

   ```bash
   make seed
   # Or: go run ./scripts/seed.go
   ```

4. **Run Server**

   ```bash
   make dev
   # Or: go run ./cmd/api/main.go
   ```

5. **Test Health Check**
   ```bash
   curl http://localhost:8080/health
   ```

### Docker

```bash
# Build and run
docker-compose up --build

# Stop
docker-compose down
```

## 🔐 Authentication

### Test Credentials

| Email             | Password    | Role    |
| ----------------- | ----------- | ------- |
| admin@pos.local   | Admin123!   | admin   |
| manager@pos.local | Manager123! | manager |
| cashier@pos.local | Cashier123! | cashier |

### Login Example

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "admin@pos.local", "password": "Admin123!"}'
```

### Using Token

```bash
curl http://localhost:8080/api/v1/products \
  -H "Authorization: Bearer <your-token>"
```

## 📡 API Endpoints

### Health

| Method | Endpoint         | Description            |
| ------ | ---------------- | ---------------------- |
| GET    | `/health`        | Health check           |
| GET    | `/api/v1/health` | Versioned health check |

### Authentication

| Method | Endpoint                | Description      | Auth |
| ------ | ----------------------- | ---------------- | ---- |
| POST   | `/api/v1/auth/login`    | Login            | No   |
| POST   | `/api/v1/auth/register` | Register         | No   |
| POST   | `/api/v1/auth/refresh`  | Refresh token    | No   |
| GET    | `/api/v1/auth/me`       | Get current user | Yes  |

### Categories

| Method | Endpoint                 | Description | Auth          |
| ------ | ------------------------ | ----------- | ------------- |
| GET    | `/api/v1/categories`     | List all    | Yes           |
| GET    | `/api/v1/categories/:id` | Get by ID   | Yes           |
| POST   | `/api/v1/categories`     | Create      | Admin/Manager |
| PUT    | `/api/v1/categories/:id` | Update      | Admin/Manager |
| DELETE | `/api/v1/categories/:id` | Delete      | Admin         |

### Products

| Method | Endpoint                     | Description  | Auth          |
| ------ | ---------------------------- | ------------ | ------------- |
| GET    | `/api/v1/products`           | List all     | Yes           |
| GET    | `/api/v1/products/:id`       | Get by ID    | Yes           |
| POST   | `/api/v1/products`           | Create       | Admin/Manager |
| PUT    | `/api/v1/products/:id`       | Update       | Admin/Manager |
| DELETE | `/api/v1/products/:id`       | Delete       | Admin         |
| PATCH  | `/api/v1/products/:id/stock` | Update stock | Yes           |

### Customers

| Method | Endpoint                | Description | Auth          |
| ------ | ----------------------- | ----------- | ------------- |
| GET    | `/api/v1/customers`     | List all    | Yes           |
| GET    | `/api/v1/customers/:id` | Get by ID   | Yes           |
| POST   | `/api/v1/customers`     | Create      | Yes           |
| PUT    | `/api/v1/customers/:id` | Update      | Yes           |
| DELETE | `/api/v1/customers/:id` | Delete      | Admin/Manager |

### Transactions

| Method | Endpoint                          | Description   | Auth          |
| ------ | --------------------------------- | ------------- | ------------- |
| GET    | `/api/v1/transactions`            | List all      | Yes           |
| GET    | `/api/v1/transactions/:id`        | Get by ID     | Yes           |
| POST   | `/api/v1/transactions`            | Create sale   | Yes           |
| PATCH  | `/api/v1/transactions/:id/status` | Update status | Admin/Manager |

### Reports

| Method | Endpoint                        | Description   | Auth          |
| ------ | ------------------------------- | ------------- | ------------- |
| GET    | `/api/v1/reports/sales/daily`   | Daily sales   | Admin/Manager |
| GET    | `/api/v1/reports/sales/monthly` | Monthly sales | Admin/Manager |
| GET    | `/api/v1/reports/products/top`  | Top products  | Admin/Manager |

## 🔧 Configuration

| Variable         | Description                          | Default        |
| ---------------- | ------------------------------------ | -------------- |
| `APP_ENV`        | Environment (development/production) | development    |
| `APP_PORT`       | Server port                          | 8080           |
| `JWT_SECRET`     | JWT signing secret                   | (change this!) |
| `DB_PATH`        | SQLite database path                 | ./data/pos.db  |
| `RATE_LIMIT_RPS` | Requests per second limit            | 100            |

## 📝 Response Format

### Success

```json
{
  "success": true,
  "message": "Data retrieved successfully",
  "data": { ... },
  "meta": {
    "page": 1,
    "per_page": 10,
    "total": 100,
    "total_pages": 10
  }
}
```

### Error

```json
{
  "success": false,
  "message": "Validation failed",
  "errors": [{ "field": "email", "message": "Email is required" }]
}
```

## 🛠️ Make Commands

```bash
make help          # Show all commands
make dev           # Run development server
make build         # Build binary
make test          # Run tests
make seed          # Seed database
make docker-build  # Build Docker image
make docker-run    # Run with Docker Compose
```

## 📄 License

MIT License - see [LICENSE](LICENSE) for details.
