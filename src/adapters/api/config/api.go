package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"net/http"
	"os"
	"vnc-api/adapters/api/endpoints/middlewares"
	"vnc-api/adapters/api/endpoints/router"
	"vnc-api/adapters/api/utils"
)

type Api interface {
	Serve()
	loadRoutes()
}

type Options struct{}

type api struct {
	group        *echo.Group
	echoInstance *echo.Echo
}

func NewApi() Api {
	serverMode := os.Getenv("SERVER_MODE")
	if serverMode != "production" && serverMode != "stage" {
		err := godotenv.Load("config/.env")
		if err != nil {
			log.Fatal("Environment variables file not found: ", err.Error())
		}
	}

	echoInstance := echo.New()
	return &api{echoInstance.Group("/api"), echoInstance}
}

func (instance *api) Serve() {
	instance.echoInstance.Use(middleware.Logger())
	instance.echoInstance.Use(middleware.Recover())
	instance.echoInstance.Use(instance.getCORSSettings())
	instance.echoInstance.Use(middlewares.GuardMiddleware)
	instance.loadRoutes()
	address := getServerAddress()
	instance.echoInstance.Logger.Fatal(instance.echoInstance.Start(address))
}

func (instance *api) loadRoutes() {
	router.New().Load(instance.group)
}

func (instance *api) getCORSSettings() echo.MiddlewareFunc {
	return middleware.CORSWithConfig(middleware.CORSConfig{
		Skipper:         middlewares.OriginInspectSkipper,
		AllowOriginFunc: middlewares.VerifyOrigin,
		AllowMethods: []string{
			http.MethodHead,
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
			http.MethodPatch,
		},
	})
}

func getServerAddress() string {
	host := utils.GetenvWithDefaultValue("HOST", "0.0.0.0")
	port := utils.GetenvWithDefaultValue("PORT", "8080")
	return fmt.Sprintf("%s:%s", host, port)
}
