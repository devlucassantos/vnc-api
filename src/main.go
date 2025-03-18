package main

import (
	"vnc-api/api/config"
	_ "vnc-api/docs"
)

// @title         Você na Câmara API
// @version       v1
// @description   Set of endpoints that make up the backend of the Você na Câmara platform, structured to enable communication between the services and the execution of the system's functionalities.
// @contact.name  Você na Câmara
// @contact.email email.vocenacamara@gmail.com
// @basePath      /api/v1
// @securityDefinitions.apikey BearerAuth
// @in   header
// @name Authorization
func main() {
	config.NewApi().Serve()
}
