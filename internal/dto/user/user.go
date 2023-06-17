package user

type User struct {
	ID                string             `json:"id"`
	Email             string             `json:"email"`
	EmailVerified     bool               `json:"emailVerified"`
	SocialConnections []SocialConnection `json:"socialConnections"`
	Roles             []Role             `json:"roles"`
}

type SocialConnection struct {
	UserID           uint
	SocialProviderID uint
	SocialMediaID    string
}

type Role struct {
	Name   string
	UserID uint
}
