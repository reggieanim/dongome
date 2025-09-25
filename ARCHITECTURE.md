# Dongome Marketplace - Complete System Architecture

## 🏗️ High-Level Architecture Overview

```mermaid
graph TB
    subgraph "External Systems"
        Client[Web/Mobile Client]
        MoMo[MoMo Payment API]
        SMS[SMS Gateway]
        Email[Email Service]
        FileStorage[File Storage S3]
    end

    subgraph "Load Balancer & Reverse Proxy"
        LB[Nginx/HAProxy]
    end

    subgraph "Modular Monolith - Dongome API"
        subgraph "HTTP Layer"
            Router[Gin HTTP Router]
            Middleware[Middleware Stack]
        end

        subgraph "Bounded Contexts"
            Users[👥 Users Context]
            Listings[🏪 Listings Context]
            Transactions[💰 Transactions Context]
            Reviews[⭐ Reviews Context]
            Notifications[📧 Notifications Context]
        end

        subgraph "Shared Kernel"
            Events[Event Bus NATS]
            Config[Configuration]
            Logger[Logging]
            Errors[Error Handling]
            DB[Database Utilities]
        end
    end

    subgraph "Background Worker"
        Worker[Event Worker Process]
        Jobs[Background Jobs]
    end

    subgraph "Infrastructure"
        Postgres[(PostgreSQL)]
        Redis[(Redis Cache)]
        NATS[NATS JetStream]
    end

    subgraph "Monitoring & Observability"
        Metrics[Prometheus Metrics]
        Tracing[Jaeger Tracing]
        Logs[Centralized Logging]
    end

    Client --> LB
    LB --> Router
    Router --> Middleware
    Middleware --> Users
    Middleware --> Listings
    Middleware --> Transactions
    Middleware --> Reviews
    Middleware --> Notifications

    Users --> Events
    Listings --> Events
    Transactions --> Events
    Reviews --> Events
    Notifications --> Events

    Events --> NATS
    Events --> Worker
    Worker --> Jobs

    Users --> DB
    Listings --> DB
    Transactions --> DB
    Reviews --> DB
    Notifications --> DB

    DB --> Postgres
    DB --> Redis

    Transactions --> MoMo
    Notifications --> SMS
    Notifications --> Email
    Listings --> FileStorage

    Router --> Metrics
    Middleware --> Tracing
    Logger --> Logs
```

## 🧩 Modular Monolith Deep Dive

### Domain-Driven Design Structure

```mermaid
graph TB
    subgraph "Modular Monolith Architecture"
        subgraph "Presentation Layer"
            HTTP[HTTP Handlers]
            Middleware[Middleware Stack]
            Validation[Request Validation]
        end

        subgraph "Application Layer"
            AppServices[Application Services]
            Commands[Command Handlers]
            Queries[Query Handlers]
            DTOs[Data Transfer Objects]
        end

        subgraph "Domain Layer"
            Aggregates[Aggregate Roots]
            Entities[Domain Entities]
            ValueObjects[Value Objects]
            DomainServices[Domain Services]
            DomainEvents[Domain Events]
            BusinessRules[Business Rules]
        end

        subgraph "Infrastructure Layer"
            Repositories[Repository Implementations]
            EventPublishers[Event Publishers]
            ExternalAPIs[External API Clients]
            Persistence[Database Access]
        end

        subgraph "Shared Kernel"
            CommonTypes[Common Types]
            Utilities[Shared Utilities]
            EventBus[Event Bus Abstraction]
            ErrorHandling[Error Handling]
        end
    end

    HTTP --> AppServices
    AppServices --> Commands
    AppServices --> Queries
    Commands --> Aggregates
    Queries --> Repositories
    Aggregates --> DomainEvents
    DomainEvents --> EventPublishers
    EventPublishers --> EventBus
    Repositories --> Persistence
    
    Aggregates --> CommonTypes
    AppServices --> Utilities
    Infrastructure --> ErrorHandling
```

## 📦 Bounded Context Details

### 1. Users Bounded Context

```mermaid
graph TB
    subgraph "Users Bounded Context"
        subgraph "Domain Layer"
            UserAggregate[User Aggregate Root]
            SellerProfile[Seller Profile Entity]
            UserEvents[Domain Events]
            UserRules[Business Rules]
        end

        subgraph "Application Layer"
            UserService[User Service]
            RegisterCommand[Register User Command]
            LoginCommand[Login Command]
            VerifyEmailCommand[Verify Email Command]
            UpgradeSellerCommand[Upgrade to Seller Command]
        end

        subgraph "Infrastructure Layer"
            UserRepository[User Repository GORM]
            UserHandlers[HTTP Handlers]
            UserEventPublisher[Event Publisher]
        end
    end

    subgraph "Domain Rules & Events"
        EmailUnique[Email Must Be Unique]
        PasswordPolicy[Password Policy]
        EmailVerification[Email Verification Required]
        SellerUpgrade[Seller Upgrade Process]
        
        UserRegistered[UserRegistered Event]
        UserEmailVerified[UserEmailVerified Event]
        UserUpgradedToSeller[UserUpgradedToSeller Event]
        UserLoggedIn[UserLoggedIn Event]
    end

    UserAggregate --> EmailUnique
    UserAggregate --> PasswordPolicy
    UserAggregate --> EmailVerification
    UserAggregate --> SellerUpgrade

    RegisterCommand --> UserRegistered
    VerifyEmailCommand --> UserEmailVerified
    UpgradeSellerCommand --> UserUpgradedToSeller
    LoginCommand --> UserLoggedIn

    UserService --> UserRepository
    UserService --> UserEventPublisher
    UserHandlers --> UserService
```

### 2. Listings Bounded Context

```mermaid
graph TB
    subgraph "Listings Bounded Context"
        subgraph "Domain Layer"
            ListingAggregate[Listing Aggregate Root]
            Category[Category Entity]
            Image[Image Value Object]
            Location[Location Value Object]
            Price[Price Value Object]
            ListingEvents[Domain Events]
        end

        subgraph "Application Layer"
            ListingService[Listing Service]
            CreateListingCommand[Create Listing Command]
            UpdateListingCommand[Update Listing Command]
            SearchListingsQuery[Search Listings Query]
            FavoriteListingCommand[Favorite Listing Command]
        end

        subgraph "Infrastructure Layer"
            ListingRepository[Listing Repository]
            CategoryRepository[Category Repository]
            SearchService[Search Service ElasticSearch]
            ImageUploadService[Image Upload Service]
            ListingHandlers[HTTP Handlers]
        end
    end

    subgraph "Domain Rules & Events"
        OnlyVerifiedSellers[Only Verified Sellers Can Create Premium Listings]
        PriceValidation[Price Must Be Positive]
        ImageLimits[Maximum 10 Images Per Listing]
        CategoryValidation[Category Must Be Active]
        
        ListingCreated[ListingCreated Event]
        ListingUpdated[ListingUpdated Event]
        ListingFavorited[ListingFavorited Event]
        ListingExpired[ListingExpired Event]
    end

    ListingAggregate --> OnlyVerifiedSellers
    ListingAggregate --> PriceValidation
    ListingAggregate --> ImageLimits
    ListingAggregate --> CategoryValidation

    CreateListingCommand --> ListingCreated
    UpdateListingCommand --> ListingUpdated
    FavoriteListingCommand --> ListingFavorited
```

### 3. Transactions Bounded Context

```mermaid
graph TB
    subgraph "Transactions Bounded Context"
        subgraph "Domain Layer"
            TransactionAggregate[Transaction Aggregate Root]
            OrderEntity[Order Entity]
            PaymentEntity[Payment Entity]
            EscrowEntity[Escrow Entity]
            DeliveryEntity[Delivery Entity]
            TransactionEvents[Domain Events]
        end

        subgraph "Application Layer"
            TransactionService[Transaction Service]
            CreateOrderCommand[Create Order Command]
            InitiatePaymentCommand[Initiate Payment Command]
            ConfirmDeliveryCommand[Confirm Delivery Command]
            ProcessRefundCommand[Process Refund Command]
        end

        subgraph "Infrastructure Layer"
            TransactionRepository[Transaction Repository]
            PaymentGateway[MoMo Payment Gateway]
            EscrowService[Escrow Service]
            NotificationService[Notification Service]
            TransactionHandlers[HTTP Handlers]
        end
    end

    subgraph "Domain Rules & Events"
        EscrowRules[Escrow Business Rules]
        PaymentValidation[Payment Validation]
        RefundPolicy[Refund Policy]
        DeliveryConfirmation[Delivery Confirmation Rules]
        
        OrderCreated[OrderCreated Event]
        PaymentInitiated[PaymentInitiated Event]
        PaymentCompleted[PaymentCompleted Event]
        DeliveryConfirmed[DeliveryConfirmed Event]
        TransactionCompleted[TransactionCompleted Event]
        RefundProcessed[RefundProcessed Event]
    end

    TransactionAggregate --> EscrowRules
    TransactionAggregate --> PaymentValidation
    TransactionAggregate --> RefundPolicy
    TransactionAggregate --> DeliveryConfirmation
```

### 4. Reviews Bounded Context

```mermaid
graph TB
    subgraph "Reviews Bounded Context"
        subgraph "Domain Layer"
            ReviewAggregate[Review Aggregate Root]
            Rating[Rating Value Object]
            Comment[Comment Value Object]
            ReviewEvents[Domain Events]
        end

        subgraph "Application Layer"
            ReviewService[Review Service]
            CreateReviewCommand[Create Review Command]
            UpdateReviewCommand[Update Review Command]
            GetSellerReviewsQuery[Get Seller Reviews Query]
        end

        subgraph "Infrastructure Layer"
            ReviewRepository[Review Repository]
            ReviewHandlers[HTTP Handlers]
        end
    end

    subgraph "Domain Rules & Events"
        OneReviewPerTransaction[One Review Per Transaction]
        RatingRange[Rating Must Be 1-5]
        ReviewPeriod[Review Within 30 Days]
        
        ReviewCreated[ReviewCreated Event]
        ReviewUpdated[ReviewUpdated Event]
        SellerRatingUpdated[SellerRatingUpdated Event]
    end
```

### 5. Notifications Bounded Context

```mermaid
graph TB
    subgraph "Notifications Bounded Context"
        subgraph "Domain Layer"
            NotificationAggregate[Notification Aggregate Root]
            Template[Template Value Object]
            Channel[Channel Value Object]
            NotificationEvents[Domain Events]
        end

        subgraph "Application Layer"
            NotificationService[Notification Service]
            SendEmailCommand[Send Email Command]
            SendSMSCommand[Send SMS Command]
            SendPushCommand[Send Push Notification Command]
        end

        subgraph "Infrastructure Layer"
            NotificationRepository[Notification Repository]
            EmailService[Email Service]
            SMSService[SMS Service]
            PushService[Push Notification Service]
            NotificationHandlers[HTTP Handlers]
        end
    end

    subgraph "External Event Handlers"
        UserRegisteredHandler[UserRegistered Event Handler]
        OrderCreatedHandler[OrderCreated Event Handler]
        PaymentCompletedHandler[PaymentCompleted Event Handler]
        ReviewCreatedHandler[ReviewCreated Event Handler]
    end
```

## 🔄 Event-Driven Architecture Flow

```mermaid
sequenceDiagram
    participant Client
    participant API
    participant Users
    participant NATS
    participant Worker
    participant Notifications
    participant External

    Client->>API: POST /users/register
    API->>Users: RegisterUser(email, password)
    Users->>Users: Create User Aggregate
    Users->>Users: Generate Verification Token
    Users->>NATS: Publish UserRegistered Event
    Users->>API: Return User + Verification Token
    API->>Client: 201 Created + User Data

    NATS->>Worker: UserRegistered Event
    Worker->>Notifications: Send Welcome Email
    Notifications->>External: Email Service API
    External-->>Notifications: Email Sent
    Notifications->>NATS: Publish EmailSent Event

    Client->>API: POST /users/verify-email
    API->>Users: VerifyEmail(token)
    Users->>Users: Mark Email as Verified
    Users->>NATS: Publish UserEmailVerified Event
    Users->>API: Return Success
    API->>Client: 200 OK

    NATS->>Worker: UserEmailVerified Event
    Worker->>Notifications: Send Verification Success Email
```

## 🏗️ Infrastructure Components

### Database Schema Design

```mermaid
erDiagram
    users {
        uuid id PK
        string email UK
        string password_hash
        string first_name
        string last_name
        string phone_number UK
        string avatar
        enum status
        enum role
        boolean email_verified
        boolean phone_verified
        string verification_token
        timestamp last_login_at
        timestamp created_at
        timestamp updated_at
    }

    seller_profiles {
        uuid id PK
        uuid user_id FK
        string business_name
        string business_address
        string business_phone
        string business_email
        string tax_number
        enum verification_status
        string verification_notes
        decimal rating
        integer total_reviews
        timestamp created_at
        timestamp updated_at
    }

    categories {
        uuid id PK
        string name
        string description
        uuid parent_id FK
        string image_url
        boolean is_active
        timestamp created_at
        timestamp updated_at
    }

    listings {
        uuid id PK
        uuid seller_id FK
        uuid category_id FK
        string title
        text description
        decimal price
        enum condition
        enum status
        json location
        json images
        json attributes
        integer views_count
        integer favorites_count
        timestamp expires_at
        timestamp created_at
        timestamp updated_at
    }

    transactions {
        uuid id PK
        uuid listing_id FK
        uuid buyer_id FK
        uuid seller_id FK
        decimal amount
        enum status
        enum payment_status
        string payment_reference
        json delivery_address
        timestamp delivery_confirmed_at
        timestamp created_at
        timestamp updated_at
    }

    reviews {
        uuid id PK
        uuid transaction_id FK
        uuid reviewer_id FK
        uuid reviewee_id FK
        integer rating
        text comment
        timestamp created_at
        timestamp updated_at
    }

    notifications {
        uuid id PK
        uuid user_id FK
        string type
        string channel
        string title
        text content
        json data
        boolean read
        timestamp sent_at
        timestamp created_at
    }

    users ||--o| seller_profiles : has
    users ||--o{ listings : creates
    listings }o--|| categories : belongs_to
    listings ||--o{ transactions : involves
    transactions ||--|| reviews : generates
    users ||--o{ notifications : receives
```

### Deployment Architecture

```mermaid
graph TB
    subgraph "Production Environment"
        subgraph "Load Balancer"
            ALB[Application Load Balancer]
        end

        subgraph "API Instances"
            API1[API Instance 1]
            API2[API Instance 2]
            API3[API Instance 3]
        end

        subgraph "Worker Instances"
            Worker1[Worker Instance 1]
            Worker2[Worker Instance 2]
        end

        subgraph "Database Cluster"
            PGMaster[(PostgreSQL Master)]
            PGSlave1[(PostgreSQL Replica 1)]
            PGSlave2[(PostgreSQL Replica 2)]
        end

        subgraph "Cache Cluster"
            RedisCluster[(Redis Cluster)]
        end

        subgraph "Message Queue"
            NATSCluster[NATS Cluster]
        end

        subgraph "External Services"
            S3[AWS S3 Storage]
            SES[AWS SES Email]
            SNS[AWS SNS SMS]
        end
    end

    ALB --> API1
    ALB --> API2
    ALB --> API3

    API1 --> PGMaster
    API2 --> PGSlave1
    API3 --> PGSlave2

    API1 --> RedisCluster
    API2 --> RedisCluster
    API3 --> RedisCluster

    API1 --> NATSCluster
    API2 --> NATSCluster
    API3 --> NATSCluster

    Worker1 --> NATSCluster
    Worker2 --> NATSCluster

    Worker1 --> S3
    Worker2 --> SES
    Worker1 --> SNS
```

## 🔧 Configuration & Environment Management

### Environment Structure
```
config/
├── local.yaml          # Local development
├── development.yaml    # Development environment
├── staging.yaml        # Staging environment
└── production.yaml     # Production environment
```

### Configuration Layers
1. **Default Values** - Hard-coded defaults in structs
2. **Configuration Files** - YAML files per environment
3. **Environment Variables** - Runtime overrides
4. **Command Line Flags** - Deployment-time overrides

## 📊 Monitoring & Observability

### Metrics Collection
- **Application Metrics**: Request latency, throughput, error rates
- **Business Metrics**: User registrations, listing creations, transaction volumes
- **Infrastructure Metrics**: CPU, memory, disk usage, database connections

### Distributed Tracing
- **Request Tracing**: End-to-end request flow across bounded contexts
- **Event Tracing**: Asynchronous event processing flows
- **Database Query Tracing**: SQL query performance analysis

### Logging Strategy
- **Structured Logging**: JSON format with consistent fields
- **Context Propagation**: Correlation IDs across service boundaries
- **Log Levels**: DEBUG, INFO, WARN, ERROR with appropriate sampling

This comprehensive architecture demonstrates a production-ready modular monolith that can scale horizontally while maintaining clean boundaries between bounded contexts, enabling future microservices extraction if needed.