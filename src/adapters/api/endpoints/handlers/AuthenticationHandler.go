package handlers

import (
	"github.com/devlucassantos/vnc-domains/src/domains/role"
	"github.com/devlucassantos/vnc-domains/src/domains/user"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"net/http"
	"strings"
	"vnc-api/adapters/api/endpoints/dto/request"
	"vnc-api/adapters/api/endpoints/dto/response"
	"vnc-api/adapters/api/endpoints/handlers/utils"
	"vnc-api/core/interfaces/services"
)

type Authentication struct {
	service services.Authentication
}

func NewAuthenticationHandler(service services.Authentication) *Authentication {
	return &Authentication{service: service}
}

// SignUp
// @ID          SignUp
// @Summary     Sign Up
// @Tags        Authentication
// @Description This request is responsible for signing the user up to the platform.
// @Accept      json
// @Produce     json
// @Param       requestBody body request.SignUp true "Request body"
// @Success 201 {object} swagger.User      "Successful request"
// @Failure 400 {object} swagger.HttpError "Badly formatted request"
// @Failure 401 {object} swagger.HttpError "Unauthorized access"
// @Failure 409 {object} swagger.HttpError "Some of the data provided is conflicting"
// @Failure 422 {object} swagger.HttpError "Some of the data provided is invalid"
// @Failure 500 {object} swagger.HttpError "An unexpected error occurred while processing the request"
// @Failure 503 {object} swagger.HttpError "Some of the services/resources are temporarily unavailable"
// @Router /auth/sign-up [POST]
func (instance Authentication) SignUp(context echo.Context) error {
	var signUpDto request.SignUp
	err := context.Bind(&signUpDto)
	if err != nil {
		log.Warn("Error assigning data from user account creation request to DTO: ", err.Error())
		return context.JSON(http.StatusBadRequest, response.NewBadRequestError())
	}

	var roles []role.Role
	userRole, err := role.NewBuilder().Code(role.InactiveUserRoleCode).Build()
	if err != nil {
		log.Errorf("Error setting account creation role for user %s: %s", signUpDto.Email, err.Error())
		return context.JSON(http.StatusInternalServerError, response.NewInternalServerError())
	}
	roles = append(roles, *userRole)

	userData, err := user.NewBuilder().
		FirstName(signUpDto.FirstName).
		LastName(signUpDto.LastName).
		Email(signUpDto.Email).
		Password(signUpDto.Password).
		Roles(roles).
		Build()
	if err != nil {
		log.Warnf("Error validating data for user %s: %s", signUpDto.Email, err.Error())
		return context.JSON(http.StatusUnprocessableEntity, response.NewHttpError(http.StatusUnprocessableEntity, err.Error()))
	}

	createdUserData, err := instance.service.SignUp(*userData)
	if err != nil {
		if strings.Contains(err.Error(), "user_email_key") {
			log.Warnf("Error registering a new user with an already registered email address(%s): %s",
				signUpDto.Email, err.Error())
			return context.JSON(http.StatusConflict, response.NewHttpError(http.StatusConflict,
				"The email address provided already belongs to an account registered on the platform"))
		} else if strings.Contains(err.Error(), "connection refused") {
			log.Error("Database unavailable: ", err.Error())
			return context.JSON(http.StatusServiceUnavailable, response.NewServiceUnavailableError())
		}

		log.Errorf("Error creating account for user %s: %s", userData.Email(), err.Error())
		return context.JSON(http.StatusInternalServerError, response.NewInternalServerError())
	}

	return context.JSON(http.StatusCreated, response.NewUser(*createdUserData))
}

// SignIn
// @ID          SignIn
// @Summary     Sign In
// @Tags        Authentication
// @Description This request is responsible for signing the user into the platform.
// @Accept      json
// @Produce     json
// @Param       requestBody body request.SignIn true "Request body"
// @Success 200 {object} swagger.User      "Successful request"
// @Failure 400 {object} swagger.HttpError "Badly formatted request"
// @Failure 401 {object} swagger.HttpError "Unauthorized access"
// @Failure 422 {object} swagger.HttpError "Some of the data provided is invalid"
// @Failure 500 {object} swagger.HttpError "An unexpected error occurred while processing the request"
// @Failure 503 {object} swagger.HttpError "Some of the services/resources are temporarily unavailable"
// @Router /auth/sign-in [POST]
func (instance Authentication) SignIn(context echo.Context) error {
	var signInDto request.SignIn
	err := context.Bind(&signInDto)
	if err != nil {
		log.Warn("Error assigning data from login request to DTO: ", err.Error())
		return context.JSON(http.StatusBadRequest, response.NewBadRequestError())
	}

	userData, err := user.NewBuilder().
		Email(signInDto.Email).
		Password(signInDto.Password).
		Build()
	if err != nil {
		log.Warnf("Error validating data for user %s: %s", signInDto.Email, err.Error())
		return context.JSON(http.StatusUnprocessableEntity, response.NewHttpError(http.StatusUnprocessableEntity, err.Error()))
	}

	savedUserData, err := instance.service.SignIn(*userData)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			log.Warn("Could not find data for user ", signInDto.Email)
			return context.JSON(http.StatusUnauthorized, response.NewHttpError(http.StatusUnauthorized,
				"The password is incorrect or the email address is not registered on the platform"))
		} else if strings.Contains(err.Error(), "incorrect password") {
			log.Warnf("The password provided for user %s is incorrect", signInDto.Email)
			return context.JSON(http.StatusUnauthorized, response.NewHttpError(http.StatusUnauthorized,
				"The password is incorrect or the email address is not registered on the platform"))
		} else if strings.Contains(err.Error(), "connection refused") {
			log.Error("Database unavailable: ", err.Error())
			return context.JSON(http.StatusServiceUnavailable, response.NewServiceUnavailableError())
		}

		log.Errorf("Error retrieving access data for user %s: %s", userData.Email(), err.Error())
		return context.JSON(http.StatusInternalServerError, response.NewInternalServerError())
	}

	return context.JSON(http.StatusOK, response.NewUser(*savedUserData))
}

// SignOut
// @ID          SignOut
// @Summary     Sign Out
// @Tags        Authentication
// @Description This request is responsible for signing the user out of the platform.
// @Security    BearerAuth
// @Produce     json
// @Success 204 {object} nil               "Successful request"
// @Failure 400 {object} swagger.HttpError "Badly formatted request"
// @Failure 401 {object} swagger.HttpError "Unauthorized access"
// @Failure 403 {object} swagger.HttpError "Access denied"
// @Failure 500 {object} swagger.HttpError "An unexpected error occurred while processing the request"
// @Failure 503 {object} swagger.HttpError "Some of the services/resources are temporarily unavailable"
// @Router /auth/sign-out [POST]
func (instance Authentication) SignOut(context echo.Context) error {
	_, accessToken := utils.ExtractToken(context.Request().Header.Get("Authorization"))

	claims, httpError := utils.ExtractTokenClaims(accessToken)
	if httpError != nil {
		log.Warn("Error extracting claims from refresh token: ", httpError.Message)
		return context.JSON(httpError.Code, httpError)
	}

	parameter, parameterDescription := "sub", "ID do usuário"
	userId, httpError := utils.ConvertFromStringToUuid(claims.Subject, parameter, parameterDescription)
	if httpError != nil {
		log.Warn("Error converting user ID: ", httpError.Message)
		return context.JSON(http.StatusBadRequest, response.NewBadRequestError())
	}

	parameter, parameterDescription = "session_id", "ID da sessão do usuário"
	sessionId, httpError := utils.ConvertFromStringToUuid(claims.SessionId, parameter, parameterDescription)
	if httpError != nil {
		log.Warn("Error converting session ID for user %s: %s", userId, httpError.Message)
		return context.JSON(http.StatusBadRequest, response.NewBadRequestError())
	}

	err := instance.service.SignOut(userId, sessionId)
	if err != nil {
		if strings.Contains(err.Error(), "connection refused") {
			log.Error("Database unavailable: ", err.Error())
			return context.JSON(http.StatusServiceUnavailable, response.NewServiceUnavailableError())
		}

		log.Errorf("Error signing user %s out of the platform: %s", userId, err.Error())
		return context.JSON(http.StatusInternalServerError, response.NewInternalServerError())
	}

	return context.NoContent(http.StatusNoContent)
}

// Refresh
// @ID          Refresh
// @Summary     Refresh access tokens
// @Tags        Authentication
// @Description This request is responsible for refreshing the user's access tokens on the platform.
// @Accept      json
// @Produce     json
// @Param       requestBody body request.RefreshTokens true "Request body"
// @Success 200 {object} swagger.User      "Successful request"
// @Failure 400 {object} swagger.HttpError "Badly formatted request"
// @Failure 401 {object} swagger.HttpError "Unauthorized access"
// @Failure 422 {object} swagger.HttpError "Some of the data provided is invalid"
// @Failure 500 {object} swagger.HttpError "An unexpected error occurred while processing the request"
// @Failure 503 {object} swagger.HttpError "Some of the services/resources are temporarily unavailable"
// @Router /auth/refresh [POST]
func (instance Authentication) Refresh(context echo.Context) error {
	var refreshTokensDto request.RefreshTokens
	err := context.Bind(&refreshTokensDto)
	if err != nil {
		log.Warn("Error assigning data from refresh tokens request to DTO: ", err.Error())
		return context.JSON(http.StatusBadRequest, response.NewBadRequestError())
	}

	httpError := utils.ValidateRefreshToken(refreshTokensDto.RefreshToken)
	if httpError != nil {
		log.Warn("Error validating refresh token: ", httpError.Message)
		return context.JSON(httpError.Code, httpError)
	}

	claims, httpError := utils.ExtractTokenClaims(refreshTokensDto.RefreshToken)
	if httpError != nil {
		log.Warn("Error extracting claims from refresh token: ", httpError.Message)
		return context.JSON(httpError.Code, httpError)
	}

	parameter, parameterDescription := "sub", "ID do usuário"
	userId, httpError := utils.ConvertFromStringToUuid(claims.Subject, parameter, parameterDescription)
	if httpError != nil {
		log.Warn("Error converting user ID: ", httpError.Message)
		return context.JSON(http.StatusBadRequest, response.NewBadRequestError())
	}

	parameter, parameterDescription = "session_id", "ID da sessão do usuário"
	sessionId, httpError := utils.ConvertFromStringToUuid(claims.SessionId, parameter, parameterDescription)
	if httpError != nil {
		log.Warnf("Error converting session ID for user %s: %s", userId, httpError.Message)
		return context.JSON(http.StatusBadRequest, response.NewBadRequestError())
	}

	savedUserData, err := instance.service.RefreshTokens(userId, sessionId, refreshTokensDto.RefreshToken)
	if err != nil {
		if strings.Contains(err.Error(), "connection refused") {
			log.Error("Database unavailable: ", err.Error())
			return context.JSON(http.StatusServiceUnavailable, response.NewServiceUnavailableError())
		}

		log.Errorf("Error refreshing tokens for user %s: %s", userId, err.Error())
		return context.JSON(http.StatusUnauthorized, response.NewUnauthorizedError())
	}

	return context.JSON(http.StatusOK, response.NewUser(*savedUserData))
}
