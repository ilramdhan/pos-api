# GoPOS API - Point of Sales Backend

A modern, scalable REST API backend for Point of Sales (POS) systems built with Go, PostgreSQL/Supabase, and best practices.

![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-Supabase-4169E1?style=flat&logo=postgresql)
![License](https://img.shields.io/badge/License-MIT-green.svg)

## ✨ Features

- **RESTful API** with versioning (`/api/v1/`)
- **JWT Authentication** with role-based access control (Admin, Manager, Cashier)
- **PostgreSQL/Supabase Database** for production-ready persistence
- **Rate Limiting** per IP address
- **Standardized Responses** with validation errors
- **Docker Ready** with multi-stage builds
- **Health Checks** for monitoring
- **Sales Reports** (daily, monthly, top products)
- **POS Features** (transactions, customers, products, categories)

## 🏗️ Architecture

```
├── cmd/api/              # Application entry point
├── internal/
│   ├── config/           # Configuration management (Viper)
│   ├── database/         # PostgreSQL connection & migrations
│   ├── dto/              # Data Transfer Objects
│   ├── handler/          # HTTP handlers
│   ├── middleware/       # Auth, CORS, Rate Limit, Logger
│   ├── models/           # Domain models
│   ├── repository/       # Data access layer (PostgreSQL)
│   ├── router/           # Route definitions
│   ├── service/          # Business logic
│   └── utils/            # Helpers (JWT, Response, Validation)
├── scripts/              # Database seed script
├── tests/                # Integration tests
└── docs/                 # API documentation (Swagger)
```

## 🚀 Quick Start

### Prerequisites

- Go 1.21 or higher
- PostgreSQL database (Supabase recommended)
- Docker (optional)

### 1. Supabase Setup

1. **Create a Supabase Project**
   - Go to [supabase.com](https://supabase.com) and create a new project
   - Note your project password (you'll need it)

2. **Get Connection String**
   - Navigate to: Settings → Database → Connection string
   - Copy the "Transaction pooler" connection string
   - Replace `[YOUR-PASSWORD]` with your database password

3. **Run Database Migration**
   - Go to SQL Editor in Supabase dashboard
   - Copy contents of `internal/database/migrations/001_init_postgres.sql`
   - Execute the SQL to create all tables

4. **Seed Initial Data (Optional)**
   - Copy contents of `scripts/seed_postgres.sql`
   - Execute in SQL Editor to add test data

### 2. Local Development

```bash
# Clone repository
git clone <repository-url>
cd goland-dasar
go mod tidy

# Configure environment
cp .env.example .env
# Edit .env and set your DB_CONN (Supabase connection string)

# Run server
make dev
# Or: go run ./cmd/api/main.go

# Test health check
curl http://localhost:8080/health
```

### 3. Docker Deployment

```bash
# Build and run
docker-compose up --build

# Stop
docker-compose down
```

## 🔧 Configuration

| Variable                   | Description                           | Default                 |
| -------------------------- | ------------------------------------- | ----------------------- |
| `APP_ENV`                  | Environment (development/production)  | development             |
| `APP_PORT`                 | Server port                           | 8080                    |
| `DB_CONN`                  | PostgreSQL/Supabase connection string | (required)              |
| `JWT_SECRET`               | JWT signing secret                    | (change in production!) |
| `JWT_EXPIRY_HOURS`         | Access token expiry                   | 24                      |
| `JWT_REFRESH_EXPIRY_HOURS` | Refresh token expiry                  | 168 (7 days)            |
| `RATE_LIMIT_RPS`           | Requests per second limit             | 100                     |
| `CORS_ALLOWED_ORIGINS`     | Allowed CORS origins                  | http://localhost:3000   |

### Example `.env` Configuration

```env
APP_ENV=development
APP_PORT=8080
DB_CONN=postgres://postgres.projectid:password@aws-0-region.pooler.supabase.com:6543/postgres?sslmode=require
JWT_SECRET=your-super-secret-jwt-key
CORS_ALLOWED_ORIGINS=http://localhost:3000,http://localhost:5173
```

## 🔐 Authentication

### Default Test Credentials

| Email               | Password    | Role    |
| ------------------- | ----------- | ------- |
| admin@gopos.local   | Admin123!   | admin   |
| manager@gopos.local | Manager123! | manager |
| cashier@gopos.local | Cashier123! | cashier |

### Login Example

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "admin@gopos.local", "password": "Admin123!"}'
```

### Using Token (Cookie-based)

After login, authentication cookies are automatically sent. For API testing:

```bash
curl http://localhost:8080/api/v1/products \
  --cookie "access_token=<your-token>"
```

## 📡 API Endpoints

### Health

| Method | Endpoint  | Description  |
| ------ | --------- | ------------ |
| GET    | `/health` | Health check |

### Authentication

| Method | Endpoint                | Description      | Auth |
| ------ | ----------------------- | ---------------- | ---- |
| POST   | `/api/v1/auth/login`    | Login            | No   |
| POST   | `/api/v1/auth/register` | Register         | No   |
| POST   | `/api/v1/auth/refresh`  | Refresh token    | No   |
| POST   | `/api/v1/auth/logout`   | Logout           | No   |
| GET    | `/api/v1/auth/me`       | Get current user | Yes  |
| PUT    | `/api/v1/auth/me`       | Update profile   | Yes  |

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

### Users (Admin Only)

| Method | Endpoint            | Description | Auth  |
| ------ | ------------------- | ----------- | ----- |
| GET    | `/api/v1/users`     | List all    | Admin |
| GET    | `/api/v1/users/:id` | Get by ID   | Admin |
| POST   | `/api/v1/users`     | Create      | Admin |
| PUT    | `/api/v1/users/:id` | Update      | Admin |
| DELETE | `/api/v1/users/:id` | Delete      | Admin |

### Reports

| Method | Endpoint                        | Description   | Auth          |
| ------ | ------------------------------- | ------------- | ------------- |
| GET    | `/api/v1/reports/sales/daily`   | Daily sales   | Admin/Manager |
| GET    | `/api/v1/reports/sales/monthly` | Monthly sales | Admin/Manager |
| GET    | `/api/v1/reports/products/top`  | Top products  | Admin/Manager |

### Dashboard

| Method | Endpoint            | Description       | Auth          |
| ------ | ------------------- | ----------------- | ------------- |
| GET    | `/api/v1/dashboard` | Dashboard summary | Admin/Manager |

### POS

| Method | Endpoint               | Description      | Auth |
| ------ | ---------------------- | ---------------- | ---- |
| GET    | `/api/v1/pos/products` | POS product list | Yes  |
| POST   | `/api/v1/pos/checkout` | Checkout         | Yes  |
| POST   | `/api/v1/pos/hold`     | Hold transaction | Yes  |
| GET    | `/api/v1/pos/held`     | Get held items   | Yes  |
| DELETE | `/api/v1/pos/held/:id` | Delete held item | Yes  |

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

## 🧪 Testing

### Running Integration Tests

```bash
# Set up test database (use a separate Supabase project or local PostgreSQL)
export TEST_DB_CONN="postgres://user:pass@host:5432/pos_test?sslmode=disable"

# Run tests
go test -v ./tests/...

# Run specific test
go test -v ./tests/... -run TestAuthLogin
```

### Test Coverage

```bash
go test -cover ./...
```

## 🛠️ Make Commands

```bash
make help          # Show all commands
make dev           # Run development server
make build         # Build binary
make test          # Run tests
make docker-build  # Build Docker image
make docker-run    # Run with Docker Compose
```

## 📚 API Documentation

Swagger documentation is available at:

- Development: `http://localhost:8080/docs/swagger.yaml`
- See `docs/swagger.yaml` for OpenAPI specification

## 🚀 Deployment

### Supabase + VPS/Cloud

1. **Set up Supabase database** (see Quick Start section)
2. **Build the binary**:
   ```bash
   CGO_ENABLED=0 GOOS=linux go build -o gopos-api ./cmd/api/
   ```
3. **Configure environment** on your server
4. **Run with systemd** or container orchestration

### Environment Variables for Production

```env
APP_ENV=production
APP_PORT=8080
DB_CONN=<your-supabase-connection-string>
JWT_SECRET=<strong-random-secret>
CORS_ALLOWED_ORIGINS=https://your-frontend-domain.com
```

## 📄 License

MIT License - see [LICENSE](LICENSE) for details.
