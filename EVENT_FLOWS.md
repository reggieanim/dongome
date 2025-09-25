# Event Flow Documentation - Dongome Marketplace

## 🔄 Complete Event Flow Scenarios

### 1. User Registration Flow

```mermaid
sequenceDiagram
    participant Client as Mobile/Web Client
    participant API as API Gateway
    participant Users as Users Context
    participant Events as Event Bus (NATS)
    participant Worker as Background Worker
    participant Notifications as Notifications Context
    participant Email as Email Service
    participant DB as PostgreSQL

    Note over Client,DB: User Registration Complete Flow

    Client->>API: POST /api/v1/users/register
    Note right of Client: {email, password, first_name, last_name}
    
    API->>Users: RegisterUser(command)
    Users->>Users: Validate business rules
    Users->>Users: Create User aggregate
    Users->>Users: Hash password
    Users->>Users: Generate verification token
    Users->>DB: Save user to database
    DB-->>Users: User saved
    
    Users->>Events: Publish UserRegistered event
    Note right of Events: {user_id, email, first_name, verification_token}
    
    Users->>API: Return user response + verification token
    API->>Client: 201 Created + User data
    
    Note over Events,Email: Asynchronous Event Processing
    
    Events->>Worker: UserRegistered event received
    Worker->>Notifications: Send welcome email
    Notifications->>Notifications: Create email notification
    Notifications->>Email: Send verification email
    Email-->>Notifications: Email sent confirmation
    Notifications->>Events: Publish EmailSent event
    Notifications->>DB: Save notification record
    
    Events->>Worker: EmailSent event received
    Worker->>Worker: Log successful email delivery
```

### 2. Email Verification Flow

```mermaid
sequenceDiagram
    participant Client as Mobile/Web Client
    participant API as API Gateway
    participant Users as Users Context
    participant Events as Event Bus (NATS)
    participant Worker as Background Worker
    participant Notifications as Notifications Context
    participant Analytics as Analytics Context
    participant DB as PostgreSQL

    Client->>API: POST /api/v1/users/verify-email
    Note right of Client: {verification_token}
    
    API->>Users: VerifyEmail(token)
    Users->>DB: Find user by verification token
    DB-->>Users: User found
    
    Users->>Users: Validate token & business rules
    Users->>Users: Mark email as verified
    Users->>Users: Update user status to ACTIVE
    Users->>DB: Update user record
    
    Users->>Events: Publish UserEmailVerified event
    Note right of Events: {user_id, email, verified_at}
    
    Users->>API: Return success response
    API->>Client: 200 OK - Email verified
    
    Note over Events,Analytics: Cross-Context Event Processing
    
    Events->>Worker: UserEmailVerified event
    
    Worker->>Notifications: Send verification success email
    Notifications->>Notifications: Create success notification
    Notifications->>DB: Save notification
    
    Worker->>Analytics: Record user verification metric
    Analytics->>Analytics: Update user engagement stats
    Analytics->>DB: Save analytics data
    
    Worker->>Users: Update user engagement score
    Users->>DB: Update user profile
```

### 3. Seller Upgrade Flow

```mermaid
sequenceDiagram
    participant Client as Mobile/Web Client
    participant API as API Gateway
    participant Users as Users Context
    participant Events as Event Bus (NATS)
    participant Worker as Background Worker
    participant Notifications as Notifications Context
    participant Listings as Listings Context
    participant Reviews as Reviews Context
    participant DB as PostgreSQL

    Client->>API: POST /api/v1/users/upgrade-to-seller
    Note right of Client: {business_name, business_address, business_phone}
    
    API->>Users: UpgradeToSeller(command)
    Users->>DB: Load user aggregate
    DB-->>Users: User loaded
    
    Users->>Users: Validate business rules
    Note right of Users: - Email must be verified<br/>- Not already a seller
    
    Users->>Users: Create SellerProfile entity
    Users->>Users: Update user role to SELLER
    Users->>DB: Save user + seller profile
    
    Users->>Events: Publish UserUpgradedToSeller event
    Note right of Events: {user_id, business_name, verification_status}
    
    Users->>API: Return seller profile response
    API->>Client: 200 OK + Seller profile
    
    Note over Events,Reviews: Multi-Context Event Processing
    
    Events->>Worker: UserUpgradedToSeller event received
    
    Worker->>Notifications: Send seller welcome email
    Notifications->>Notifications: Create seller onboarding notification
    Notifications->>DB: Save notification
    
    Worker->>Listings: Initialize seller listing capabilities
    Listings->>Listings: Create seller listing preferences
    Listings->>DB: Save listing settings
    
    Worker->>Reviews: Initialize seller review profile
    Reviews->>Reviews: Create seller review aggregate
    Reviews->>DB: Save review profile
    
    Note over Worker: Admin Notification for Manual Verification
    Worker->>Notifications: Notify admin of pending seller verification
    Notifications->>Notifications: Create admin notification
```

### 4. Listing Creation Flow

```mermaid
sequenceDiagram
    participant Client as Mobile/Web Client
    participant API as API Gateway
    participant Listings as Listings Context
    participant Events as Event Bus (NATS)
    participant Worker as Background Worker
    participant Users as Users Context
    participant Notifications as Notifications Context
    participant Search as Search Service
    participant Storage as File Storage
    participant DB as PostgreSQL

    Client->>API: POST /api/v1/listings
    Note right of Client: {title, description, price, category_id, images[]}
    
    API->>Listings: CreateListing(command)
    Listings->>DB: Validate seller exists & verified
    DB-->>Listings: Seller validated
    
    Listings->>Listings: Validate business rules
    Note right of Listings: - Price > 0<br/>- Category is active<br/>- Image limits
    
    Listings->>Storage: Upload listing images
    Storage-->>Listings: Image URLs returned
    
    Listings->>Listings: Create Listing aggregate
    Listings->>DB: Save listing
    
    Listings->>Events: Publish ListingCreated event
    Note right of Events: {listing_id, seller_id, title, category_id, price}
    
    Listings->>API: Return listing response
    API->>Client: 201 Created + Listing data
    
    Note over Events,Search: Asynchronous Processing Pipeline
    
    Events->>Worker: ListingCreated event received
    
    Worker->>Search: Index listing for search
    Search->>Search: Add to search index
    Note right of Search: ElasticSearch/Algolia
    
    Worker->>Users: Update seller listing count
    Users->>DB: Increment seller stats
    
    Worker->>Notifications: Notify followers of new listing
    Notifications->>DB: Query seller followers
    Notifications->>Notifications: Create bulk notifications
    Notifications->>DB: Save notifications
    
    Events->>Events: Publish ListingIndexed event
    Events->>Worker: ListingIndexed event
    Worker->>Worker: Log successful listing creation
```

### 5. Transaction/Purchase Flow

```mermaid
sequenceDiagram
    participant Client as Buyer Client
    participant API as API Gateway
    participant Transactions as Transactions Context
    participant Events as Event Bus (NATS)
    participant Worker as Background Worker
    participant Listings as Listings Context
    participant Users as Users Context
    participant Notifications as Notifications Context
    participant Payment as MoMo Payment Gateway
    participant Escrow as Escrow Service
    participant DB as PostgreSQL

    Client->>API: POST /api/v1/transactions
    Note right of Client: {listing_id, delivery_address, payment_method}
    
    API->>Transactions: CreateTransaction(command)
    Transactions->>DB: Load listing & validate availability
    DB-->>Transactions: Listing available
    
    Transactions->>Transactions: Validate business rules
    Note right of Transactions: - Listing is active<br/>- Buyer != Seller<br/>- Sufficient funds
    
    Transactions->>Transactions: Create Transaction aggregate
    Transactions->>DB: Save transaction
    
    Transactions->>Events: Publish TransactionCreated event
    Note right of Events: {transaction_id, listing_id, buyer_id, seller_id, amount}
    
    Transactions->>Payment: Initiate payment
    Payment-->>Transactions: Payment reference returned
    
    Transactions->>API: Return transaction + payment details
    API->>Client: 201 Created + Payment instructions
    
    Note over Events,Escrow: Cross-Service Coordination
    
    Events->>Worker: TransactionCreated event received
    
    Worker->>Listings: Reserve listing
    Listings->>Listings: Update listing status to RESERVED
    Listings->>DB: Save listing update
    Listings->>Events: Publish ListingReserved event
    
    Worker->>Notifications: Notify seller of new order
    Notifications->>Notifications: Create seller notification
    Notifications->>DB: Save notification
    
    Worker->>Escrow: Setup escrow account
    Escrow->>Escrow: Create escrow for transaction
    
    Note over Payment,DB: Payment Completion Flow
    
    Payment->>API: Payment webhook notification
    API->>Transactions: ProcessPayment(webhook_data)
    
    Transactions->>Transactions: Validate payment
    Transactions->>Transactions: Update transaction status
    Transactions->>DB: Save transaction update
    
    Transactions->>Events: Publish PaymentCompleted event
    Events->>Worker: PaymentCompleted event
    
    Worker->>Escrow: Fund escrow account
    Worker->>Notifications: Notify buyer & seller
    Worker->>Listings: Update listing status to SOLD
```

### 6. Review Creation Flow

```mermaid
sequenceDiagram
    participant Client as Buyer Client
    participant API as API Gateway
    participant Reviews as Reviews Context
    participant Events as Event Bus (NATS)
    participant Worker as Background Worker
    participant Users as Users Context
    participant Transactions as Transactions Context
    participant Notifications as Notifications Context
    participant Analytics as Analytics Context
    participant DB as PostgreSQL

    Client->>API: POST /api/v1/reviews
    Note right of Client: {transaction_id, rating, comment}
    
    API->>Reviews: CreateReview(command)
    Reviews->>DB: Validate transaction exists & completed
    DB-->>Reviews: Transaction validated
    
    Reviews->>Reviews: Validate business rules
    Note right of Reviews: - One review per transaction<br/>- Rating 1-5<br/>- Within 30 days
    
    Reviews->>Reviews: Create Review aggregate
    Reviews->>DB: Save review
    
    Reviews->>Events: Publish ReviewCreated event
    Note right of Events: {review_id, transaction_id, reviewer_id, reviewee_id, rating}
    
    Reviews->>API: Return review response
    API->>Client: 201 Created + Review data
    
    Note over Events,Analytics: Rating Calculation & Updates
    
    Events->>Worker: ReviewCreated event received
    
    Worker->>Users: Update seller rating
    Users->>DB: Query all seller reviews
    Users->>Users: Calculate new average rating
    Users->>DB: Update seller profile
    Users->>Events: Publish SellerRatingUpdated event
    
    Worker->>Notifications: Notify seller of new review
    Notifications->>Notifications: Create review notification
    Notifications->>DB: Save notification
    
    Worker->>Analytics: Update review analytics
    Analytics->>Analytics: Update marketplace metrics
    Analytics->>DB: Save analytics data
    
    Worker->>Transactions: Mark transaction as reviewed
    Transactions->>DB: Update transaction status
```

## 🔄 Event Processing Patterns

### Event Processing Guarantees

```mermaid
graph TB
    subgraph "Event Processing Pipeline"
        Event[Domain Event Published]
        Queue[NATS JetStream Queue]
        Worker[Background Worker]
        Handler[Event Handler]
        
        Event --> Queue
        Queue --> Worker
        Worker --> Handler
        
        subgraph "Failure Handling"
            Retry[Automatic Retry]
            DLQ[Dead Letter Queue]
            Alert[Admin Alert]
        end
        
        Handler -->|Failure| Retry
        Retry -->|Max Retries| DLQ
        DLQ --> Alert
    end
    
    subgraph "Processing Guarantees"
        AtLeastOnce["At-Least-Once Delivery"]
        Idempotent["Idempotent Handlers"]
        Ordering["Event Ordering"]
    end
```

### Event Handler Patterns

1. **Immediate Processing**: Critical events processed synchronously
2. **Background Processing**: Non-critical events processed asynchronously
3. **Batch Processing**: Multiple events processed together for efficiency
4. **Saga Pattern**: Multi-step business processes coordinated via events

### Event Versioning Strategy

```go
// Event versioning for backward compatibility
type UserRegisteredEvent struct {
    Version   int    `json:"version"`   // Event schema version
    EventID   string `json:"event_id"`
    EventType string `json:"event_type"`
    // ... event data
}

// Handler supports multiple versions
func (h *UserEventHandler) Handle(ctx context.Context, event map[string]interface{}) error {
    version := int(event["version"].(float64))
    
    switch version {
    case 1:
        return h.handleV1(ctx, event)
    case 2:
        return h.handleV2(ctx, event)
    default:
        return fmt.Errorf("unsupported event version: %d", version)
    }
}
```

## 📊 Event Monitoring & Observability

### Event Metrics Collected

1. **Event Volume**: Events published/consumed per second
2. **Processing Latency**: Time from publish to successful processing
3. **Error Rates**: Failed event processing percentage
4. **Queue Depth**: Pending events in each queue
5. **Handler Performance**: Processing time per event handler

### Event Tracing

Each event carries correlation IDs for distributed tracing:

```json
{
  "event_id": "evt_123",
  "correlation_id": "req_456",
  "causation_id": "evt_789",
  "event_type": "UserRegistered",
  "aggregate_id": "user_abc",
  "timestamp": "2025-09-25T10:30:00Z",
  "data": { ... }
}
```

This comprehensive event flow documentation shows how the Dongome marketplace achieves loose coupling between bounded contexts while maintaining strong consistency within each context through carefully orchestrated event-driven communication.