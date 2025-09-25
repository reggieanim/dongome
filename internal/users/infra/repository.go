package infra

import (
	"dongome/internal/users/domain"
	"dongome/pkg/errors"

	"gorm.io/gorm"
)

// UserGORMRepository implements UserRepository using GORM
type UserGORMRepository struct {
	db *gorm.DB
}

// NewUserGORMRepository creates a new user repository
func NewUserGORMRepository(db *gorm.DB) *UserGORMRepository {
	return &UserGORMRepository{
		db: db,
	}
}

// Save saves a user to the database
func (r *UserGORMRepository) Save(user *domain.User) error {
	if err := r.db.Create(user).Error; err != nil {
		return err
	}
	return nil
}

// FindByID finds a user by ID
func (r *UserGORMRepository) FindByID(id string) (*domain.User, error) {
	var user domain.User
	err := r.db.Preload("SellerProfile").First(&user, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NotFoundError("user not found")
		}
		return nil, err
	}
	return &user, nil
}

// FindByEmail finds a user by email
func (r *UserGORMRepository) FindByEmail(email string) (*domain.User, error) {
	var user domain.User
	err := r.db.Preload("SellerProfile").First(&user, "email = ?", email).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NotFoundError("user not found")
		}
		return nil, err
	}
	return &user, nil
}

// FindByVerificationToken finds a user by verification token
func (r *UserGORMRepository) FindByVerificationToken(token string) (*domain.User, error) {
	var user domain.User
	err := r.db.First(&user, "verification_token = ?", token).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NotFoundError("user not found")
		}
		return nil, err
	}
	return &user, nil
}

// Update updates a user in the database
func (r *UserGORMRepository) Update(user *domain.User) error {
	return r.db.Session(&gorm.Session{FullSaveAssociations: true}).Save(user).Error
}

// Delete deletes a user from the database
func (r *UserGORMRepository) Delete(id string) error {
	return r.db.Delete(&domain.User{}, "id = ?", id).Error
}
