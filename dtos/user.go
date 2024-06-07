package dtos

type CreateUserParams struct {
	Name            string `json:"name" binding:"required"`
	Email           string `json:"email" binding:"required"`
	Password        string `json:"password" binding:"required,eqfield=ConfirmPassword"`
	ConfirmPassword string `json:"confirmPassword" binding:"required"`
	Token           string `json:"token" binding:"required"`
}

type LoginParams struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	Token    string `json:"token" binding:"required"`
}

type GetUserFilter struct {
	ID       uint
	Email    string
	Password string
}
