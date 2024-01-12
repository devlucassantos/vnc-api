package response

import (
	"github.com/google/uuid"
	"time"
)

type SwaggerError struct {
	Message string `json:"message" example:"Parâmetro inválido: ID da Proposição"`
}

type SwaggerNewsPagination struct {
	Page         int           `json:"page"           example:"1"`
	ItensPerPage int           `json:"itens_per_page" example:"25"`
	Total        int           `json:"total"          example:"4562"`
	Data         []SwaggerNews `json:"data"`
}

type SwaggerNews struct {
	Id        uuid.UUID `json:"id"         example:"b27947d6-3224-4479-8da4-7917ae16b34d"`
	Title     string    `json:"title"      example:"Novo projeto de lei impulsiona crescimento do setor portuário até 2028"`
	Content   string    `json:"content"    example:"Foi sancionado o projeto de lei que altera a Lei n° 11.033 para prorrogar o Regime Tributário..."`
	Type      string    `json:"type"       example:"Proposição"`
	CreatedAt time.Time `json:"created_at" example:"2024-01-05T20:25:19Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2024-01-05T20:25:19Z"`
}

type SwaggerNewsletter struct {
	Id            uuid.UUID                      `json:"id"                example:"7963a6dd-f0b8-4065-8e56-6d2a79626db7"`
	Title         string                         `json:"title"             example:"Proposta inovadora busca impulsionar o crescimento empresarial"`
	Content       string                         `json:"content"           example:"O presidente enviou ao Congresso Nacional um projeto de lei que permite a concessão de descontos fiscais..."`
	ReferenceDate time.Time                      `json:"reference_date"    example:"2023-12-23T16:34:14Z"`
	Propositions  []SwaggerNewsletterProposition `json:"propositions"`
	CreatedAt     time.Time                      `json:"created_at"        example:"2023-12-24T19:15:22Z"`
	UpdatedAt     time.Time                      `json:"updated_at"        example:"2023-12-24T19:15:22Z"`
}

type SwaggerNewsletterProposition struct {
	Id              uuid.UUID `json:"id"                example:"9dc67bd9-674f-4e4d-9536-07485335c362"`
	Code            int       `json:"code"              example:"9465723"`
	OriginalTextUrl string    `json:"original_text_url" example:"https://www.camara.leg.br/proposicoesWeb/prop_mostrarintegra?codteor=4865485"`
	Title           string    `json:"title"             example:"Requerimento de Votação Nominal-Destaque de Emenda"`
	Content         string    `json:"content"           example:"O presente requerimento foi elaborado pelos deputados..."`
	SubmittedAt     time.Time `json:"submitted_at"      example:"2023-08-09T14:25:00Z"`
	ImageUrl        string    `json:"image_url"         example:"https://www.vnc.com.br/news/proposition/image/87624.jpg"`
	CreatedAt       time.Time `json:"created_at"        example:"2023-08-09T14:55:00Z"`
	UpdatedAt       time.Time `json:"updated_at"        example:"2023-08-09T14:55:00Z"`
}

type SwaggerProposition struct {
	Id              uuid.UUID                    `json:"id"                example:"9dc67bd9-674f-4e4d-9536-07485335c362"`
	Code            int                          `json:"code"              example:"9465723"`
	OriginalTextUrl string                       `json:"original_text_url" example:"https://www.camara.leg.br/proposicoesWeb/prop_mostrarintegra?codteor=4865485"`
	Title           string                       `json:"title"             example:"Requerimento de Votação Nominal-Destaque de Emenda"`
	Content         string                       `json:"content"           example:"O presente requerimento foi elaborado pelos deputados..."`
	SubmittedAt     time.Time                    `json:"submitted_at"      example:"2023-08-09T14:25:00Z"`
	ImageUrl        string                       `json:"image_url"         example:"https://www.vnc.com.br/news/proposition/image/87624.jpg"`
	Deputies        []SwaggerDeputy              `json:"deputies"`
	Organizations   []SwaggerOrganization        `json:"organizations"`
	Newsletter      SwaggerNewsletterProposition `json:"newsletter"`
	CreatedAt       time.Time                    `json:"created_at"        example:"2023-08-09T14:55:00Z"`
	UpdatedAt       time.Time                    `json:"updated_at"        example:"2023-08-09T14:55:00Z"`
}

type SwaggerDeputy struct {
	Id                    uuid.UUID     `json:"id"             example:"a4b04454-f426-44d2-843e-1331510b19ad"`
	Code                  int           `json:"code"           example:"87624"`
	Cpf                   string        `json:"cpf"            example:"12365478955"`
	Name                  string        `json:"name"           example:"José da Silva Santos"`
	ElectoralName         string        `json:"electoral_name" example:"José do Povo"`
	ImageUrl              string        `json:"image_url"      example:"https://www.camara.leg.br/internet/deputado/bandep/87624.jpg"`
	CreatedAt             time.Time     `json:"created_at"     example:"2022-08-07T15:55:00Z"`
	UpdatedAt             time.Time     `json:"updated_at"     example:"2022-08-07T15:55:00Z"`
	Party                 *SwaggerParty `json:"party"`
	PartyInTheProposition *SwaggerParty `json:"party_in_the_proposal"`
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
	Nickname  string    `json:"nickname"   example:"VNC"`
	Type      string    `json:"type"       example:"Org da Câmara"`
	CreatedAt time.Time `json:"created_at" example:"2020-12-13T18:37:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2020-12-13T18:37:00Z"`
}

type SwaggerResources struct {
	Parties       []SwaggerParty          `json:"parties"`
	Deputies      []SwaggerDeputyResource `json:"deputies"`
	Organizations []SwaggerOrganization   `json:"organizations"`
}

type SwaggerDeputyResource struct {
	Id            uuid.UUID     `json:"id"             example:"a4b04454-f426-44d2-843e-1331510b19ad"`
	Code          int           `json:"code"           example:"87624"`
	Cpf           string        `json:"cpf"            example:"12365478955"`
	Name          string        `json:"name"           example:"José da Silva Santos"`
	ElectoralName string        `json:"electoral_name" example:"José do Povo"`
	ImageUrl      string        `json:"image_url"      example:"https://www.camara.leg.br/internet/deputado/bandep/87624.jpg"`
	CreatedAt     time.Time     `json:"created_at"     example:"2022-08-07T15:55:00Z"`
	UpdatedAt     time.Time     `json:"updated_at"     example:"2022-08-07T15:55:00Z"`
	Party         *SwaggerParty `json:"party"`
}
