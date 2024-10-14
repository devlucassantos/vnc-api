package main

import (
	"vnc-api/api/config"
	_ "vnc-api/docs"
)

// @title         Você na Câmara API
// @version       v1
// @description   Set of routes responsible for managing data manipulation in Você na Câmara applications.
// @contact.name  Você na Câmara
// @contact.email email.vocenacamara@gmail.com
// @basePath      /api/v1
// @securityDefinitions.apikey BearerAuth
// @in   header
// @name Authorization
func main() {
	config.NewApi().Serve()
}
