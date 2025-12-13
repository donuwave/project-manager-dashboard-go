package dto

type UserDTO struct {
	ID      string `json:"id" format:"uuid"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Country string `json:"country,omitempty"`
}

type CreateUserRequest struct {
	Email   string `json:"email" example:"john@doe.com"`
	Name    string `json:"name" example:"John"`
	Country string `json:"country,omitempty" example:"DE"`
}

type ErrorDTO struct {
	Message string `json:"message"`
}
