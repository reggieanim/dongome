# Technical Implementation Patterns - Dongome Modular Monolith

## 🏗️ Hexagonal Architecture Implementation

### Directory Structure Pattern
```
internal/
├── {context}/
│   ├── domain/           # Domain Layer (Pure Business Logic)
│   │   ├── {aggregate}.go        # Aggregate Root
│   │   ├── {entity}.go          # Domain Entities
│   │   ├── {value_object}.go    # Value Objects
│   │   ├── events.go            # Domain Events
│   │   ├── repository.go        # Repository Interface (Port)
│   │   └── service.go           # Domain Services
│   ├── app/              # Application Layer (Use Cases)
│   │   ├── service.go           # Application Services
│   │   ├── commands.go          # Command DTOs
│   │   ├── queries.go           # Query DTOs
│   │   └── handlers.go          # Command/Query Handlers
│   └── infra/            # Infrastructure Layer (Adapters)
│       ├── repository.go        # Repository Implementation
│       ├── handlers.go          # HTTP Handlers
│       ├── clients.go           # External API Clients
│       └── mappers.go           # Domain ↔ DTO Mapping
```

## 🎯 Domain Layer Implementation Patterns

### 1. Aggregate Root Pattern

```go
// internal/users/domain/user.go
package domain

import (
    "time"
    "dongome/pkg/errors"
    "github.com/google/uuid"
)

// User represents the aggregate root for user-related operations
type User struct {
    // Identity
    ID                string     `gorm:"type:uuid;primary_key"`
    Email             string     `gorm:"uniqueIndex;not null"`
    
    // Value Objects
    Password          Password   // Encapsulated password logic
    Profile           Profile    // User profile information
    Location          Location   // Geographic location
    
    // State
    Status            UserStatus
    Role              UserRole
    EmailVerified     bool
    PhoneVerified     bool
    
    // Audit
    CreatedAt         time.Time
    UpdatedAt         time.Time
    
    // Associations
    SellerProfile     *SellerProfile `gorm:"foreignKey:UserID"`
    
    // Domain Events (not persisted)
    events            []DomainEvent `gorm:"-"`
}

// Business Methods
func (u *User) VerifyEmail(token string) error {
    if u.VerificationToken != token {
        return errors.ValidationError("invalid verification token")
    }
    
    if u.EmailVerified {
        return errors.BusinessRuleError("email already verified")
    }
    
    u.EmailVerified = true
    u.Status = UserStatusActive
    u.AddEvent(NewUserEmailVerifiedEvent(u.ID, u.Email))
    
    return nil
}

func (u *User) UpgradeToSeller(businessName, businessAddress string) error {
    if !u.EmailVerified {
        return errors.BusinessRuleError("email must be verified to become seller")
    }
    
    if u.Role == UserRoleSeller {
        return errors.BusinessRuleError("user is already a seller")
    }
    
    profile := &SellerProfile{
        ID:              uuid.New().String(),
        UserID:          u.ID,
        BusinessName:    businessName,
        BusinessAddress: businessAddress,
        VerificationStatus: VerificationStatusPending,
        CreatedAt:       time.Now(),
    }
    
    u.SellerProfile = profile
    u.Role = UserRoleSeller
    u.AddEvent(NewUserUpgradedToSellerEvent(u.ID, businessName))
    
    return nil
}

// Event Management
func (u *User) AddEvent(event DomainEvent) {
    u.events = append(u.events, event)
}

func (u *User) GetEvents() []DomainEvent {
    return u.events
}

func (u *User) ClearEvents() {
    u.events = nil
}
```

### 2. Value Objects Pattern

```go
// internal/users/domain/password.go
package domain

import (
    "golang.org/x/crypto/bcrypt"
    "dongome/pkg/errors"
)

// Password represents a value object for password handling
type Password struct {
    hash string
}

// NewPassword creates a new password value object
func NewPassword(plaintext string) (Password, error) {
    if len(plaintext) < 8 {
        return Password{}, errors.ValidationError("password must be at least 8 characters")
    }
    
    hash, err := bcrypt.GenerateFromPassword([]byte(plaintext), bcrypt.DefaultCost)
    if err != nil {
        return Password{}, err
    }
    
    return Password{hash: string(hash)}, nil
}

// Verify checks if the provided plaintext matches the stored hash
func (p Password) Verify(plaintext string) bool {
    return bcrypt.CompareHashAndPassword([]byte(p.hash), []byte(plaintext)) == nil
}

// Hash returns the password hash for storage
func (p Password) Hash() string {
    return p.hash
}
```

### 3. Domain Events Pattern

```go
// internal/users/domain/events.go
package domain

import (
    "time"
    "github.com/google/uuid"
)

// DomainEvent interface for all domain events
type DomainEvent interface {
    EventID() string
    EventType() string
    AggregateID() string
    OccurredAt() time.Time
    EventData() interface{}
}

// Base event struct
type baseEvent struct {
    eventID     string
    eventType   string
    aggregateID string
    occurredAt  time.Time
}

func (e baseEvent) EventID() string     { return e.eventID }
func (e baseEvent) EventType() string   { return e.eventType }
func (e baseEvent) AggregateID() string { return e.aggregateID }
func (e baseEvent) OccurredAt() time.Time { return e.occurredAt }

// UserRegisteredEvent represents a user registration event
type UserRegisteredEvent struct {
    baseEvent
    Email     string `json:"email"`
    FirstName string `json:"first_name"`
    LastName  string `json:"last_name"`
}

func NewUserRegisteredEvent(userID, email, firstName, lastName string) UserRegisteredEvent {
    return UserRegisteredEvent{
        baseEvent: baseEvent{
            eventID:     uuid.New().String(),
            eventType:   "UserRegistered",
            aggregateID: userID,
            occurredAt:  time.Now(),
        },
        Email:     email,
        FirstName: firstName,
        LastName:  lastName,
    }
}

func (e UserRegisteredEvent) EventData() interface{} {
    return struct {
        Email     string `json:"email"`
        FirstName string `json:"first_name"`
        LastName  string `json:"last_name"`
    }{
        Email:     e.Email,
        FirstName: e.FirstName,
        LastName:  e.LastName,
    }
}
```

## 🔄 Application Layer Implementation Patterns

### 1. Application Services Pattern

```go
// internal/users/app/service.go
package app

import (
    "context"
    "dongome/internal/users/domain"
    "dongome/pkg/events"
    "dongome/pkg/errors"
)

// UserService handles user-related use cases
type UserService struct {
    repo        domain.UserRepository
    eventBus    events.EventBus
}

func NewUserService(repo domain.UserRepository, eventBus events.EventBus) *UserService {
    return &UserService{
        repo:     repo,
        eventBus: eventBus,
    }
}

// RegisterUser handles user registration use case
func (s *UserService) RegisterUser(ctx context.Context, cmd RegisterUserCommand) (*UserResponse, error) {
    // Check if user already exists
    existingUser, _ := s.repo.FindByEmail(ctx, cmd.Email)
    if existingUser != nil {
        return nil, errors.ConflictError("user with this email already exists")
    }
    
    // Create new user aggregate
    user, err := domain.NewUser(cmd.Email, cmd.Password, cmd.FirstName, cmd.LastName)
    if err != nil {
        return nil, err
    }
    
    // Save to repository
    if err := s.repo.Save(ctx, user); err != nil {
        return nil, err
    }
    
    // Publish domain events
    for _, event := range user.GetEvents() {
        if err := s.eventBus.Publish(ctx, event); err != nil {
            // Log error but don't fail the operation
            // Events will be retried by background workers
        }
    }
    
    user.ClearEvents()
    
    return &UserResponse{
        ID:        user.ID,
        Email:     user.Email,
        FirstName: user.FirstName,
        LastName:  user.LastName,
        Status:    string(user.Status),
        Role:      string(user.Role),
        CreatedAt: user.CreatedAt,
    }, nil
}

// LoginUser handles user authentication
func (s *UserService) LoginUser(ctx context.Context, cmd LoginUserCommand) (*LoginResponse, error) {
    user, err := s.repo.FindByEmail(ctx, cmd.Email)
    if err != nil {
        return nil, errors.NotFoundError("user not found")
    }
    
    if !user.Password.Verify(cmd.Password) {
        return nil, errors.UnauthorizedError("invalid credentials")
    }
    
    if user.Status == domain.UserStatusSuspended {
        return nil, errors.ForbiddenError("account suspended")
    }
    
    // Update last login
    user.UpdateLastLogin()
    if err := s.repo.Save(ctx, user); err != nil {
        return nil, err
    }
    
    // Publish login event
    event := domain.NewUserLoggedInEvent(user.ID, user.Email)
    s.eventBus.Publish(ctx, event)
    
    // Generate JWT token (simplified)
    token := s.generateJWT(user)
    
    return &LoginResponse{
        Token: token,
        User: UserResponse{
            ID:        user.ID,
            Email:     user.Email,
            FirstName: user.FirstName,
            LastName:  user.LastName,
            Status:    string(user.Status),
            Role:      string(user.Role),
        },
    }, nil
}
```

### 2. Command/Query Objects Pattern

```go
// internal/users/app/commands.go
package app

import "time"

// Commands for state-changing operations
type RegisterUserCommand struct {
    Email     string `json:"email" validate:"required,email"`
    Password  string `json:"password" validate:"required,min=8"`
    FirstName string `json:"first_name" validate:"required"`
    LastName  string `json:"last_name" validate:"required"`
}

type LoginUserCommand struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required"`
}

type VerifyEmailCommand struct {
    Token string `json:"token" validate:"required"`
}

type UpgradeToSellerCommand struct {
    UserID          string `json:"user_id" validate:"required"`
    BusinessName    string `json:"business_name" validate:"required"`
    BusinessAddress string `json:"business_address" validate:"required"`
    BusinessPhone   string `json:"business_phone"`
    BusinessEmail   string `json:"business_email" validate:"email"`
}

// Queries for data retrieval
type GetUserQuery struct {
    UserID string `validate:"required"`
}

type SearchUsersQuery struct {
    Email    string
    Role     string
    Status   string
    Limit    int `validate:"max=100"`
    Offset   int
}

// Response objects
type UserResponse struct {
    ID               string     `json:"id"`
    Email            string     `json:"email"`
    FirstName        string     `json:"first_name"`
    LastName         string     `json:"last_name"`
    Status           string     `json:"status"`
    Role             string     `json:"role"`
    EmailVerified    bool       `json:"email_verified"`
    PhoneVerified    bool       `json:"phone_verified"`
    LastLoginAt      *time.Time `json:"last_login_at"`
    CreatedAt        time.Time  `json:"created_at"`
    SellerProfile    *SellerProfileResponse `json:"seller_profile,omitempty"`
}

type LoginResponse struct {
    Token string       `json:"token"`
    User  UserResponse `json:"user"`
}
```

## 🔌 Infrastructure Layer Implementation Patterns

### 1. Repository Implementation Pattern

```go
// internal/users/infra/repository.go
package infra

import (
    "context"
    "dongome/internal/users/domain"
    "gorm.io/gorm"
)

// userRepository implements domain.UserRepository
type userRepository struct {
    db *gorm.DB
}

func NewUserRepository(db *gorm.DB) domain.UserRepository {
    return &userRepository{db: db}
}

func (r *userRepository) Save(ctx context.Context, user *domain.User) error {
    return r.db.WithContext(ctx).Save(user).Error
}

func (r *userRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
    var user domain.User
    err := r.db.WithContext(ctx).
        Preload("SellerProfile").
        First(&user, "id = ?", id).Error
    
    if err != nil {
        if err == gorm.ErrRecordNotFound {
            return nil, domain.ErrUserNotFound
        }
        return nil, err
    }
    
    return &user, nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
    var user domain.User
    err := r.db.WithContext(ctx).
        Preload("SellerProfile").
        First(&user, "email = ?", email).Error
    
    if err != nil {
        if err == gorm.ErrRecordNotFound {
            return nil, domain.ErrUserNotFound
        }
        return nil, err
    }
    
    return &user, nil
}

func (r *userRepository) Search(ctx context.Context, criteria domain.SearchCriteria) ([]*domain.User, error) {
    query := r.db.WithContext(ctx)
    
    if criteria.Email != "" {
        query = query.Where("email ILIKE ?", "%"+criteria.Email+"%")
    }
    
    if criteria.Role != "" {
        query = query.Where("role = ?", criteria.Role)
    }
    
    if criteria.Status != "" {
        query = query.Where("status = ?", criteria.Status)
    }
    
    var users []*domain.User
    err := query.
        Preload("SellerProfile").
        Limit(criteria.Limit).
        Offset(criteria.Offset).
        Find(&users).Error
    
    return users, err
}
```

### 2. HTTP Handlers Pattern

```go
// internal/users/infra/handlers.go
package infra

import (
    "net/http"
    "dongome/internal/users/app"
    "github.com/gin-gonic/gin"
)

// UserHandlers contains HTTP handlers for user operations
type UserHandlers struct {
    userService *app.UserService
}

func NewUserHandlers(userService *app.UserService) *UserHandlers {
    return &UserHandlers{userService: userService}
}

// RegisterRoutes registers all user routes
func (h *UserHandlers) RegisterRoutes(router *gin.RouterGroup) {
    users := router.Group("/users")
    {
        users.POST("/register", h.Register)
        users.POST("/login", h.Login)
        users.POST("/verify-email", h.VerifyEmail)
        users.POST("/upgrade-to-seller", h.UpgradeToSeller)
        users.GET("/:id", h.GetUser)
        users.GET("", h.SearchUsers)
    }
}

// Register handles user registration
func (h *UserHandlers) Register(c *gin.Context) {
    var cmd app.RegisterUserCommand
    
    if err := c.ShouldBindJSON(&cmd); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    response, err := h.userService.RegisterUser(c.Request.Context(), cmd)
    if err != nil {
        h.handleError(c, err)
        return
    }
    
    c.JSON(http.StatusCreated, response)
}

// Login handles user authentication
func (h *UserHandlers) Login(c *gin.Context) {
    var cmd app.LoginUserCommand
    
    if err := c.ShouldBindJSON(&cmd); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    response, err := h.userService.LoginUser(c.Request.Context(), cmd)
    if err != nil {
        h.handleError(c, err)
        return
    }
    
    c.JSON(http.StatusOK, response)
}

// Error handling
func (h *UserHandlers) handleError(c *gin.Context, err error) {
    switch err.(type) {
    case *errors.ValidationError:
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
    case *errors.ConflictError:
        c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
    case *errors.NotFoundError:
        c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
    case *errors.UnauthorizedError:
        c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
    case *errors.ForbiddenError:
        c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
    default:
        c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
    }
}
```

## 🌐 Event-Driven Communication Patterns

### 1. Event Bus Implementation

```go
// pkg/events/eventbus.go
package events

import (
    "context"
    "encoding/json"
    "fmt"
    "dongome/pkg/logger"
    "github.com/nats-io/nats.go"
)

// EventBus handles event publishing and subscription
type EventBus struct {
    conn   *nats.Conn
    js     nats.JetStreamContext
    logger *logger.Logger
}

func NewEventBus(conn *nats.Conn, logger *logger.Logger) (*EventBus, error) {
    js, err := conn.JetStream()
    if err != nil {
        return nil, err
    }
    
    return &EventBus{
        conn:   conn,
        js:     js,
        logger: logger,
    }, nil
}

// Publish publishes a domain event
func (eb *EventBus) Publish(ctx context.Context, event DomainEvent) error {
    subject := fmt.Sprintf("events.%s.%s", event.AggregateType(), event.EventType())
    
    data, err := json.Marshal(event)
    if err != nil {
        return err
    }
    
    _, err = eb.js.PublishAsync(subject, data)
    if err != nil {
        eb.logger.Error("Failed to publish event",
            "event_type", event.EventType(),
            "aggregate_id", event.AggregateID(),
            "error", err,
        )
        return err
    }
    
    eb.logger.Info("Event published",
        "event_type", event.EventType(),
        "aggregate_id", event.AggregateID(),
        "subject", subject,
    )
    
    return nil
}

// Subscribe subscribes to events with a handler
func (eb *EventBus) Subscribe(subject string, handler EventHandler) error {
    _, err := eb.js.Subscribe(subject, func(msg *nats.Msg) {
        var event map[string]interface{}
        if err := json.Unmarshal(msg.Data, &event); err != nil {
            eb.logger.Error("Failed to unmarshal event", "error", err)
            return
        }
        
        ctx := context.Background()
        if err := handler.Handle(ctx, event); err != nil {
            eb.logger.Error("Event handler failed",
                "subject", subject,
                "error", err,
            )
            // Don't ack the message so it can be retried
            return
        }
        
        msg.Ack()
    })
    
    return err
}
```

### 2. Event Handlers Pattern

```go
// cmd/worker/handlers.go
package main

import (
    "context"
    "dongome/internal/notifications/app"
    "dongome/pkg/logger"
)

// UserRegisteredHandler handles UserRegistered events
type UserRegisteredHandler struct {
    notificationService *app.NotificationService
    logger             *logger.Logger
}

func NewUserRegisteredHandler(ns *app.NotificationService, logger *logger.Logger) *UserRegisteredHandler {
    return &UserRegisteredHandler{
        notificationService: ns,
        logger:             logger,
    }
}

func (h *UserRegisteredHandler) Handle(ctx context.Context, event map[string]interface{}) error {
    // Extract event data
    email, ok := event["email"].(string)
    if !ok {
        return fmt.Errorf("invalid email in event data")
    }
    
    firstName, ok := event["first_name"].(string)
    if !ok {
        return fmt.Errorf("invalid first_name in event data")
    }
    
    // Send welcome email
    cmd := app.SendEmailCommand{
        To:       email,
        Subject:  "Welcome to Dongome!",
        Template: "welcome_email",
        Data: map[string]interface{}{
            "first_name": firstName,
        },
    }
    
    if err := h.notificationService.SendEmail(ctx, cmd); err != nil {
        h.logger.Error("Failed to send welcome email",
            "email", email,
            "error", err,
        )
        return err
    }
    
    h.logger.Info("Welcome email sent",
        "email", email,
    )
    
    return nil
}
```

## 🔧 Configuration & Dependency Injection Patterns

### 1. Dependency Injection Container

```go
// cmd/api/main.go
package main

import (
    "dongome/internal/users/app"
    "dongome/internal/users/infra"
    "dongome/pkg/config"
    "dongome/pkg/db"
    "dongome/pkg/events"
    "dongome/pkg/logger"
)

// Container holds all application dependencies
type Container struct {
    Config  *config.Config
    Logger  *logger.Logger
    DB      *gorm.DB
    EventBus *events.EventBus
    
    // Repositories
    UserRepository domain.UserRepository
    
    // Services
    UserService *app.UserService
    
    // Handlers
    UserHandlers *infra.UserHandlers
}

// NewContainer creates a new dependency injection container
func NewContainer() (*Container, error) {
    // Load configuration
    cfg, err := config.Load()
    if err != nil {
        return nil, err
    }
    
    // Initialize logger
    logger := logger.New(cfg.Log.Level, cfg.Log.Format)
    
    // Initialize database
    database, err := db.Connect(cfg.Database)
    if err != nil {
        return nil, err
    }
    
    // Initialize event bus
    eventBus, err := events.NewEventBus(cfg.NATS, logger)
    if err != nil {
        return nil, err
    }
    
    // Initialize repositories
    userRepo := infra.NewUserRepository(database)
    
    // Initialize services
    userService := app.NewUserService(userRepo, eventBus)
    
    // Initialize handlers
    userHandlers := infra.NewUserHandlers(userService)
    
    return &Container{
        Config:         cfg,
        Logger:         logger,
        DB:            database,
        EventBus:      eventBus,
        UserRepository: userRepo,
        UserService:   userService,
        UserHandlers:  userHandlers,
    }, nil
}
```

This technical implementation guide shows the specific patterns used throughout the Dongome modular monolith, demonstrating how Domain-Driven Design and Hexagonal Architecture principles are applied in practice with Go.