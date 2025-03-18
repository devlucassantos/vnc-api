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
		log.Error("Error extracting authorization token from request: ", httpError.Message)
		return uuid.Nil
	}

	if claims.Subject == "" {
		log.Warn("Invalid token: The provided token does not contain the user ID")
		return uuid.Nil
	}

	userId, err := uuid.Parse(claims.Subject)
	if err != nil || userId == uuid.Nil {
		log.Warn("Invalid token: The user ID contained in the provided token is invalid")
		return uuid.Nil
	}

	return userId
}

func getAuthorizationClaims(authorizationHeader string) (*user.Claims, *response.HttpError) {
	_, token := ExtractToken(authorizationHeader)
	authenticationClaims, err := ExtractTokenClaims(token)
	if err != nil {
		log.Warn("Error extracting claims from access token: ", err.Message)
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
		log.Warn("Error separating the access token parts: The number of token parts is invalid")
		return nil, response.NewUnauthorizedError()
	}

	payload := parts[1]
	payloadBytes, err := jwt.DecodeSegment(payload)
	if err != nil {
		log.Warn("Error decoding token: ", err.Error())
		return nil, response.NewUnauthorizedError()
	}

	var claims user.Claims
	err = json.Unmarshal(payloadBytes, &claims)
	if err != nil {
		log.Warn("Error assigning token data to user claims: ", err.Error())
		return nil, response.NewUnauthorizedError()
	}

	return &claims, nil
}

func ValidateRefreshToken(refreshToken string) *response.HttpError {
	publicKey := os.Getenv("SERVER_REFRESH_TOKEN_PUBLIC_KEY")
	publicKeyBytes, err := base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		log.Error("Error decoding server refresh token public key: ", err.Error())
		return response.NewInternalServerError()
	}

	rsaPublicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyBytes)
	if err != nil {
		log.Error("Error building the server refresh token public key: ", err.Error())
		return response.NewInternalServerError()
	}

	splitToken := strings.Split(refreshToken, ".")
	if len(splitToken) != 3 {
		log.Warn("Error splitting the refresh token parts: The number of token parts is invalid")
		return response.NewUnauthorizedError()
	}

	err = jwt.SigningMethodRS256.Verify(strings.Join(splitToken[0:2], "."), splitToken[2], rsaPublicKey)
	if err != nil {
		log.Warn("The refresh token provided is not authentic")
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
		log.Warn("Error extracting claims from access token: ", httpError.Message)
		return nil
	}

	return claims.Roles
}

func authorizationIsValid(tokenType, accessToken string) bool {
	publicKeyBytes, err := base64.StdEncoding.DecodeString(os.Getenv("SERVER_ACCESS_TOKEN_PUBLIC_KEY"))
	if err != nil {
		log.Error("Error decoding the server access token public key: ", err.Error())
		return false
	}

	rsaPublicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyBytes)
	if err != nil {
		log.Error("Error building the server access token public key: ", err.Error())
		return false
	}

	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		return rsaPublicKey, nil
	})
	if err != nil {
		log.Warn("Error converting the provided token: ", err.Error())
		return false
	}

	if !token.Valid || token.Claims.Valid() != nil {
		log.Warn("The token provided is invalid or expired")
		return false
	}

	if strings.ToLower(tokenType) != "bearer" {
		log.Warn("The token type used is not supported: ", tokenType)
		return false
	}

	return true
}
