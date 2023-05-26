package domain

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID               uint   `gorm:"primaryKey"`
	Email            string `gorm:"type:varchar(256);unique"`
	Password         string
	Salt             string
	RefreshToken     string
	VerificationCode *string
	EmailVerified    bool `gorm:"default:false"`
	SocialProvider   *SocialProvider
	Roles            []Role `gorm:"foreignkey:UserID"`
}

type Role struct {
	gorm.Model
	Name   string `gorm:"type:varchar(256)"`
	UserID string `gorm:"type:varchar(256)"`
}

type SocialProvider struct {
	gorm.Model
	Name   string `gorm:"type:varchar(256)"`
	UserID string `gorm:"type:varchar(256)"`
}
