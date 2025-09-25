package domain
package domain_test

import (
	"testing"
	"time"

	"dongome/internal/users/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewUser(t *testing.T) {
	tests := []struct {
		name      string
		email     string
		password  string
		firstName string
		lastName  string
		wantErr   bool
	}{
		{
			name:      "valid user creation",
			email:     "test@example.com",
			password:  "password123",
			firstName: "John",
			lastName:  "Doe",
			wantErr:   false,
		},
		{
			name:      "empty email",
			email:     "",
			password:  "password123",
			firstName: "John",
			lastName:  "Doe",
			wantErr:   true,
		},
		{
			name:      "empty password",
			email:     "test@example.com",
			password:  "",
			firstName: "John",
			lastName:  "Doe",
			wantErr:   true,
		},
		{
			name:      "empty first name",
			email:     "test@example.com",
			password:  "password123",
			firstName: "",
			lastName:  "Doe",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := domain.NewUser(tt.email, tt.password, tt.firstName, tt.lastName)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, user)
			} else {
				require.NoError(t, err)
				require.NotNil(t, user)

				assert.Equal(t, tt.email, user.Email)
				assert.Equal(t, tt.firstName, user.FirstName)
				assert.Equal(t, tt.lastName, user.LastName)
				assert.Equal(t, domain.UserStatusPending, user.Status)
				assert.Equal(t, domain.UserRoleBuyer, user.Role)
				assert.False(t, user.EmailVerified)
				assert.False(t, user.PhoneVerified)
				assert.NotEmpty(t, user.ID)
				assert.NotEmpty(t, user.PasswordHash)
				assert.NotEmpty(t, user.VerificationToken)
				assert.NotZero(t, user.CreatedAt)
				assert.NotZero(t, user.UpdatedAt)
			}
		})
	}
}

func TestUser_ValidatePassword(t *testing.T) {
	user, err := domain.NewUser("test@example.com", "password123", "John", "Doe")
	require.NoError(t, err)

	// Test correct password
	err = user.ValidatePassword("password123")
	assert.NoError(t, err)

	// Test incorrect password
	err = user.ValidatePassword("wrongpassword")
	assert.Error(t, err)
}

func TestUser_VerifyEmail(t *testing.T) {
	user, err := domain.NewUser("test@example.com", "password123", "John", "Doe")
	require.NoError(t, err)

	// Initially not verified
	assert.False(t, user.EmailVerified)
	assert.Equal(t, domain.UserStatusPending, user.Status)
	assert.NotEmpty(t, user.VerificationToken)

	// Verify email
	user.VerifyEmail()

	assert.True(t, user.EmailVerified)
	assert.Equal(t, domain.UserStatusActive, user.Status)
	assert.Empty(t, user.VerificationToken)
}

func TestUser_UpgradeToSeller(t *testing.T) {
	user, err := domain.NewUser("test@example.com", "password123", "John", "Doe")
	require.NoError(t, err)

	// Cannot upgrade unverified user
	err = user.UpgradeToSeller("Test Business", "123 Business St")
	assert.Error(t, err)

	// Verify email first
	user.VerifyEmail()

	// Now upgrade to seller
	err = user.UpgradeToSeller("Test Business", "123 Business St")
	assert.NoError(t, err)

	assert.Equal(t, domain.UserRoleSeller, user.Role)
	require.NotNil(t, user.SellerProfile)
	assert.Equal(t, "Test Business", user.SellerProfile.BusinessName)
	assert.Equal(t, "123 Business St", user.SellerProfile.BusinessAddress)
	assert.Equal(t, domain.VerificationStatusPending, user.SellerProfile.VerificationStatus)

	// Cannot upgrade again
	err = user.UpgradeToSeller("Another Business", "456 Another St")
	assert.Error(t, err)
}

func TestUser_UpdateLastLogin(t *testing.T) {
	user, err := domain.NewUser("test@example.com", "password123", "John", "Doe")
	require.NoError(t, err)

	// Initially no last login
	assert.Nil(t, user.LastLoginAt)

	// Update last login
	beforeUpdate := time.Now()
	user.UpdateLastLogin()
	afterUpdate := time.Now()

	require.NotNil(t, user.LastLoginAt)
	assert.True(t, user.LastLoginAt.After(beforeUpdate) || user.LastLoginAt.Equal(beforeUpdate))
	assert.True(t, user.LastLoginAt.Before(afterUpdate) || user.LastLoginAt.Equal(afterUpdate))
}

func TestUser_BusinessLogic(t *testing.T) {
	user, err := domain.NewUser("test@example.com", "password123", "John", "Doe")
	require.NoError(t, err)

	// Test FullName
	assert.Equal(t, "John Doe", user.FullName())

	// Test IsActive
	assert.False(t, user.IsActive()) // Pending users are not active
	user.VerifyEmail()
	assert.True(t, user.IsActive())

	// Test IsSeller
	assert.False(t, user.IsSeller())
	user.UpgradeToSeller("Test Business", "123 Business St")
	assert.True(t, user.IsSeller())

	// Test IsVerifiedSeller
	assert.False(t, user.IsVerifiedSeller()) // Still pending verification
	user.SellerProfile.VerificationStatus = domain.VerificationStatusApproved
	assert.True(t, user.IsVerifiedSeller())

	// Test Suspend
	user.Suspend("Terms violation")
	assert.Equal(t, domain.UserStatusSuspended, user.Status)
	assert.False(t, user.IsActive())

	// Test Activate
	user.Activate()
	assert.Equal(t, domain.UserStatusActive, user.Status)
	assert.True(t, user.IsActive())
}