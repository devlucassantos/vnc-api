package response

import (
	"github.com/devlucassantos/vnc-domains/src/domains/user"
	"github.com/google/uuid"
	"time"
)

type User struct {
	Id           uuid.UUID `json:"id"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	Email        string    `json:"email"`
	Roles        []string  `json:"roles"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func NewUser(userData user.User) *User {
	var roles []string
	for _, roleData := range userData.Roles() {
		roles = append(roles, roleData.Code())
	}

	return &User{
		Id:           userData.Id(),
		FirstName:    userData.FirstName(),
		LastName:     userData.LastName(),
		Email:        userData.Email(),
		Roles:        roles,
		AccessToken:  userData.AccessToken(),
		RefreshToken: userData.RefreshToken(),
		CreatedAt:    userData.CreatedAt(),
		UpdatedAt:    userData.UpdatedAt(),
	}
}
