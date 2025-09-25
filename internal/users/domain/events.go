package domain

import (
	"time"
)

// Event types
const (
	UserRegisteredEvent       = "user.registered"
	UserEmailVerifiedEvent    = "user.email_verified"
	UserUpgradedToSellerEvent = "user.upgraded_to_seller"
	SellerVerifiedEvent       = "seller.verified"
	UserSuspendedEvent        = "user.suspended"
	UserActivatedEvent        = "user.activated"
	UserLoggedInEvent         = "user.logged_in"
)

// UserRegistered represents the event when a user registers
type UserRegistered struct {
	UserID            string    `json:"user_id"`
	Email             string    `json:"email"`
	FirstName         string    `json:"first_name"`
	LastName          string    `json:"last_name"`
	Role              UserRole  `json:"role"`
	VerificationToken string    `json:"verification_token"`
	Timestamp         time.Time `json:"timestamp"`
}

// UserEmailVerified represents the event when a user verifies their email
type UserEmailVerified struct {
	UserID    string    `json:"user_id"`
	Email     string    `json:"email"`
	Timestamp time.Time `json:"timestamp"`
}

// UserUpgradedToSeller represents the event when a user becomes a seller
type UserUpgradedToSeller struct {
	UserID       string    `json:"user_id"`
	Email        string    `json:"email"`
	BusinessName string    `json:"business_name"`
	Timestamp    time.Time `json:"timestamp"`
}

// SellerVerified represents the event when a seller gets verified
type SellerVerified struct {
	UserID       string    `json:"user_id"`
	SellerID     string    `json:"seller_id"`
	Email        string    `json:"email"`
	BusinessName string    `json:"business_name"`
	Timestamp    time.Time `json:"timestamp"`
}

// UserSuspended represents the event when a user is suspended
type UserSuspended struct {
	UserID    string    `json:"user_id"`
	Email     string    `json:"email"`
	Reason    string    `json:"reason"`
	Timestamp time.Time `json:"timestamp"`
}

// UserActivated represents the event when a user is activated
type UserActivated struct {
	UserID    string    `json:"user_id"`
	Email     string    `json:"email"`
	Timestamp time.Time `json:"timestamp"`
}

// UserLoggedIn represents the event when a user logs in
type UserLoggedIn struct {
	UserID    string    `json:"user_id"`
	Email     string    `json:"email"`
	Timestamp time.Time `json:"timestamp"`
}
