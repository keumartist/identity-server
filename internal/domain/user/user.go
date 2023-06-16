package user

import (
	"errors"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID                uint   `gorm:"primaryKey"`
	Email             string `gorm:"type:varchar(256);unique"`
	Password          string
	RefreshToken      string
	VerificationCode  *string
	EmailVerified     bool               `gorm:"default:false"`
	SocialConnections []SocialConnection `gorm:"foreignkey:UserID"`
	Roles             []Role             `gorm:"foreignkey:UserID"`
}

func (u *User) ChangePassword(oldPassword, newPassword string) error {
	if u.Password != oldPassword {
		return errors.New("The provided old password does not match The current password")
	}

	u.Password = newPassword
	return nil
}

func (u *User) VerifyEmail(code string) error {
	if u.VerificationCode == nil || *u.VerificationCode != code {
		return errors.New("The provided verification code does not match")
	}

	u.EmailVerified = true
	return nil
}
