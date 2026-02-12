package user

import (
	"context"
	"errors"

	"github.com/SalmanKhandakarhub/scalable-backend-architecture/pkg/utils"
)

type ServiceInterface interface {
	GetById(ctx context.Context, id uint) (*User, error)
	GetAll(ctx context.Context) ([]User, error)
	Create(ctx context.Context, req CreateUserRequest) (*User, error)
	Update(ctx context.Context, id uint, req UpdateUserRequest) (*User, error)
	Delete(ctx context.Context, id uint) error
	Login(ctx context.Context, req LoginRequest) (*User, string, error)
	ChangePassword(ctx context.Context, id uint, req ChangePasswordRequest) error
	GetWithPagination(ctx context.Context, page, pageSize int) ([]User, int64, error)
	Search(ctx context.Context, keyword string) ([]User, error)
	SearchWithPagination(ctx context.Context, keyword string, page, pageSize int) ([]User, int64, error)
	GetStats(ctx context.Context) (map[string]any, error)
}

type Service struct {
	repo RepositoryInterface
}

func NewService(repo RepositoryInterface) ServiceInterface {
	return &Service{
		repo: repo,
	}
}

func (s *Service) GetById(ctx context.Context, id uint) (*User, error) {
	return s.repo.GetById(ctx, id)
}

func (s *Service) GetAll(ctx context.Context) ([]User, error) {
	return s.repo.GetAll(ctx)
}

// Create - Create new user with validation
func (s *Service) Create(ctx context.Context, req CreateUserRequest) (*User, error) {
	// Check if email already exists
	exists, err := s.repo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("email already exists")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// Create user
	user := &User{
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Email:        req.Email,
		Password:     hashedPassword,
		Age:          req.Age,
		ContactNo:    req.ContactNo,
		ProfileImage: req.ProfileImage,
		Address:      req.Address,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// Update - Update user
func (s *Service) Update(ctx context.Context, id uint, req UpdateUserRequest) (*User, error) {
	// Get existing user
	user, err := s.repo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update fields only if provided
	if req.FirstName != nil {
		user.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		user.LastName = *req.LastName
	}
	if req.Age != nil {
		user.Age = req.Age
	}
	if req.ContactNo != nil {
		user.ContactNo = req.ContactNo
	}
	if req.ProfileImage != nil {
		user.ProfileImage = req.ProfileImage
	}
	if req.Address != nil {
		user.Address = req.Address
	}

	if err := s.repo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// Delete - Delete user
func (s *Service) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}

// Login - Authenticate user
func (s *Service) Login(ctx context.Context, req LoginRequest) (*User, string, error) {
	// Find user by email
	user, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, "", errors.New("invalid email or password")
	}

	// Check password
	if !utils.CheckPassword(req.Password, user.Password) {
		return nil, "", errors.New("invalid email or password")
	}

	// Generate JWT token
	token, err := utils.GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, "", errors.New("failed to generate token")
	}

	return user, token, nil
}

// ChangePassword - Change user password
func (s *Service) ChangePassword(ctx context.Context, id uint, req ChangePasswordRequest) error {
	// Get user
	user, err := s.repo.GetById(ctx, id)
	if err != nil {
		return err
	}

	// Verify old password
	if !utils.CheckPassword(req.OldPassword, user.Password) {
		return errors.New("invalid old password")
	}

	// Hash new password
	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return errors.New("failed to hash password")
	}

	// Update password
	updates := map[string]interface{}{
		"password": hashedPassword,
	}

	return s.repo.UpdateFields(ctx, id, updates)
}

// GetWithPagination - Get users with pagination
func (s *Service) GetWithPagination(ctx context.Context, page, pageSize int) ([]User, int64, error) {
	return s.repo.GetWithPagination(ctx, page, pageSize)
}

// Search - Search users
func (s *Service) Search(ctx context.Context, keyword string) ([]User, error) {
	return s.repo.Search(ctx, keyword)
}

// SearchWithPagination - Search with pagination
func (s *Service) SearchWithPagination(ctx context.Context, keyword string, page, pageSize int) ([]User, int64, error) {
	return s.repo.SearchWithPagination(ctx, keyword, page, pageSize)
}

// GetStats - Get user statistics
func (s *Service) GetStats(ctx context.Context) (map[string]any, error) {
	total, err := s.repo.Count(ctx)
	if err != nil {
		return nil, err
	}

	stats := map[string]any{
		"total_users": total,
	}

	return stats, nil
}
