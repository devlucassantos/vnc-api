package swagger

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	Id           uuid.UUID `json:"id"            example:"54b094b23-f426-89d2-843e-1335410b18df"`
	FirstName    string    `json:"first_name"    example:"Lucas"`
	LastName     string    `json:"last_name"     example:"Santos"`
	Email        string    `json:"email"         example:"example@email.com"`
	Roles        []string  `json:"roles"         example:"user"`
	AccessToken  string    `json:"access_token"  example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InByb2YiL..."`
	RefreshToken string    `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InByb2YiL..."`
	CreatedAt    time.Time `json:"created_at"    example:"2024-01-06T14:52:23.723959Z"`
	UpdatedAt    time.Time `json:"updated_at"    example:"2024-01-06T16:02:29.854799Z"`
}
