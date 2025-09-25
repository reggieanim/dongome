# Getting Started with Dongome

## Quick Demo Setup

### 1. Start Infrastructure
```bash
# Start PostgreSQL, Redis, and NATS
make run-infrastructure

# Or start everything with Docker
make run-docker
```

### 2. Run Database Migrations
```bash
make migrate-up
```

### 3. Test the Event Flow
```bash
# Run the NATS event flow demonstration
go run examples/event_flow_demo.go
```

### 4. Start the API Server
```bash
# In one terminal
make run-api
```

### 5. Start the Worker
```bash
# In another terminal
make run-worker
```

### 6. Test API Endpoints

**Health Check:**
```bash
curl http://localhost:8080/health
```

**Register a User:**
```bash
curl -X POST http://localhost:8080/api/v1/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john.doe@example.com",
    "password": "securepassword123",
    "first_name": "John",
    "last_name": "Doe"
  }'
```

**Login a User:**
```bash
curl -X POST http://localhost:8080/api/v1/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john.doe@example.com",
    "password": "securepassword123"
  }'
```

**Verify Email:**
```bash
curl -X POST http://localhost:8080/api/v1/users/verify-email \
  -H "Content-Type: application/json" \
  -d '{
    "token": "verification-token-from-registration"
  }'
```

### 7. Monitor Event Processing

Watch the console outputs from both the API server and worker to see:

1. **API Server**: Real-time event handling for user interactions
2. **Worker**: Background processing of events
3. **NATS**: Event flow between bounded contexts

### 8. View NATS Monitoring
```bash
# NATS monitoring dashboard
open http://localhost:8222
```

## Architecture Highlights Demonstrated

### ✅ Domain-Driven Design
- **Bounded Contexts**: Users, Listings, Transactions, Reviews, Notifications
- **Domain Events**: UserRegistered, UserEmailVerified, etc.
- **Aggregates**: User with SellerProfile
- **Value Objects**: Location, UserStatus, UserRole

### ✅ Hexagonal Architecture
- **Domain Layer**: Pure business logic (internal/users/domain)
- **Application Layer**: Use cases (internal/users/app)  
- **Infrastructure Layer**: External concerns (internal/users/infra)

### ✅ Event-Driven Architecture
- **Event Bus**: NATS with JetStream for reliability
- **Event Publishing**: Domain events published on state changes
- **Event Handling**: Cross-context communication via events
- **Event Persistence**: JetStream provides durability

### ✅ Clean Code Principles
- **Dependency Inversion**: Interfaces in domain, implementations in infra
- **Single Responsibility**: Each module has one clear purpose
- **Open/Closed**: Easy to extend without modifying existing code

### ✅ Modern Go Practices
- **Module Structure**: Following Go 1.22+ conventions
- **Error Handling**: Custom domain errors with HTTP mapping
- **Configuration Management**: Viper with environment overrides
- **Structured Logging**: Zap for production-ready logging
- **Database Migrations**: golang-migrate for schema versioning

### ✅ Production Readiness
- **Docker Support**: Multi-stage builds for API and Worker
- **Health Checks**: Monitoring endpoints for load balancers
- **Graceful Shutdown**: Proper cleanup on termination signals
- **Test Coverage**: Unit tests for domain logic
- **Build Automation**: Comprehensive Makefile

## Event Flow Example

When a user registers:

1. **HTTP Request** → API Handler
2. **Domain Logic** → User aggregate created
3. **Event Publishing** → UserRegistered event to NATS
4. **Event Distribution** → Multiple subscribers process:
   - **Notifications**: Send welcome email
   - **Analytics**: Record metrics
   - **Worker**: Background processing

This demonstrates true decoupling - the user registration doesn't need to know about notifications, analytics, or other concerns. They're all handled asynchronously via events.

## Next Steps

1. **Extend Domain Models**: Add Listing, Transaction, Review entities
2. **Implement Payment Integration**: Add MoMo payment processing
3. **Add Authentication**: JWT middleware for protected endpoints  
4. **Implement Search**: Full-text search for listings
5. **Add File Upload**: Image handling for listings
6. **Implement Caching**: Redis for frequently accessed data
7. **Add Rate Limiting**: Protect against abuse
8. **Implement Monitoring**: Prometheus metrics and tracing

This scaffold provides a solid foundation for building a production-ready marketplace application with modern Go practices! 🚀