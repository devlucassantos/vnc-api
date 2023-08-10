package response

import (
	"github.com/google/uuid"
	"time"
	"vnc-read-api/core/domains/organization"
)

type Organization struct {
	Id        uuid.UUID `json:"id"`
	Code      int       `json:"code"`
	Name      string    `json:"name"`
	Acronym   string    `json:"acronym"`
	Nickname  string    `json:"nickname"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewOrganization(organization organization.Organization) *Organization {
	return &Organization{
		Id:        organization.Id(),
		Code:      organization.Code(),
		Name:      organization.Name(),
		Acronym:   organization.Acronym(),
		Nickname:  organization.Nickname(),
		CreatedAt: organization.CreatedAt(),
		UpdatedAt: organization.UpdatedAt(),
	}
}
