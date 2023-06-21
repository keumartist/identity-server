package user

import (
	userdomain "art-sso/internal/domain/user"
	"errors"

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
	result := r.db.Create(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *MySQLUserRepository) GetUserByID(id string) (*userdomain.User, error) {
	var user userdomain.User
	result := r.db.First(&user, id)

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
			return nil, nil
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

func (r *MySQLUserRepository) CreateUnverifiedUser(user *userdomain.User, verificationCode string) error {
	user.VerificationCode = &verificationCode
	user.EmailVerified = false

	result := r.db.Create(user)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *MySQLUserRepository) UpdateVerificationCode(user *userdomain.User, verificationCode string) error {
	user.VerificationCode = &verificationCode
	user.EmailVerified = false

	result := r.db.Model(user).Updates(map[string]interface{}{
		"verification_code": user.VerificationCode,
		"email_verified":    user.EmailVerified,
	})

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *MySQLUserRepository) VerifyUser(email, verificationCode string) error {
	var user userdomain.User
	result := r.db.Where("email = ? AND verification_code = ?", email, verificationCode).First(&user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return errors.New("Invalid email or verification code")
		} else {
			return result.Error
		}
	}

	user.EmailVerified = true

	empty := ""
	user.VerificationCode = &empty

	updateResult := r.db.Save(&user)
	if updateResult.Error != nil {
		return updateResult.Error
	}
	return nil
}
