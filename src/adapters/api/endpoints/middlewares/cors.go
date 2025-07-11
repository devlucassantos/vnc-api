package middlewares

import (
	"github.com/labstack/echo/v4"
	"strings"
	"vnc-api/adapters/api/utils"
)

func VerifyOrigin(origin string) (bool, error) {
	allowedOrigins := strings.Split(utils.GetenvWithDefaultValue("SERVER_ALLOWED_HOSTS", "*"), ",")
	for _, allowedOrigin := range allowedOrigins {
		if allowedOrigin == "*" || origin == allowedOrigin {
			return true, nil
		}
	}

	return false, &echo.HTTPError{Code: 401, Message: "Unauthorized access"}
}

func OriginInspectSkipper(context echo.Context) bool {
	return false
}
