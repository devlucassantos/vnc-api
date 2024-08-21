package utils

import (
	"encoding/base64"
	"encoding/json"
	"github.com/devlucassantos/vnc-domains/src/domains/role"
	"github.com/devlucassantos/vnc-domains/src/domains/user"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"os"
	"strings"
	"vnc-api/api/endpoints/dto/response"
)

func GetUserIdFromAuthorizationHeader(context echo.Context) uuid.UUID {
	_, accessToken := ExtractToken(context.Request().Header.Get("Authorization"))
	if accessToken == "" {
		return uuid.Nil
	}

	claims, httpError := getAuthorizationClaims(context.Request().Header.Get("Authorization"))
	if httpError != nil {
		log.Error("Erro ao extrair token de autorização da requisição: ", httpError.Message)
		return uuid.Nil
	}

	if claims.Subject == "" {
		log.Error("Token inválido: O token informador não possui o ID do usuário")
		return uuid.Nil
	}

	userId, err := uuid.Parse(claims.Subject)
	if err != nil || userId == uuid.Nil {
		log.Error("Token inválido: O ID do usuário do token informado é inválido")
		return uuid.Nil
	}

	return userId
}

func getAuthorizationClaims(authorizationHeader string) (*user.Claims, *response.HttpError) {
	_, token := ExtractToken(authorizationHeader)
	authenticationClaims, err := ExtractTokenClaims(token)
	if err != nil {
		log.Error("Erro ao extrair as claims do token de acesso: " + err.Message)
		return nil, err
	}

	return authenticationClaims, nil
}

func ExtractToken(authorizationHeader string) (tokenType string, accessToken string) {
	authorization := strings.Split(strings.TrimSpace(authorizationHeader), " ")
	if len(authorization) < 2 {
		return "", ""
	}

	return authorization[0], authorization[1]
}

func ExtractTokenClaims(token string) (*user.Claims, *response.HttpError) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		log.Error("Erro ao separar as partes do token de acesso: A quantidade de partes do token é inválida.")
		return nil, response.NewUnauthorizedError()
	}

	payload := parts[1]
	payloadBytes, err := jwt.DecodeSegment(payload)
	if err != nil {
		log.Error("Erro ao decodificar token: " + err.Error())
		return nil, response.NewUnauthorizedError()
	}

	var claims user.Claims
	err = json.Unmarshal(payloadBytes, &claims)
	if err != nil {
		log.Error("Erro ao atribuir os dados do token a entidade de claims: " + err.Error())
		return nil, response.NewUnauthorizedError()
	}

	return &claims, nil
}

func ValidateRefreshToken(refreshToken string) *response.HttpError {
	publicKey := os.Getenv("SERVER_REFRESH_TOKEN_PUBLIC_KEY")
	publicKeyBytes, err := base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		log.Error("Erro durante a decodificação da chave pública do servidor: ", err.Error())
		return response.NewInternalServerError()
	}

	rsaPublicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyBytes)
	if err != nil {
		log.Error("Erro durante a construção da chave pública do servidor: ", err.Error())
		return response.NewInternalServerError()
	}

	splitToken := strings.Split(refreshToken, ".")
	if len(splitToken) != 3 {
		log.Error("Erro ao separar as partes do token de atualização: A quantidade de partes do token é inválida.")
		return response.NewUnauthorizedError()
	}

	err = jwt.SigningMethodRS256.Verify(strings.Join(splitToken[0:2], "."), splitToken[2], rsaPublicKey)
	if err != nil {
		log.Error("O token de atualização não é autêntico.")
		return response.NewUnauthorizedError()
	}

	return nil
}

func ExtractUserAuthorizationRoles(authorizationHeader string) []string {
	tokenType, accessToken := ExtractToken(authorizationHeader)
	if tokenType == "" || accessToken == "" {
		return []string{role.AnonymousRoleCode}
	}

	accessTokenIsValid := authorizationIsValid(tokenType, accessToken)
	if !accessTokenIsValid {
		return nil
	}

	claims, httpError := ExtractTokenClaims(accessToken)
	if httpError != nil {
		log.Error("Erro ao extrair claims do token de acesso: ", httpError.Message)
		return nil
	}

	return claims.Roles
}

func authorizationIsValid(tokenType, accessToken string) bool {
	publicKeyBytes, err := base64.StdEncoding.DecodeString(os.Getenv("SERVER_ACCESS_TOKEN_PUBLIC_KEY"))
	if err != nil {
		log.Error("Erro durante a decodificação da chave pública de acesso: ", err.Error())
		return false
	}

	rsaPublicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyBytes)
	if err != nil {
		log.Error("Erro durante a construção da chave pública de acesso: ", err.Error())
		return false
	}

	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		return rsaPublicKey, nil
	})
	if err != nil {
		log.Error("Erro durante a conversão do token: ", err.Error())
		return false
	}

	if !token.Valid || token.Claims.Valid() != nil {
		log.Error("O token informado é inválido ou está expirado.")
		return false
	}

	if strings.ToLower(tokenType) != "bearer" {
		log.Error("O tipo do token utilizado não é suportado: ", tokenType)
		return false
	}

	return true
}
