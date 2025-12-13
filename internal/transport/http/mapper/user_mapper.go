package mapper

import (
	"project-manager-dashboard-go/internal/app/usecase/user"
	"project-manager-dashboard-go/internal/transport/http/dto"
)

func ToUserDTO(u user.User) dto.UserDTO {
	return dto.UserDTO{
		ID:      u.ID.String(),
		Email:   u.Email,
		Name:    u.Name,
		Country: u.Country,
	}
}
