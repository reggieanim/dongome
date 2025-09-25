package app

import (
	"context"
	"time"

	"dongome/internal/users/domain"
	"dongome/pkg/errors"
	"dongome/pkg/events"
)

// RegisterUserCommand represents the command to register a user
type RegisterUserCommand struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
}

// LoginCommand represents the command to login a user
type LoginCommand struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// UpgradeToSellerCommand represents the command to upgrade user to seller
type UpgradeToSellerCommand struct {
	UserID          string `json:"user_id" binding:"required"`
	BusinessName    string `json:"business_name" binding:"required"`
	BusinessAddress string `json:"business_address" binding:"required"`
}

// UserService handles user-related use cases
type UserService struct {
	userRepo domain.UserRepository
	eventBus events.EventBus
}

// NewUserService creates a new user service
func NewUserService(userRepo domain.UserRepository, eventBus events.EventBus) *UserService {
	return &UserService{
		userRepo: userRepo,
		eventBus: eventBus,
	}
}

// RegisterUser registers a new user
func (s *UserService) RegisterUser(ctx context.Context, cmd RegisterUserCommand) (*domain.User, error) {
	// Check if user already exists
	existing, _ := s.userRepo.FindByEmail(cmd.Email)
	if existing != nil {
		return nil, errors.ConflictError("user with this email already exists")
	}

	// Create new user
	user, err := domain.NewUser(cmd.Email, cmd.Password, cmd.FirstName, cmd.LastName)
	if err != nil {
		return nil, err
	}

	// Save user
	if err := s.userRepo.Save(user); err != nil {
		return nil, err
	}

	// Publish UserRegistered event
	event, err := events.NewEvent(
		domain.UserRegisteredEvent,
		user.ID,
		domain.UserRegistered{
			UserID:            user.ID,
			Email:             user.Email,
			FirstName:         user.FirstName,
			LastName:          user.LastName,
			Role:              user.Role,
			VerificationToken: user.VerificationToken,
			Timestamp:         time.Now(),
		},
	)
	if err != nil {
		return nil, err
	}

	if err := s.eventBus.Publish(ctx, event); err != nil {
		// Log error but don't fail the operation
		// Consider implementing eventual consistency or retry mechanism
	}

	return user, nil
}

// LoginUser authenticates a user
func (s *UserService) LoginUser(ctx context.Context, cmd LoginCommand) (*domain.User, error) {
	// Find user by email
	user, err := s.userRepo.FindByEmail(cmd.Email)
	if err != nil {
		return nil, errors.UnauthorizedError("invalid credentials")
	}

	// Validate password
	if err := user.ValidatePassword(cmd.Password); err != nil {
		return nil, errors.UnauthorizedError("invalid credentials")
	}

	// Check if user is active
	if !user.IsActive() {
		return nil, errors.UnauthorizedError("account is not active")
	}

	// Update last login
	user.UpdateLastLogin()
	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	// Publish UserLoggedIn event
	event, err := events.NewEvent(
		domain.UserLoggedInEvent,
		user.ID,
		domain.UserLoggedIn{
			UserID:    user.ID,
			Email:     user.Email,
			Timestamp: time.Now(),
		},
	)
	if err != nil {
		return nil, err
	}

	if err := s.eventBus.Publish(ctx, event); err != nil {
		// Log error but don't fail the operation
	}

	return user, nil
}

// VerifyEmail verifies a user's email
func (s *UserService) VerifyEmail(ctx context.Context, token string) error {
	// Find user by verification token
	user, err := s.userRepo.FindByVerificationToken(token)
	if err != nil {
		return errors.NotFoundError("invalid verification token")
	}

	// Verify email
	user.VerifyEmail()
	if err := s.userRepo.Update(user); err != nil {
		return err
	}

	// Publish UserEmailVerified event
	event, err := events.NewEvent(
		domain.UserEmailVerifiedEvent,
		user.ID,
		domain.UserEmailVerified{
			UserID:    user.ID,
			Email:     user.Email,
			Timestamp: time.Now(),
		},
	)
	if err != nil {
		return err
	}

	return s.eventBus.Publish(ctx, event)
}

// UpgradeToSeller upgrades a user to seller
func (s *UserService) UpgradeToSeller(ctx context.Context, cmd UpgradeToSellerCommand) error {
	// Find user
	user, err := s.userRepo.FindByID(cmd.UserID)
	if err != nil {
		return errors.NotFoundError("user not found")
	}

	// Upgrade to seller
	if err := user.UpgradeToSeller(cmd.BusinessName, cmd.BusinessAddress); err != nil {
		return err
	}

	// Update user
	if err := s.userRepo.Update(user); err != nil {
		return err
	}

	// Publish UserUpgradedToSeller event
	event, err := events.NewEvent(
		domain.UserUpgradedToSellerEvent,
		user.ID,
		domain.UserUpgradedToSeller{
			UserID:       user.ID,
			Email:        user.Email,
			BusinessName: cmd.BusinessName,
			Timestamp:    time.Now(),
		},
	)
	if err != nil {
		return err
	}

	return s.eventBus.Publish(ctx, event)
}

// GetUser retrieves a user by ID
func (s *UserService) GetUser(ctx context.Context, userID string) (*domain.User, error) {
	return s.userRepo.FindByID(userID)
}

// GetUserByEmail retrieves a user by email
func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	return s.userRepo.FindByEmail(email)
}
