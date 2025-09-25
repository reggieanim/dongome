package infra

import (
	"net/http"

	"dongome/internal/users/app"
	"dongome/pkg/errors"

	"github.com/gin-gonic/gin"
)

// UserHandler handles HTTP requests for users
type UserHandler struct {
	userService *app.UserService
}

// NewUserHandler creates a new user handler
func NewUserHandler(userService *app.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// RegisterRoutes registers user routes
func (h *UserHandler) RegisterRoutes(r *gin.RouterGroup) {
	users := r.Group("/users")
	{
		users.POST("/register", h.RegisterUser)
		users.POST("/login", h.LoginUser)
		users.POST("/verify-email", h.VerifyEmail)
		users.POST("/:id/upgrade-to-seller", h.UpgradeToSeller)
		users.GET("/:id", h.GetUser)
	}
}

// RegisterUser handles user registration
func (h *UserHandler) RegisterUser(c *gin.Context) {
	var cmd app.RegisterUserCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userService.RegisterUser(c.Request.Context(), cmd)
	if err != nil {
		if domainErr, ok := err.(*errors.DomainError); ok {
			c.JSON(domainErr.HTTPStatusCode(), gin.H{"error": domainErr.Message, "code": domainErr.Code})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	// Don't return sensitive information
	response := gin.H{
		"id":             user.ID,
		"email":          user.Email,
		"first_name":     user.FirstName,
		"last_name":      user.LastName,
		"status":         user.Status,
		"role":           user.Role,
		"email_verified": user.EmailVerified,
		"phone_verified": user.PhoneVerified,
		"created_at":     user.CreatedAt,
	}

	c.JSON(http.StatusCreated, response)
}

// LoginUser handles user login
func (h *UserHandler) LoginUser(c *gin.Context) {
	var cmd app.LoginCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userService.LoginUser(c.Request.Context(), cmd)
	if err != nil {
		if domainErr, ok := err.(*errors.DomainError); ok {
			c.JSON(domainErr.HTTPStatusCode(), gin.H{"error": domainErr.Message, "code": domainErr.Code})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	// In a real app, you'd generate and return a JWT token here
	response := gin.H{
		"id":             user.ID,
		"email":          user.Email,
		"first_name":     user.FirstName,
		"last_name":      user.LastName,
		"status":         user.Status,
		"role":           user.Role,
		"email_verified": user.EmailVerified,
		"phone_verified": user.PhoneVerified,
		"last_login_at":  user.LastLoginAt,
		"seller_profile": user.SellerProfile,
	}

	c.JSON(http.StatusOK, response)
}

// VerifyEmail handles email verification
func (h *UserHandler) VerifyEmail(c *gin.Context) {
	type VerifyEmailRequest struct {
		Token string `json:"token" binding:"required"`
	}

	var req VerifyEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.userService.VerifyEmail(c.Request.Context(), req.Token)
	if err != nil {
		if domainErr, ok := err.(*errors.DomainError); ok {
			c.JSON(domainErr.HTTPStatusCode(), gin.H{"error": domainErr.Message, "code": domainErr.Code})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "email verified successfully"})
}

// UpgradeToSeller handles user upgrade to seller
func (h *UserHandler) UpgradeToSeller(c *gin.Context) {
	userID := c.Param("id")

	var cmd app.UpgradeToSellerCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cmd.UserID = userID

	err := h.userService.UpgradeToSeller(c.Request.Context(), cmd)
	if err != nil {
		if domainErr, ok := err.(*errors.DomainError); ok {
			c.JSON(domainErr.HTTPStatusCode(), gin.H{"error": domainErr.Message, "code": domainErr.Code})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "successfully upgraded to seller"})
}

// GetUser handles getting user by ID
func (h *UserHandler) GetUser(c *gin.Context) {
	userID := c.Param("id")

	user, err := h.userService.GetUser(c.Request.Context(), userID)
	if err != nil {
		if domainErr, ok := err.(*errors.DomainError); ok {
			c.JSON(domainErr.HTTPStatusCode(), gin.H{"error": domainErr.Message, "code": domainErr.Code})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	// Don't return sensitive information
	response := gin.H{
		"id":             user.ID,
		"email":          user.Email,
		"first_name":     user.FirstName,
		"last_name":      user.LastName,
		"phone_number":   user.PhoneNumber,
		"avatar":         user.Avatar,
		"status":         user.Status,
		"role":           user.Role,
		"email_verified": user.EmailVerified,
		"phone_verified": user.PhoneVerified,
		"last_login_at":  user.LastLoginAt,
		"created_at":     user.CreatedAt,
		"seller_profile": user.SellerProfile,
	}

	c.JSON(http.StatusOK, response)
}
