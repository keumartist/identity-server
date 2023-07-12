package user

type CreateUserInput struct {
	Email    string
	Password string
}

type UpdateUserProfileInput struct {
	ID    string
	Email *string
	Name  *string
}

type GetUserByIDInput struct {
	ID string
}

type GetUserByEmailInput struct {
	Email string
}

type DeleteUserInput struct {
	ID string
}

type RefreshAccessTokenInput struct {
	Token string
}
