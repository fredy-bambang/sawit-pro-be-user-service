package repository

import (
	"github.com/SawitProRecruitment/UserService/models"
	"github.com/jinzhu/gorm"
)

type PgUserRepository struct {
	DB *gorm.DB
}

// UserRepository is an interface for user repository
type UserRepository interface {
	Create(user *models.User) error
	FindByPhone(phone string) (*models.User, error)
	FindByID(id int) (*models.User, error)
	Update(user *models.User) error
}

// Create creates a new user
func (r *PgUserRepository) Create(user *models.User) error {
	return r.DB.Create(user).Error
}

// FindByPhone finds a user by phone number
func (r *PgUserRepository) FindByPhone(phone string) (*models.User, error) {
	var user models.User
	err := r.DB.Where("phone_number = ?", phone).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByID finds a user by id
func (r *PgUserRepository) FindByID(id int) (*models.User, error) {
	var user models.User
	err := r.DB.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update updates a user
func (r *PgUserRepository) Update(user *models.User) error {
	return r.DB.Save(user).Error
}

// NewPgUserRepository creates new postgress user repository
func NewPgUserRepository(db *gorm.DB) *PgUserRepository {
	return &PgUserRepository{DB: db}
}
