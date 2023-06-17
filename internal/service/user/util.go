package user

import (
	domain "art-sso/internal/domain/user"
	dto "art-sso/internal/dto/user"
	"fmt"
)

func UserDomainToDto(user *domain.User) dto.User {
	var socialConnections []dto.SocialConnection
	for _, sc := range user.SocialConnections {
		socialConnections = append(socialConnections, dto.SocialConnection{
			UserID:           sc.UserID,
			SocialProviderID: sc.SocialProviderID,
			SocialMediaID:    sc.SocialMediaID,
		})
	}

	var roles []dto.Role
	for _, role := range user.Roles {
		roles = append(roles, dto.Role{
			Name:   role.Name,
			UserID: role.UserID,
		})
	}

	return dto.User{
		ID:                fmt.Sprintf("%d", user.ID),
		Email:             user.Email,
		EmailVerified:     user.EmailVerified,
		SocialConnections: socialConnections,
		Roles:             roles,
	}
}
