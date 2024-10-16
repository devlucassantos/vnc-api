package handlers

import (
	"github.com/devlucassantos/vnc-domains/src/domains/user"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"net/http"
	"strings"
	"vnc-api/api/endpoints/dto/request"
	"vnc-api/api/endpoints/dto/response"
	"vnc-api/api/endpoints/handlers/utils"
	"vnc-api/core/interfaces/services"
)

type User struct {
	service services.User
}

func NewUserHandler(service services.User) *User {
	return &User{service: service}
}

// ResendActivationEmail
// @ID          ResendActivationEmail
// @Summary     Resend user account activation email
// @Tags        Users
// @Description This request is responsible for resending the user account activation email. Each time it is sent, a new code is generated, invalidating the previous code.
// @Security    BearerAuth
// @Produce     json
// @Success 204 {object} nil                       "Successful request"
// @Failure 400 {object} response.SwaggerHttpError "Badly formatted request"
// @Failure 401 {object} response.SwaggerHttpError "Unauthorized access"
// @Failure 403 {object} response.SwaggerHttpError "Access denied"
// @Failure 409 {object} response.SwaggerHttpError "Some of the data provided is conflicting"
// @Failure 500 {object} response.SwaggerHttpError "An unexpected error occurred while processing the request"
// @Failure 503 {object} response.SwaggerHttpError "Some of the services/resources are temporarily unavailable"
// @Router /user/resend-activation-email [PATCH]
func (instance User) ResendActivationEmail(context echo.Context) error {
	userId := utils.GetUserIdFromAuthorizationHeader(context)

	err := instance.service.ResendUserAccountActivationEmail(userId)
	if err != nil {
		if strings.Contains(err.Error(), "connection refused") {
			log.Error("Database unavailable: ", err.Error())
			return context.JSON(http.StatusServiceUnavailable, response.NewServiceUnavailableError())
		} else if strings.Contains(err.Error(), "active account") {
			log.Errorf("The account for user %s is already active: %s", userId, err.Error())
			return context.JSON(http.StatusConflict, response.NewHttpError(http.StatusConflict,
				"Account active, email could not be resent"))
		}

		log.Errorf("Error resending the account activation email for user %s: %s", userId, err.Error())
		return context.JSON(http.StatusInternalServerError, response.NewInternalServerError())
	}

	return context.NoContent(http.StatusNoContent)
}

// ActivateAccount
// @ID          ActivateAccount
// @Summary     Activate user account
// @Tags        Users
// @Description This request is responsible for activating the user's account, proving that the email address provided during registration really exists and belongs to the user.
// @Security    BearerAuth
// @Accept      json
// @Produce     json
// @Param       body body request.UserAccountActivation true "Request body"
// @Success 200 {array}  response.SwaggerUser      "Successful request"
// @Failure 400 {object} response.SwaggerHttpError "Badly formatted request"
// @Failure 401 {object} response.SwaggerHttpError "Unauthorized access"
// @Failure 403 {object} response.SwaggerHttpError "Access denied"
// @Failure 409 {object} response.SwaggerHttpError "Some of the data provided is conflicting"
// @Failure 500 {object} response.SwaggerHttpError "An unexpected error occurred while processing the request"
// @Failure 503 {object} response.SwaggerHttpError "Some of the services/resources are temporarily unavailable"
// @Router /user/activate-account [PATCH]
func (instance User) ActivateAccount(context echo.Context) error {
	userId := utils.GetUserIdFromAuthorizationHeader(context)

	var UserAccountActivationDto request.UserAccountActivation
	err := context.Bind(&UserAccountActivationDto)
	if err != nil {
		log.Error("Error assigning data from account activation request to DTO: ", err.Error())
		return context.JSON(http.StatusBadRequest, response.NewBadRequestError())
	}

	userData, err := user.NewBuilder().Id(userId).ActivationCode(UserAccountActivationDto.ActivationCode).Build()
	if err != nil {
		log.Errorf("Error validating data for user %s: %s", userId, err.Error())
		return context.JSON(http.StatusUnprocessableEntity, response.NewHttpError(http.StatusUnprocessableEntity, err.Error()))
	}

	userData, err = instance.service.ActivateUserAccount(*userData)
	if err != nil {
		if strings.Contains(err.Error(), "connection refused") {
			log.Error("Database unavailable: ", err.Error())
			return context.JSON(http.StatusServiceUnavailable, response.NewServiceUnavailableError())
		} else if strings.Contains(err.Error(), "active account") {
			log.Errorf("The account for user %s is already active: %s", userId, err.Error())
			return context.JSON(http.StatusConflict, response.NewHttpError(http.StatusConflict,
				"Active account, could not proceed"))
		} else if strings.Contains(err.Error(), "invalid activation code") {
			log.Errorf("The activation code provided for the user account %s is invalid: %s", userId, err.Error())
			return context.JSON(http.StatusBadRequest, response.NewHttpError(http.StatusBadRequest,
				"The account activation code provided is invalid"))
		}

		log.Errorf("Error activating account for user %s: %s", userId, err.Error())
		return context.JSON(http.StatusInternalServerError, response.NewInternalServerError())
	}

	return context.JSON(http.StatusOK, response.NewUser(*userData))
}
