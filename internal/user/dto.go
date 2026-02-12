package user

import "time"

// Request DTO for creating user
type CreateUserRequest struct {
	FirstName    string  `json:"first_name" binding:"required, min=2, max=64"`
	LastName     string  `json:"last_name" binding:"required,min=2,max=64"`
	Email        string  `json:"email" binding:"required,email"`
	Password     string  `json:"password" binding:"required,min=6"`
	Age          *int    `json:"age,omitempty" binding:"omitempty,gte=0,lte=150"`
	ContactNo    *string `json:"contact_no,omitempty" binding:"omitempty,min=10,max=20"`
	ProfileImage *string `json:"profile_image,omitempty"`
	Address      *string `json:"address,omitempty"`
}

// Request DTO for updating user
type UpdateUserRequest struct {
	FirstName    *string `json:"first_name,omitempty" binding:"omitempty,min=2,max=64"`
	LastName     *string `json:"last_name,omitempty" binding:"omitempty,min=2,max=64"`
	Age          *int    `json:"age,omitempty" binding:"omitempty,gte=0,lte=150"`
	ContactNo    *string `json:"contact_no,omitempty" binding:"omitempty,min=10,max=20"`
	ProfileImage *string `json:"profile_image,omitempty"`
	Address      *string `json:"address,omitempty"`
}

// Request DTO for change password
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required, min=6"`
}

// Request DTO for login
type LoginRequest struct {
	Email    string `json:"email" binding:"required, email"`
	Password string `json:"password" binding:"required, min=6"`
}

// Response DTO for user
type UserResponse struct {
	ID           uint      `json:"id"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	FullName     string    `json:"full_name"`
	Email        string    `json:"email"`
	Age          *int      `json:"age,omitempty"`
	ContactNo    *string   `json:"contact_no,omitempty"`
	ProfileImage *string   `json:"profile_image,omitempty"`
	Address      *string   `json:"address,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Response DTO for login
type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

// PaginationQuery - Query parameters for pagination
type PaginationQuery struct {
	Page     int `form:"page,default=1" binding:"omitempty,min=1"`
	PageSize int `form:"page_size,default=10" binding:"omitempty,min=1,max=100"`
}

// SearchQuery - Query parameters for search
type SearchQuery struct {
	Keyword string `form:"keyword" binding:"omitempty,min=1"`
	PaginationQuery
}

// ToUserResponse - Convert User model to UserResponse
func ToUserResponse(user *User) UserResponse {
	return UserResponse{
		ID:           user.ID,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		FullName:     user.FullName(),
		Email:        user.Email,
		Age:          user.Age,
		ContactNo:    user.ContactNo,
		ProfileImage: user.ProfileImage,
		Address:      user.Address,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}
}

// ToUserResponseList - Convert User list to UserResponse list
func ToUserResponseList(users []User) []UserResponse {
	responses := make([]UserResponse, len(users))
	for i, user := range users {
		responses[i] = ToUserResponse(&user)
	}
	return responses
}
