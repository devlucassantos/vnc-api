package swagger

import "github.com/google/uuid"

type Party struct {
	Id               uuid.UUID `json:"id"                example:"9bb4028f-6fa8-493a-9fe8-e3bbd341c794"`
	Name             string    `json:"name"              example:"Partido Você na Câmara"`
	Acronym          string    `json:"acronym"           example:"PVNC"`
	ImageUrl         string    `json:"image_url"         example:"https://www.camara.leg.br/internet/Deputado/img/partidos/VNC.gif"`
	ImageDescription string    `json:"image_description" example:"Logo do Partido Você na Câmara (PVNC)"`
}
