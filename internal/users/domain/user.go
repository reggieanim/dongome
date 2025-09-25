package domain

import (
	"time"

	"dongome/pkg/errors"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// UserStatus represents the status of a user
type UserStatus string

const (
	UserStatusPending   UserStatus = "pending"
	UserStatusActive    UserStatus = "active"
	UserStatusSuspended UserStatus = "suspended"
	UserStatusDeactive  UserStatus = "deactive"
)

// UserRole represents the role of a user
type UserRole string

const (
	UserRoleBuyer  UserRole = "buyer"
	UserRoleSeller UserRole = "seller"
	UserRoleAdmin  UserRole = "admin"
)

// VerificationStatus represents seller verification status
type VerificationStatus string

const (
	VerificationStatusPending  VerificationStatus = "pending"
	VerificationStatusApproved VerificationStatus = "approved"
	VerificationStatusRejected VerificationStatus = "rejected"
)

// User represents a user aggregate root
type User struct {
	ID                string     `gorm:"type:uuid;primary_key" json:"id"`
	Email             string     `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash      string     `gorm:"not null" json:"-"`
	FirstName         string     `gorm:"not null" json:"first_name"`
	LastName          string     `gorm:"not null" json:"last_name"`
	PhoneNumber       string     `gorm:"uniqueIndex" json:"phone_number"`
	Avatar            string     `json:"avatar"`
	Status            UserStatus `gorm:"default:'pending'" json:"status"`
	Role              UserRole   `gorm:"default:'buyer'" json:"role"`
	EmailVerified     bool       `gorm:"default:false" json:"email_verified"`
	PhoneVerified     bool       `gorm:"default:false" json:"phone_verified"`
	VerificationToken string     `json:"-"`
	LastLoginAt       *time.Time `json:"last_login_at"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`

	// Seller-specific fields
	SellerProfile *SellerProfile `gorm:"foreignKey:UserID" json:"seller_profile,omitempty"`
}

// SellerProfile represents seller-specific information
type SellerProfile struct {
	ID                 string             `gorm:"type:uuid;primary_key" json:"id"`
	UserID             string             `gorm:"type:uuid;not null;uniqueIndex" json:"user_id"`
	BusinessName       string             `json:"business_name"`
	BusinessAddress    string             `json:"business_address"`
	BusinessPhone      string             `json:"business_phone"`
	BusinessEmail      string             `json:"business_email"`
	TaxNumber          string             `json:"tax_number"`
	VerificationStatus VerificationStatus `gorm:"default:'pending'" json:"verification_status"`
	VerificationNotes  string             `json:"verification_notes"`
	Rating             float64            `gorm:"default:0" json:"rating"`
	TotalReviews       int                `gorm:"default:0" json:"total_reviews"`
	CreatedAt          time.Time          `json:"created_at"`
	UpdatedAt          time.Time          `json:"updated_at"`
}

// NewUser creates a new user
func NewUser(email, password, firstName, lastName string) (*User, error) {
	if email == "" {
		return nil, errors.ValidationError("email is required")
	}
	if password == "" {
		return nil, errors.ValidationError("password is required")
	}
	if firstName == "" {
		return nil, errors.ValidationError("first name is required")
	}
	if lastName == "" {
		return nil, errors.ValidationError("last name is required")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Generate verification token
	verificationToken := uuid.New().String()

	return &User{
		ID:                uuid.New().String(),
		Email:             email,
		PasswordHash:      string(hashedPassword),
		FirstName:         firstName,
		LastName:          lastName,
		Status:            UserStatusPending,
		Role:              UserRoleBuyer,
		EmailVerified:     false,
		PhoneVerified:     false,
		VerificationToken: verificationToken,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}, nil
}

// ValidatePassword checks if the provided password matches the user's password
func (u *User) ValidatePassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
}

// VerifyEmail marks the user's email as verified
func (u *User) VerifyEmail() {
	u.EmailVerified = true
	u.Status = UserStatusActive
	u.VerificationToken = ""
	u.UpdatedAt = time.Now()
}

// UpgradeToSeller upgrades a buyer to seller
func (u *User) UpgradeToSeller(businessName, businessAddress string) error {
	if u.Role == UserRoleSeller {
		return errors.ValidationError("user is already a seller")
	}

	if !u.EmailVerified {
		return errors.ValidationError("email must be verified to become a seller")
	}

	u.Role = UserRoleSeller
	u.SellerProfile = &SellerProfile{
		ID:                 uuid.New().String(),
		UserID:             u.ID,
		BusinessName:       businessName,
		BusinessAddress:    businessAddress,
		VerificationStatus: VerificationStatusPending,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}
	u.UpdatedAt = time.Now()

	return nil
}

// Suspend suspends the user account
func (u *User) Suspend(reason string) {
	u.Status = UserStatusSuspended
	u.UpdatedAt = time.Now()
}

// Activate activates the user account
func (u *User) Activate() {
	u.Status = UserStatusActive
	u.UpdatedAt = time.Now()
}

// UpdateLastLogin updates the last login timestamp
func (u *User) UpdateLastLogin() {
	now := time.Now()
	u.LastLoginAt = &now
	u.UpdatedAt = time.Now()
}

// FullName returns the user's full name
func (u *User) FullName() string {
	return u.FirstName + " " + u.LastName
}

// IsActive checks if the user is active
func (u *User) IsActive() bool {
	return u.Status == UserStatusActive
}

// IsSeller checks if the user is a seller
func (u *User) IsSeller() bool {
	return u.Role == UserRoleSeller
}

// IsVerifiedSeller checks if the user is a verified seller
func (u *User) IsVerifiedSeller() bool {
	return u.IsSeller() && u.SellerProfile != nil &&
		u.SellerProfile.VerificationStatus == VerificationStatusApproved
}

// UserRepository defines the interface for user persistence
type UserRepository interface {
	Save(user *User) error
	FindByID(id string) (*User, error)
	FindByEmail(email string) (*User, error)
	FindByVerificationToken(token string) (*User, error)
	Update(user *User) error
	Delete(id string) error
}
