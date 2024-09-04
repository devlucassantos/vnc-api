package middlewares

import (
	"github.com/casbin/casbin/v2"
	"github.com/devlucassantos/vnc-domains/src/domains/role"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"net/http"
	"os"
	"strings"
	"vnc-api/api/config/diconteiner"
	"vnc-api/api/endpoints/dto/response"
	"vnc-api/api/endpoints/handlers/utils"
)

func GuardMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	authModel := os.Getenv("SERVER_CASBIN_AUTHORIZATION_MODEL")
	authPolicy := os.Getenv("SERVER_CASBIN_AUTHORIZATION_POLICY")

	enforcer, err := casbin.NewEnforcer(authModel, authPolicy)
	if err != nil {
		log.Error("Erro ao construir o enforcer: ", err)
	}

	authService := diconteiner.GetAuthenticationService()
	return func(context echo.Context) error {
		authHeader := context.Request().Header.Get("Authorization")
		method := context.Request().Method
		path := context.Request().URL.Path
		roles := utils.ExtractUserAuthorizationRoles(authHeader)
		if roles == nil {
			return context.JSON(http.StatusUnauthorized, response.NewUnauthorizedError())
		}

		for index, userRole := range roles {
			hasAccess, err := enforcer.Enforce(userRole, path, method)
			if err != nil {
				log.Errorf("Erro ao verificar a permissão do usuário ao recurso: Método: %s; Path: %s; Roles: %s)",
					method, path, roles)
				return context.JSON(http.StatusUnauthorized, response.NewUnauthorizedError())
			} else if userRole == role.AnonymousRoleCode {
				if hasAccess {
					return next(context)
				}

				return context.JSON(http.StatusUnauthorized, response.NewUnauthorizedError())
			} else if !hasAccess {
				if index+1 != len(roles) {
					continue
				}

				log.Errorf("Usuário não autorizado para acessar o recurso (Método: %s; Path: %s; Roles: %s)",
					method, path, roles)
				return context.JSON(http.StatusForbidden, response.NewForbiddenError())
			}

			_, accessToken := utils.ExtractToken(authHeader)

			claims, httpError := utils.ExtractTokenClaims(accessToken)
			if httpError != nil {
				log.Error("Erro ao extrair as claims do token de acesso: " + httpError.Message)
				return context.JSON(http.StatusUnauthorized, response.NewUnauthorizedError())
			}

			userId, err := uuid.Parse(claims.Subject)
			if err != nil {
				log.Error("Erro ao converter ID do usuário do token de acesso: " + err.Error())
				return context.JSON(http.StatusUnauthorized, response.NewUnauthorizedError())
			}

			sessionId, err := uuid.Parse(claims.SessionId)
			if err != nil {
				log.Error("Erro ao converter ID da sessão do usuário do token de acesso: " + err.Error())
				return context.JSON(http.StatusUnauthorized, response.NewUnauthorizedError())
			}

			exists, err := authService.SessionExists(userId, sessionId, accessToken)
			if err != nil {
				if strings.Contains(err.Error(), "connection refused") {
					log.Error("Banco de dados indisponível: ", err.Error())
					return context.JSON(http.StatusServiceUnavailable, response.NewServiceUnavailableError())
				}

				return context.JSON(http.StatusUnauthorized, response.NewUnauthorizedError())
			} else if !exists {
				return context.JSON(http.StatusUnauthorized, response.NewUnauthorizedError())
			}
		}

		return next(context)
	}
}
