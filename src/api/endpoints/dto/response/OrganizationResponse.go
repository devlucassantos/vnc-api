package response

import (
	"github.com/devlucassantos/vnc-domains/src/domains/organization"
	"github.com/google/uuid"
	"time"
)

type Organization struct {
	Id        uuid.UUID `json:"id,omitempty"`
	Code      int       `json:"code,omitempty"`
	Name      string    `json:"name,omitempty"`
	Acronym   string    `json:"acronym,omitempty"`
	Type      string    `json:"type,omitempty"`
	Nickname  string    `json:"nickname,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

func NewOrganization(organization organization.Organization) *Organization {
	return &Organization{
		Id:        organization.Id(),
		Code:      organization.Code(),
		Name:      organization.Name(),
		Acronym:   organization.Acronym(),
		Type:      organization.Type(),
		Nickname:  organization.Nickname(),
		CreatedAt: organization.CreatedAt(),
		UpdatedAt: organization.UpdatedAt(),
	}
}
