package user

type CreateUserInput struct {
	Email    string
	Password string
}

type UpdateUserInput struct {
	ID       string
	Email    *string
	Password *string
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
