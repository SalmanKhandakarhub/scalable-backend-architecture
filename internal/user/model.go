package user

import (
	"strings"

	"github.com/SalmanKhandakarhub/scalable-backend-architecture/internal/models"
)

type User struct {
	models.BaseModel
	FirstName    string  `gorm:"type:varchar(64);not null" json:"first_name" binding:"required"`
	LastName     string  `gorm:"type:varchar(64);not null" json:"last_name" binding:"required"`
	Email        string  `gorm:"type:varchar(255);uniqueIndex;not null" json:"email" binding:"required,email"`
	Password     string  `gorm:"type:varchar(255);not null" json:"-"`
	Age          *int    `gorm:"type:int" json:"age,omitempty" binding:"omitempty,gte=0,lte=150"`
	ContactNo    *string `gorm:"type:varchar(20)" json:"contact_no,omitempty"`
	ProfileImage *string `gorm:"type:varchar(255)" json:"profile_image,omitempty"`
	Address      *string `gorm:"type:text" json:"address,omitempty"`
}

func (User) TableName() string {
	return "users"
}

// FullName is not stored in DB, computed on the fly
func (u *User) FullName() string {
	return strings.TrimSpace(u.FirstName + " " + u.LastName)
}
