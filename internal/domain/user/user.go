package domain

import (
	"errors"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID                uint   `gorm:"primaryKey"`
	Email             string `gorm:"type:varchar(256);unique"`
	Password          string
	Salt              string
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

type Role struct {
	gorm.Model
	Name   string `gorm:"type:varchar(256)"`
	UserID uint
}

func (r *Role) AssignToUser(user *User) error {
	for _, existingRole := range user.Roles {
		if existingRole.Name == r.Name {
			return errors.New("The user already has this role")
		}
	}

	user.Roles = append(user.Roles, *r)
	return nil
}

type SocialProvider struct {
	gorm.Model
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"type:varchar(256)"`
}

type SocialConnection struct {
	gorm.Model
	UserID           uint
	SocialProviderID uint
	SocialMediaID    string `gorm:"type:varchar(256)"`
}

func (sc *SocialConnection) Connect(user *User, provider *SocialProvider) error {
	for _, existingConnection := range user.SocialConnections {
		if existingConnection.SocialProviderID == provider.ID {
			return errors.New("The user is already connected to this social media account")
		}
	}

	user.SocialConnections = append(user.SocialConnections, *sc)
	return nil
}
