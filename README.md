# Dongome - Marketplace Application

A modern marketplace application built with Go, following Domain-Driven Design (DDD) and Hexagonal Architecture principles. This is a Jiji.ng clone showcasing best practices for building scalable, maintainable Go applications.

## 🏗️ Architecture

This application follows a **modular monolith** architecture with clear separation of concerns:

- **Domain-Driven Design (DDD)**: Business logic is organized into bounded contexts
- **Hexagonal Architecture**: Clean separation between business logic and infrastructure
- **Event-Driven Architecture**: Loose coupling between modules via NATS messaging
- **CQRS patterns**: Clear separation of command and query responsibilities

## 📁 Project Structure

```
dongome/
├── cmd/                          # Application entry points
│   ├── api/                      # REST API server
│   └── worker/                   # Background worker
├── internal/                     # Private application code
│   ├── users/                    # User bounded context
│   │   ├── domain/               # Domain entities, value objects, events
│   │   ├── app/                  # Use cases, application services
│   │   └── infra/                # Repositories, HTTP handlers, external services
│   ├── listings/                 # Listings bounded context
│   ├── transactions/             # Transaction bounded context
│   ├── reviews/                  # Review bounded context
│   └── notifications/            # Notification bounded context
├── pkg/                          # Shared kernel - reusable packages
│   ├── config/                   # Configuration management
│   ├── logger/                   # Structured logging
│   ├── errors/                   # Domain error types
│   ├── events/                   # Event bus abstraction
│   └── db/                       # Database utilities
├── migrations/                   # Database migrations
├── docker/                       # Docker configurations
├── config/                       # Configuration files
└── Makefile                      # Build automation
```

## 🚀 Features

### Core Modules

1. **Users Module**
   - User registration and authentication
   - Email verification
   - Seller profile management and verification
   - Role-based access control

2. **Listings Module**
   - Product listing creation and management
   - Category hierarchies
   - Image uploads and management
   - Search and filtering
   - Favorites system

3. **Transactions Module**
   - Escrow payment system
   - Mobile Money (MoMo) integration
   - Order lifecycle management
   - Payment processing

4. **Reviews Module**
   - Seller ratings and feedback
   - Review management
   - Rating aggregation

5. **Notifications Module**
   - Email notifications
   - SMS notifications
   - Real-time notifications via NATS

## 🛠️ Tech Stack

- **Language**: Go 1.22+
- **Framework**: Gin (HTTP router)
- **Database**: PostgreSQL with GORM ORM
- **Cache**: Redis
- **Message Broker**: NATS with JetStream
- **Configuration**: Viper
- **Logging**: Zap
- **Migrations**: golang-migrate
- **Containerization**: Docker & Docker Compose

## 🏃‍♂️ Quick Start

### Prerequisites

- Go 1.22+
- Docker and Docker Compose
- Make

### 1. Clone the Repository

```bash
git clone <repository-url>
cd dongome
```

### 2. Install Development Tools

```bash
make install-tools
```

### 3. Start Infrastructure Services

```bash
make run-infrastructure
```

This will start PostgreSQL, Redis, and NATS using Docker Compose.

### 4. Run Database Migrations

```bash
make migrate-up
```

### 5. Start the API Server

```bash
make run-api
```

The API will be available at `http://localhost:8080`

### 6. Start the Worker (Optional)

In another terminal:

```bash
make run-worker
```

## 🐳 Docker Development

### Start Everything with Docker

```bash
make run-docker
```

This starts all services including the API, worker, database, cache, and message broker.

### Stop All Services

```bash
make stop-docker
```

### View Logs

```bash
make logs-docker
```

## 📊 API Endpoints

### Health Check
```
GET /health
```

### User Management
```
POST   /api/v1/users/register          # Register new user
POST   /api/v1/users/login             # Login user
POST   /api/v1/users/verify-email      # Verify email
POST   /api/v1/users/{id}/upgrade-to-seller  # Upgrade to seller
GET    /api/v1/users/{id}              # Get user profile
```

## 🏗️ Development Workflow

### Running Tests

```bash
make test                    # Run all tests
make test-coverage          # Run tests with coverage report
```

### Code Quality

```bash
make lint                   # Run linter
make format                 # Format code
make vet                    # Run go vet
```

### Database Operations

```bash
make migrate-up             # Run migrations
make migrate-down           # Rollback migrations
make migrate-create name=add_new_table  # Create new migration
make db-reset              # Reset database
```

### Building

```bash
make build                 # Build all binaries
make build-api             # Build API only
make build-worker          # Build worker only
```

## 🔄 Event-Driven Architecture

The application uses NATS for event-driven communication between bounded contexts. Here's an example flow:

### UserRegistered Event Flow

1. **User Registration**: A user registers via the API
2. **Domain Event**: `UserRegistered` event is published to NATS
3. **Event Handlers**: 
   - **API Server**: Sends welcome email
   - **Worker**: Processes background tasks (analytics, etc.)
   - **Notifications**: Queues verification email

```go
// Example event publishing
event, _ := events.NewEvent(
    domain.UserRegisteredEvent,
    user.ID,
    domain.UserRegistered{
        UserID:    user.ID,
        Email:     user.Email,
        FirstName: user.FirstName,
        LastName:  user.LastName,
        // ... other fields
    },
)
eventBus.Publish(ctx, event)
```

```go
// Example event handling
func handleUserRegistered(ctx context.Context, event *events.Event) error {
    var userData domain.UserRegistered
    events.ParseEventData(event, &userData)
    
    // Process the event
    // - Send welcome email
    // - Create user analytics record
    // - Initialize user preferences
    
    return nil
}
```

## 🔧 Configuration

Configuration is managed through:
1. YAML files in `/config` directory
2. Environment variables (override YAML)
3. Command-line flags (highest priority)

### Environment Variables

Key environment variables:

```bash
# Server
PORT=8080
SERVER_HOST=localhost

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=dongome
DB_PASSWORD=password
DB_NAME=dongome_db

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379

# NATS
NATS_URL=nats://localhost:4222

# JWT
JWT_SECRET=your-secret-key

# MoMo Integration
MOMO_API_KEY=your-api-key
MOMO_API_SECRET=your-api-secret
```

## 🏛️ Domain-Driven Design

### Bounded Contexts

1. **Users Context**: User management, authentication, seller profiles
2. **Listings Context**: Product listings, categories, search
3. **Transactions Context**: Payments, orders, escrow
4. **Reviews Context**: Ratings, feedback
5. **Notifications Context**: Email, SMS, real-time notifications

### Domain Events

- `UserRegistered`: New user account created
- `UserEmailVerified`: User verified their email
- `UserUpgradedToSeller`: User became a seller
- `ListingCreated`: New listing published
- `OrderPlaced`: New order created
- `PaymentCompleted`: Payment processed successfully

## 🚀 Deployment

### Production Build

```bash
make build
make docker-build
```

### Environment-Specific Configs

Create environment-specific config files:
- `config/config.yaml` (default)
- `config/config.staging.yaml`
- `config/config.production.yaml`

## 🧪 Testing

The application includes comprehensive tests:

- **Unit Tests**: Domain logic and business rules
- **Integration Tests**: Database and external service integration
- **API Tests**: HTTP endpoint testing
- **Event Tests**: NATS event flow testing

## 📈 Monitoring & Observability

- **Structured Logging**: Zap logger with JSON output
- **Health Checks**: `/health` endpoint for load balancer
- **NATS Monitoring**: Available at `http://localhost:8222`
- **Metrics**: Ready for Prometheus integration

## 🤝 Contributing

1. Follow Go best practices and project conventions
2. Write tests for new features
3. Update documentation for API changes
4. Follow the existing code structure
5. Use meaningful commit messages

## 📝 License

[Add your license here]

## 🙏 Acknowledgments

- Inspired by modern Go application architecture patterns
- Built with best practices from the Go community
- Follows principles from Domain-Driven Design and Clean Architecture

---

**Happy coding! 🚀**
