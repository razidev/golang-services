package memberships

import "gorm.io/gorm"

type (
	User struct {
		gorm.Model
		Email     string `gorm:"unioque;not null"`
		Username  string `gorm:"unioque;not null"`
		Password  string `gorm:"not null"`
		CreatedBy string `gorm:"not null"`
		UpdateBy  string `gorm:"not null"`
	}
)

type (
	SignUpRequest struct {
		Email    string `json:"email"`
		Username string `json:"username"`
		Password string `json:"password"`
	}

	LoginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
)

type (
	LoginResponse struct {
		AccessToken string `json:"accessToken"`
	}
)
