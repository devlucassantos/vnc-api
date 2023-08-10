package main

import (
	"vnc-read-api/api/config"
	_ "vnc-read-api/docs"
)

// @Title       VNC Read API
// @Version     v1
// @Description Este repositório é responsável pela leitura dos dados nas bases de dados da Plataforma Você na Câmara.
// @BasePath    /api/v1
func main() {
	config.NewServer()
}
