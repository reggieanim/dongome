package domain

import (
	"time"

	"dongome/pkg/errors"

	"github.com/google/uuid"
)

// ListingStatus represents the status of a listing
type ListingStatus string

const (
	ListingStatusDraft    ListingStatus = "draft"
	ListingStatusActive   ListingStatus = "active"
	ListingStatusInactive ListingStatus = "inactive"
	ListingStatusSold     ListingStatus = "sold"
	ListingStatusExpired  ListingStatus = "expired"
)

// Condition represents the condition of an item
type Condition string

const (
	ConditionNew      Condition = "new"
	ConditionLikeNew  Condition = "like_new"
	ConditionGood     Condition = "good"
	ConditionFair     Condition = "fair"
	ConditionPoor     Condition = "poor"
	ConditionForParts Condition = "for_parts"
)

// Category represents a listing category
type Category struct {
	ID          string     `gorm:"type:uuid;primary_key" json:"id"`
	Name        string     `gorm:"not null" json:"name"`
	Description string     `json:"description"`
	ParentID    *string    `gorm:"type:uuid" json:"parent_id,omitempty"`
	Parent      *Category  `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Children    []Category `gorm:"foreignKey:ParentID" json:"children,omitempty"`
	ImageURL    string     `json:"image_url"`
	IsActive    bool       `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// Listing represents a marketplace listing aggregate root
type Listing struct {
	ID             string             `gorm:"type:uuid;primary_key" json:"id"`
	SellerID       string             `gorm:"type:uuid;not null;index" json:"seller_id"`
	CategoryID     string             `gorm:"type:uuid;not null" json:"category_id"`
	Category       Category           `gorm:"foreignKey:CategoryID" json:"category"`
	Title          string             `gorm:"not null" json:"title"`
	Description    string             `gorm:"type:text" json:"description"`
	Price          float64            `gorm:"not null" json:"price"`
	Currency       string             `gorm:"default:'GHS'" json:"currency"`
	Condition      Condition          `gorm:"not null" json:"condition"`
	Status         ListingStatus      `gorm:"default:'draft'" json:"status"`
	Location       Location           `gorm:"embedded" json:"location"`
	Images         []ListingImage     `gorm:"foreignKey:ListingID" json:"images"`
	Attributes     []ListingAttribute `gorm:"foreignKey:ListingID" json:"attributes"`
	Tags           []ListingTag       `gorm:"many2many:listing_tags;" json:"tags"`
	ViewsCount     int                `gorm:"default:0" json:"views_count"`
	FavoritesCount int                `gorm:"default:0" json:"favorites_count"`
	IsNegotiable   bool               `gorm:"default:true" json:"is_negotiable"`
	IsPromoted     bool               `gorm:"default:false" json:"is_promoted"`
	PromotedUntil  *time.Time         `json:"promoted_until,omitempty"`
	ExpiresAt      time.Time          `json:"expires_at"`
	CreatedAt      time.Time          `json:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at"`
}

// Location represents geographical location
type Location struct {
	Region    string  `gorm:"not null" json:"region"`
	City      string  `gorm:"not null" json:"city"`
	Area      string  `json:"area"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// ListingImage represents a listing image
type ListingImage struct {
	ID        string    `gorm:"type:uuid;primary_key" json:"id"`
	ListingID string    `gorm:"type:uuid;not null" json:"listing_id"`
	URL       string    `gorm:"not null" json:"url"`
	Caption   string    `json:"caption"`
	Order     int       `gorm:"default:0" json:"order"`
	CreatedAt time.Time `json:"created_at"`
}

// ListingAttribute represents dynamic attributes for listings
type ListingAttribute struct {
	ID        string    `gorm:"type:uuid;primary_key" json:"id"`
	ListingID string    `gorm:"type:uuid;not null" json:"listing_id"`
	Key       string    `gorm:"not null" json:"key"`
	Value     string    `gorm:"not null" json:"value"`
	CreatedAt time.Time `json:"created_at"`
}

// ListingTag represents tags for listings
type ListingTag struct {
	ID        string    `gorm:"type:uuid;primary_key" json:"id"`
	Name      string    `gorm:"uniqueIndex;not null" json:"name"`
	Color     string    `json:"color"`
	CreatedAt time.Time `json:"created_at"`
}

// NewListing creates a new listing
func NewListing(sellerID, categoryID, title, description string, price float64, condition Condition, location Location) (*Listing, error) {
	if sellerID == "" {
		return nil, errors.ValidationError("seller ID is required")
	}
	if categoryID == "" {
		return nil, errors.ValidationError("category ID is required")
	}
	if title == "" {
		return nil, errors.ValidationError("title is required")
	}
	if price <= 0 {
		return nil, errors.ValidationError("price must be greater than 0")
	}

	// Set expiration to 30 days from now
	expiresAt := time.Now().AddDate(0, 0, 30)

	return &Listing{
		ID:             uuid.New().String(),
		SellerID:       sellerID,
		CategoryID:     categoryID,
		Title:          title,
		Description:    description,
		Price:          price,
		Currency:       "GHS",
		Condition:      condition,
		Status:         ListingStatusDraft,
		Location:       location,
		Images:         []ListingImage{},
		Attributes:     []ListingAttribute{},
		Tags:           []ListingTag{},
		ViewsCount:     0,
		FavoritesCount: 0,
		IsNegotiable:   true,
		IsPromoted:     false,
		ExpiresAt:      expiresAt,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}, nil
}

// Activate activates the listing
func (l *Listing) Activate() error {
	if l.Status == ListingStatusSold {
		return errors.ValidationError("cannot activate sold listing")
	}

	l.Status = ListingStatusActive
	l.UpdatedAt = time.Now()
	return nil
}

// Deactivate deactivates the listing
func (l *Listing) Deactivate() {
	l.Status = ListingStatusInactive
	l.UpdatedAt = time.Now()
}

// MarkAsSold marks the listing as sold
func (l *Listing) MarkAsSold() {
	l.Status = ListingStatusSold
	l.UpdatedAt = time.Now()
}

// IncrementViews increments the view count
func (l *Listing) IncrementViews() {
	l.ViewsCount++
	l.UpdatedAt = time.Now()
}

// IncrementFavorites increments the favorites count
func (l *Listing) IncrementFavorites() {
	l.FavoritesCount++
	l.UpdatedAt = time.Now()
}

// DecrementFavorites decrements the favorites count
func (l *Listing) DecrementFavorites() {
	if l.FavoritesCount > 0 {
		l.FavoritesCount--
		l.UpdatedAt = time.Now()
	}
}

// AddImage adds an image to the listing
func (l *Listing) AddImage(url, caption string) {
	image := ListingImage{
		ID:        uuid.New().String(),
		ListingID: l.ID,
		URL:       url,
		Caption:   caption,
		Order:     len(l.Images),
		CreatedAt: time.Now(),
	}
	l.Images = append(l.Images, image)
	l.UpdatedAt = time.Now()
}

// AddAttribute adds an attribute to the listing
func (l *Listing) AddAttribute(key, value string) {
	attribute := ListingAttribute{
		ID:        uuid.New().String(),
		ListingID: l.ID,
		Key:       key,
		Value:     value,
		CreatedAt: time.Now(),
	}
	l.Attributes = append(l.Attributes, attribute)
	l.UpdatedAt = time.Now()
}

// Promote promotes the listing
func (l *Listing) Promote(duration time.Duration) {
	l.IsPromoted = true
	promotedUntil := time.Now().Add(duration)
	l.PromotedUntil = &promotedUntil
	l.UpdatedAt = time.Now()
}

// IsExpired checks if the listing has expired
func (l *Listing) IsExpired() bool {
	return time.Now().After(l.ExpiresAt)
}

// IsActive checks if the listing is active
func (l *Listing) IsActive() bool {
	return l.Status == ListingStatusActive && !l.IsExpired()
}

// ListingRepository defines the interface for listing persistence
type ListingRepository interface {
	Save(listing *Listing) error
	FindByID(id string) (*Listing, error)
	FindBySeller(sellerID string, limit, offset int) ([]*Listing, error)
	FindByCategory(categoryID string, limit, offset int) ([]*Listing, error)
	Search(query string, filters map[string]interface{}, limit, offset int) ([]*Listing, error)
	Update(listing *Listing) error
	Delete(id string) error
}

// CategoryRepository defines the interface for category persistence
type CategoryRepository interface {
	Save(category *Category) error
	FindByID(id string) (*Category, error)
	FindAll() ([]*Category, error)
	FindByParent(parentID string) ([]*Category, error)
	Update(category *Category) error
	Delete(id string) error
}
