package response

import (
	"github.com/google/uuid"
	"time"
)

type SwaggerError struct {
	Message string `json:"message" example:"Parâmetro inválido: ID da Proposição"`
}

type SwaggerPropositionPagination struct {
	Page         int                  `json:"page"           example:"1"`
	ItensPerPage int                  `json:"itens_per_page" example:"25"`
	Total        int                  `json:"total"          example:"4562"`
	Data         []SwaggerProposition `json:"data"`
}

type SwaggerProposition struct {
	Id              uuid.UUID             `json:"id"                example:"9dc67bd9-674f-4e4d-9536-07485335c362"`
	Code            int                   `json:"code"              example:"9465723"`
	OriginalTextUrl string                `json:"original_text_url" example:"https://www.camara.leg.br/proposicoesWeb/prop_mostrarintegra?codteor=4865485"`
	Title           string                `json:"title"             example:"Requerimento de Votação Nominal-Destaque de Emenda"`
	Summary         string                `json:"summary"           example:"O presente requerimento foi elaborado pelos deputados..."`
	SubmittedAt     time.Time             `json:"submitted_at"      example:"2023-08-09T14:25:00Z"`
	Deputies        []SwaggerDeputy       `json:"deputies"`
	Organizations   []SwaggerOrganization `json:"organizations"`
	Keywords        []SwaggerKeyword      `json:"keywords"`
	CreatedAt       time.Time             `json:"created_at"        example:"2023-08-09T14:55:00Z"`
	UpdatedAt       time.Time             `json:"updated_at"        example:"2023-08-09T14:55:00Z"`
}

type SwaggerDeputy struct {
	Id            uuid.UUID `json:"id"             example:"a4b04454-f426-44d2-843e-1331510b19ad"`
	Code          int       `json:"code"           example:"87624"`
	Cpf           string    `json:"cpf"            example:"12365478955"`
	Name          string    `json:"name"           example:"José da Silva Santos"`
	ElectoralName string    `json:"electoral_name" example:"José do Povo"`
	ImageUrl      string    `json:"image_url"      example:"https://www.camara.leg.br/internet/deputado/bandep/87624.jpg"`
	CreatedAt     time.Time `json:"created_at"     example:"2022-08-07T15:55:00Z"`
	UpdatedAt     time.Time `json:"updated_at"     example:"2022-08-07T15:55:00Z"`
	Party         *Party    `json:"party"`
}

type SwaggerParty struct {
	Id        uuid.UUID `json:"id"         example:"9bb4028f-6fa8-493a-9fe8-e3bbd341c794"`
	Code      int       `json:"code"       example:"78965"`
	Name      string    `json:"name"       example:"Partido Você na Câmara"`
	Acronym   string    `json:"acronym"    example:"PVNC"`
	ImageUrl  string    `json:"image_url"  example:"https://www.camara.leg.br/internet/Deputado/img/partidos/VNC.gif"`
	CreatedAt time.Time `json:"created_at" example:"2021-05-19T18:10:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2021-05-19T18:10:00Z"`
}

type SwaggerOrganization struct {
	Id        uuid.UUID `json:"id"         example:"9d543e6f-20e3-4895-83e4-26b6a976580e"`
	Code      int       `json:"code"       example:"73214"`
	Name      string    `json:"name"       example:"Organização Você na Câmara"`
	Acronym   string    `json:"acronym"    example:"OVNC"`
	Nickname  string    `json:"nickname"   example:"Org da Câmara"`
	CreatedAt time.Time `json:"created_at" example:"2020-12-13T18:37:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2020-12-13T18:37:00Z"`
}

type SwaggerKeyword struct {
	Id        uuid.UUID `json:"id"         example:"ae01ec73-306b-4726-9ae7-4a90594ae994"`
	Keyword   string    `json:"keyword"    example:"Educação"`
	CreatedAt time.Time `json:"created_at" example:"2019-01-02T09:44:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2019-02-06T09:44:00Z"`
}
