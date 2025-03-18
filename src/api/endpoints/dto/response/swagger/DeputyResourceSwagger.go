package swagger

import "github.com/google/uuid"

type DeputyResource struct {
	Id               uuid.UUID `json:"id"                example:"a4b04454-f426-44d2-843e-1331510b19ad"`
	Name             string    `json:"name"              example:"José da Silva Santos"`
	ElectoralName    string    `json:"electoral_name"    example:"José do Povo"`
	ImageUrl         string    `json:"image_url"         example:"https://www.camara.leg.br/internet/deputado/bandep/87624.jpg"`
	ImageDescription string    `json:"image_description" example:"Foto do(a) deputado(a) federal José do Povo (PVNC-AL)"`
	Party            Party     `json:"party"`
	FederatedUnit    string    `json:"federated_unit"    example:"AL"`
}
