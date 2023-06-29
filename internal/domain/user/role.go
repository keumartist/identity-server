package user

import (
	"gorm.io/gorm"
)

type RoleName string

const (
	RoleNormal RoleName = "normal"
	RoleAdmin  RoleName = "admin"
)

type Role struct {
	gorm.Model
	Name   RoleName `gorm:"type:varchar(256)"`
	UserID uint
}
