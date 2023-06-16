package user

import (
	"errors"

	"gorm.io/gorm"
)

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
