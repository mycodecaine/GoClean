# GoClean - Clean Architecture Go API

A comprehensive Golang project implementing Clean Architecture, Domain-Driven Design (DDD), CQRS, REST API, gRPC, and modern development practices.

## 🏗️ Architecture

This project follows **Clean Architecture** principles with **Domain-Driven Design (DDD)** patterns:

```
├── cmd/                    # Application entry points
│   ├── http/              # HTTP server main
│   └── grpc/              # gRPC server main
├── internal/              # Private application code
│   ├── domain/            # Domain layer (entities, repositories, services)
│   ├── application/       # Application layer (commands, queries, DTOs)
│   ├── infrastructure/    # Infrastructure layer (database, cache, auth)
│   └── interfaces/        # Interface layer (HTTP handlers, gRPC handlers)
├── pkg/                   # Public packages
│   ├── config/           # Configuration management
│   └── logger/           # Logging utilities
├── api/                  # API definitions
│   ├── proto/           # Protocol buffer definitions
│   └── swagger/         # Swagger documentation
├── test/                # Test utilities and fixtures
└── deployments/         # Docker and deployment configurations
```

## 🚀 Features

- **Clean Architecture** with clear separation of concerns
- **Domain-Driven Design (DDD)** with aggregates and domain services
- **Aggregate Root Pattern** with domain events system
- **Soft Delete Pattern** with restore capabilities (Entity Framework style)
- **CQRS (Command Query Responsibility Segregation)** pattern
- **Domain Events** with event dispatcher and handlers
- **REST API** with Echo framework
- **gRPC API** with Protocol Buffers
- **JWT Authentication** with Keycloak integration
- **Database** with PostgreSQL and GORM
- **Caching** with Redis
- **Docker Compose** for local development
- **Swagger** API documentation
- **Comprehensive logging** with structured logging
- **Graceful shutdown** and error handling

## 🛠️ Technologies

### Core
- **Go 1.21+** - Programming language
- **Echo v4** - HTTP web framework
- **gRPC** - High-performance RPC framework
- **GORM** - ORM library for Go

### Infrastructure
- **PostgreSQL** - Primary database
- **Redis** - Caching and session storage
- **Keycloak** - Identity and access management

### DevOps
- **Docker** & **Docker Compose** - Containerization
- **Swagger/OpenAPI** - API documentation

## 📋 Prerequisites

- Go 1.21 or higher
- Docker and Docker Compose
- Protocol Buffer Compiler (protoc) - for gRPC development
- Make (optional, for build automation)

## 🚀 Quick Start

### 1. Clone the Repository

```bash
git clone <repository-url>
cd GoClean
```

### 2. Start Infrastructure Services

```bash
docker-compose up -d postgres redis keycloak
```

Wait for all services to be healthy. You can check status with:

```bash
docker-compose ps
```

### 3. Set Environment Variables

Create a `.env` file or set environment variables:

```bash
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=goclean
DB_SSL_MODE=disable

# Redis
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
REDIS_DB=0

# Keycloak
KEYCLOAK_BASE_URL=http://localhost:8081
KEYCLOAK_REALM=goclean
KEYCLOAK_CLIENT_ID=goclean-api
KEYCLOAK_CLIENT_SECRET=your-client-secret

# Server
HTTP_HOST=localhost
HTTP_PORT=8080
GRPC_HOST=localhost
GRPC_PORT=9090

# Application
APP_NAME=GoClean
APP_VERSION=1.0.0
APP_ENV=development
LOG_LEVEL=info
```

### 4. Install Dependencies

```bash
go mod tidy
```

### 5. Generate Protocol Buffers (if needed)

```bash
# Install protoc-gen-go and protoc-gen-go-grpc if not installed
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Generate protobuf code
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    api/proto/goclean.proto
```

### 6. Run the Application

```bash
# Run tests (all passing!)
go test ./test -v

# Start the HTTP server with full soft delete support
go run cmd/http/main.go

# Start the gRPC server
go run cmd/grpc/main.go

# Build production binaries
go build -o bin/goclean ./cmd/http
go build -o bin/grpc-server ./cmd/grpc
```

## 📖 API Documentation

### REST API

Once the HTTP server is running, you can access:

- **Swagger UI**: http://localhost:8080/swagger/index.html
- **Health Check**: http://localhost:8080/health

### gRPC API

The gRPC server runs on `localhost:9090` by default. Use tools like:
- **grpcurl** for command-line testing
- **Postman** (with gRPC support)
- **BloomRPC** or **Evans** for GUI clients

## 🔐 Authentication

The application uses **Keycloak** for authentication:

1. **Keycloak Admin Console**: http://localhost:8081
   - Username: `admin`
   - Password: `admin`

2. **Default Test Users**:
   - **Admin**: `admin@goclean.com` / `admin123`
   - **User**: `test@goclean.com` / `test123`

3. **Getting JWT Token**:
```bash
curl -X POST http://localhost:8081/realms/goclean/protocol/openid-connect/token \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "grant_type=password" \
  -d "client_id=goclean-api" \
  -d "username=admin@goclean.com" \
  -d "password=admin123"
```

## 🧪 Testing

### Unit Tests

```bash
# Run all tests (includes domain services, aggregate root, and soft delete tests)
go test ./test -v

# Run specific domain tests
go test ./internal/domain/... -v

# Run all tests recursively
go test ./...
```

### Integration Tests

```bash
go test -tags=integration ./test/...
```

### API Testing

Use the provided test users or create new ones through Keycloak.

## 🐳 Docker

### Build Application Image

```bash
docker build -t goclean-api .
```

### Run Full Stack

```bash
docker-compose up
```

This starts:
- PostgreSQL database
- Redis cache  
- Keycloak identity server
- GoClean API (if uncommented in docker-compose.yml)

## 📚 Project Structure Details

### Domain Layer
- **Entities**: Core business objects (User, Product, Order)
- **Value Objects**: OrderStatus and other value types
- **Repository Interfaces**: Data access contracts
- **Domain Services**: Business logic that doesn't belong to entities

### Application Layer
- **Commands**: Write operations (CreateUser, CreateProduct, etc.)
- **Queries**: Read operations (GetUser, ListProducts, etc.)
- **DTOs**: Data transfer objects for API contracts
- **Handlers**: Command and query handlers implementing CQRS

### Infrastructure Layer
- **Persistence**: GORM implementations of repositories
- **Cache**: Redis client and caching logic
- **Auth**: Keycloak integration and JWT handling

### Interface Layer
- **HTTP**: REST API handlers and middleware
- **gRPC**: Protocol buffer service implementations

## 🎯 Enhanced Domain Features

### Aggregate Root Pattern
The application implements the **Aggregate Root** pattern with domain events:

- **BaseEntity**: Foundation for all entities with soft delete capabilities
- **AggregateRoot**: Manages domain events and entity lifecycle
- **Domain Events**: Event-driven architecture for business logic
- **Factory Methods**: `NewUser()`, `NewProduct()`, `NewOrder()` for proper entity creation

### Soft Delete Pattern (Entity Framework Style)
Full soft delete implementation with:

- **Soft Delete**: Mark entities as deleted without physical removal
- **Restore Operations**: Bring back soft-deleted entities
- **Query Methods**: 
  - `GetByIDIncludeDeleted()` - Get entities including deleted ones
  - `ListIncludeDeleted()` - List all entities including deleted
  - `ListDeleted()` - List only deleted entities
  - `SoftDelete()` - Mark entity as deleted
  - `Restore()` - Restore deleted entity

### REST API Soft Delete Endpoints
```bash
# Soft delete a user
DELETE /users/{id}/soft-delete

# Restore a deleted user  
POST /users/{id}/restore

# List deleted users
GET /users/deleted

# Get user including deleted ones
GET /users/{id}?include_deleted=true
```

## 🔧 Configuration

The application supports environment-based configuration:

- **Development**: Uses `.env` file or environment variables
- **Production**: Uses environment variables only
- **Docker**: Uses docker-compose environment settings

## 🚀 Deployment

### Environment Preparation

1. Set up PostgreSQL database
2. Set up Redis instance
3. Configure Keycloak realm and clients
4. Set environment variables

### Application Deployment

```bash
# Build binaries
go build -o bin/goclean ./cmd/http
go build -o bin/grpc-server ./cmd/grpc

# Run HTTP server
./bin/goclean

# Run gRPC server  
./bin/grpc-server
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Clean Architecture by Robert C. Martin
- Domain-Driven Design by Eric Evans
- CQRS pattern and Event Sourcing concepts
- Echo framework community
- gRPC and Protocol Buffers
- Keycloak community

## 📞 Support

For support and questions:
- Create an issue in the repository
- Check the documentation in `/docs`
- Review the API documentation at `/swagger`