package main

import (
	"vnc-api/api/config"
	_ "vnc-api/docs"
)

// @title         API da Plataforma Você na Câmara
// @version       v1
// @description   Conjunto de rotas responsável por gerenciar a manipulação de dados da Plataforma Você na Câmara.
// @contact.name  Você na Câmara
// @contact.email email.vocenacamara@gmail.com
// @basePath      /api/v1
// @securityDefinitions.apikey BearerAuth
// @in   header
// @name Authorization
func main() {
	config.NewApi().Serve()
}
