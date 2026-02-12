package user

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

type RepositoryInterface interface {
	GetById(ctx context.Context, id uint) (*User, error)
	GetAll(ctx context.Context) ([]User, error)
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
	UpdateFields(ctx context.Context, id uint, updates map[string]interface{}) error
	Delete(ctx context.Context, id uint) error
	FindByEmail(ctx context.Context, email string) (*User, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	ExistsByEmailExcludingId(ctx context.Context, email string, userId uint) (bool, error)
	GetWithPagination(ctx context.Context, page, pageSize int) ([]User, int64, error)
	Search(ctx context.Context, keyword string) ([]User, error)
	SearchWithPagination(ctx context.Context, keyword string, page, pageSize int) ([]User, int64, error)
	Count(ctx context.Context) (int64, error)
}

type Repository struct {
	DB *gorm.DB
}

func NewRepository(db *gorm.DB) RepositoryInterface {
	return &Repository{
		DB: db,
	}
}

// Get user by Id
func (r *Repository) GetById(ctx context.Context, id uint) (*User, error) {
	var user User
	err := r.DB.WithContext(ctx).First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("User not found")
		}
		return nil, err
	}
	return &user, err
}

// Get all Users
func (r *Repository) GetAll(ctx context.Context) ([]User, error) {
	var users []User
	err := r.DB.WithContext(ctx).Order("id DESC").Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

// Create a new user
func (r *Repository) Create(ctx context.Context, user *User) error {
	return r.DB.WithContext(ctx).Create(user).Error
}

// Update a existing user
func (r *Repository) Update(ctx context.Context, user *User) error {
	return r.DB.WithContext(ctx).Save(user).Error
}

// Update specific fields in existing user
func (r *Repository) UpdateFields(ctx context.Context, id uint, updates map[string]interface{}) error {
	result := r.DB.WithContext(ctx).Model(&User{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("User not Found")
	}
	return nil
}

// Delete - Delete user by ID (soft delete)
func (r *Repository) Delete(ctx context.Context, id uint) error {
	result := r.DB.WithContext(ctx).Delete(&User{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}
	return nil
}

// Get user by email
func (r *Repository) FindByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	err := r.DB.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

// Check if email exists or not
func (r *Repository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	err := r.DB.WithContext(ctx).Model(&User{}).
		Where("email = ?", email).
		Count(&count).Error
	return count > 0, err
}

// Check if email exists excluding specific user ID
func (r *Repository) ExistsByEmailExcludingId(ctx context.Context, email string, userId uint) (bool, error) {
	var count int64
	err := r.DB.WithContext(ctx).Model(&User{}).
		Where("email = ? AND id != ?", email, userId).
		Count(&count).Error
	return count > 0, err
}

// Get users with pagination
func (r *Repository) GetWithPagination(ctx context.Context, page, pageSize int) ([]User, int64, error) {
	var users []User
	var total int64

	offset := (page - 1) * pageSize

	if err := r.DB.WithContext(ctx).Model(&User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.DB.WithContext(ctx).
		Offset(offset).
		Limit(pageSize).
		Order("id DESC").
		Find(&users).Error
	if err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

// Search users by first name, last name, or email
func (r *Repository) Search(ctx context.Context, keyword string) ([]User, error) {
	var users []User
	searchPattern := "%" + keyword + "%"
	err := r.DB.WithContext(ctx).
		Where("first_name ILIKE ? OR last_name ILIKE ? OR email ILIKE ?",
			searchPattern, searchPattern, searchPattern).
		Order("id DESC").
		Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

// SearchWithPagination - Search with pagination
func (r *Repository) SearchWithPagination(ctx context.Context, keyword string, page, pageSize int) ([]User, int64, error) {
	var users []User
	var total int64

	offset := (page - 1) * pageSize
	searchPattern := "%" + keyword + "%"

	query := r.DB.WithContext(ctx).Where(
		"first_name ILIKE ? OR last_name ILIKE ? OR email ILIKE ?",
		searchPattern, searchPattern, searchPattern,
	)

	if err := query.Model(&User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// Get total user count
func (r *Repository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.DB.WithContext(ctx).Model(&User{}).Count(&count).Error
	return count, err
}
