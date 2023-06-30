package user

import (
	"errors"

	"gorm.io/gorm"
)

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

const (
	GoogleProviderID = 0
)
