package user

import (
	userdomain "art-sso/internal/domain/user"
	"errors"
	"time"

	"gorm.io/gorm"
)

type MySQLUserRepository struct {
	db *gorm.DB
}

func NewMySQLUserRepository(db *gorm.DB) UserRepository {
	return &MySQLUserRepository{
		db: db,
	}
}

func (r *MySQLUserRepository) CreateUser(user *userdomain.User) error {
	normalRole := userdomain.Role{Name: userdomain.RoleNormal}
	user.Roles = append(user.Roles, normalRole)

	result := r.db.Create(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *MySQLUserRepository) GetUserByID(id string) (*userdomain.User, error) {
	var user userdomain.User
	result := r.db.Preload("Roles").Preload("SocialConnections").First(&user, id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, result.Error
		}
	}

	return &user, nil
}

func (r *MySQLUserRepository) GetUserByEmail(email string) (*userdomain.User, error) {
	var user userdomain.User
	result := r.db.Where("email = ?", email).First(&user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		} else {
			return nil, result.Error
		}
	}

	return &user, nil
}

func (r *MySQLUserRepository) UpdateUser(user *userdomain.User) error {
	result := r.db.Save(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *MySQLUserRepository) UpdateUserProfile(user *userdomain.User) error {
	result := r.db.Model(user).Updates(map[string]interface{}{
		"email": user.Email,
		"name":  user.Name,
	})

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *MySQLUserRepository) DeleteUser(user *userdomain.User) error {
	result := r.db.Delete(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *MySQLUserRepository) CreateUnverifiedUser(user *userdomain.User, verificationCode string, expireAt time.Time) error {
	user.VerificationCode = &verificationCode
	user.VerificationCodeExpireAt = &expireAt
	user.EmailVerified = false

	normalRole := userdomain.Role{Name: userdomain.RoleNormal}
	user.Roles = append(user.Roles, normalRole)

	result := r.db.Create(user)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *MySQLUserRepository) UpdateVerificationCode(user *userdomain.User, verificationCode string, expireAt time.Time) error {
	user.VerificationCode = &verificationCode
	user.VerificationCodeExpireAt = &expireAt
	user.EmailVerified = false

	result := r.db.Model(user).Updates(map[string]interface{}{
		"verification_code":           user.VerificationCode,
		"email_verified":              user.EmailVerified,
		"verification_code_expire_at": user.VerificationCodeExpireAt,
	})

	if result.Error != nil {
		return result.Error
	}

	return nil
}
func (r *MySQLUserRepository) UpdateRefreshToken(user *userdomain.User, refreshToken string) error {
	user.RefreshToken = refreshToken

	result := r.db.Model(user).Updates(map[string]interface{}{
		"refresh_token": refreshToken,
	})

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *MySQLUserRepository) VerifyUserEmail(user *userdomain.User) error {
	user.EmailVerified = true
	empty := ""
	user.VerificationCode = &empty
	user.VerificationCodeExpireAt = nil

	result := r.db.Save(user)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *MySQLUserRepository) AssignRoleToUser(user *userdomain.User, role userdomain.RoleName) error {
	if role != userdomain.RoleNormal && role != userdomain.RoleAdmin {
		return errors.New("Invalid role")
	}

	for _, existingRole := range user.Roles {
		if existingRole.Name == role {
			return errors.New("The user already has this role")
		}
	}

	newRole := userdomain.Role{Name: role, UserID: user.ID}
	user.Roles = append(user.Roles, newRole)

	result := r.db.Save(user)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
