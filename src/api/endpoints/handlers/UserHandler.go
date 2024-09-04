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
// @Summary     Reenviar email de ativação da conta do usuário
// @Tags        Usuário
// @Description Esta requisição é responsável por reenviar o email de ativação da conta do usuário. A cada envio, um novo código é gerado, invalidando o código anterior.
// @Security    BearerAuth
// @Produce     json
// @Success 204 {object} nil                       "Requisição realizada com sucesso."
// @Failure 400 {object} response.SwaggerHttpError "Requisição mal formulada."
// @Failure 401 {object} response.SwaggerHttpError "Acesso não autorizado."
// @Failure 403 {object} response.SwaggerHttpError "Acesso negado."
// @Failure 409 {object} response.SwaggerHttpError "A conta do usuário já está ativa."
// @Failure 500 {object} response.SwaggerHttpError "Ocorreu um erro inesperado durante o processamento da requisição."
// @Failure 503 {object} response.SwaggerHttpError "Algum dos serviços/recursos está temporariamente indisponível."
// @Router /user/resend-activation-email [PATCH]
func (instance User) ResendActivationEmail(context echo.Context) error {
	userId := utils.GetUserIdFromAuthorizationHeader(context)

	err := instance.service.ResendUserAccountActivationEmail(userId)
	if err != nil {
		if strings.Contains(err.Error(), "connection refused") {
			log.Error("Banco de dados indisponível: ", err.Error())
			return context.JSON(http.StatusServiceUnavailable, response.NewServiceUnavailableError())
		} else if strings.Contains(err.Error(), "conta ativa") {
			log.Errorf("A conta do usuário %s já está ativa: %s", userId, err.Error())
			return context.JSON(http.StatusConflict, response.NewHttpError(http.StatusConflict,
				"Conta ativa, não foi possível reenviar o email."))
		}

		log.Errorf("Erro ao reenviar o email de ativação da conta do usuário %s: %s", userId, err.Error())
		return context.JSON(http.StatusInternalServerError, response.NewInternalServerError())
	}

	return context.NoContent(http.StatusNoContent)
}

// ActivateAccount
// @ID          ActivateAccount
// @Summary     Ativar conta do usuário
// @Tags        Usuário
// @Description Esta requisição é responsável por ativar a conta do usuário, comprovando que o endereço de email fornecido no cadastro realmente existe e pertence ao usuário.
// @Security    BearerAuth
// @Accept      json
// @Produce     json
// @Param       body body request.UserAccountActivation true "JSON com todos os dados necessários para que a ativação da conta seja realizada."
// @Success 200 {array}  response.SwaggerUser      "Requisição realizada com sucesso."
// @Failure 400 {object} response.SwaggerHttpError "Requisição mal formulada ou o código de ativação informado é inválido."
// @Failure 401 {object} response.SwaggerHttpError "Acesso não autorizado."
// @Failure 403 {object} response.SwaggerHttpError "Acesso negado."
// @Failure 409 {object} response.SwaggerHttpError "A conta do usuário já está ativa."
// @Failure 500 {object} response.SwaggerHttpError "Ocorreu um erro inesperado durante o processamento da requisição."
// @Failure 503 {object} response.SwaggerHttpError "Algum dos serviços/recursos está temporariamente indisponível."
// @Router /user/activate-account [PATCH]
func (instance User) ActivateAccount(context echo.Context) error {
	userId := utils.GetUserIdFromAuthorizationHeader(context)

	var UserAccountActivationDto request.UserAccountActivation
	err := context.Bind(&UserAccountActivationDto)
	if err != nil {
		log.Error("Erro ao atribuir os dados da requisição de ativação da conta ao DTO: ", err.Error())
		return context.JSON(http.StatusBadRequest, response.NewBadRequestError())
	}

	userData, err := user.NewBuilder().Id(userId).ActivationCode(UserAccountActivationDto.ActivationCode).Build()
	if err != nil {
		log.Errorf("Erro ao validar os dados do usuário %s: %s", userId, err.Error())
		return context.JSON(http.StatusUnprocessableEntity, response.NewHttpError(http.StatusUnprocessableEntity, err.Error()))
	}

	userData, err = instance.service.ActivateUserAccount(*userData)
	if err != nil {
		if strings.Contains(err.Error(), "connection refused") {
			log.Error("Banco de dados indisponível: ", err.Error())
			return context.JSON(http.StatusServiceUnavailable, response.NewServiceUnavailableError())
		} else if strings.Contains(err.Error(), "conta ativa") {
			log.Errorf("A conta do usuário %s já está ativa: %s", userId, err.Error())
			return context.JSON(http.StatusConflict, response.NewHttpError(http.StatusConflict,
				"Conta ativa, não foi possível prosseguir."))
		} else if strings.Contains(err.Error(), "código de ativação inválido") {
			log.Errorf("O código de ativação informado para a conta do usuário %s é inválido: %s", userId, err.Error())
			return context.JSON(http.StatusBadRequest, response.NewHttpError(http.StatusBadRequest,
				"O código de ativação da conta informado é inválido."))
		}

		log.Errorf("Erro ao reenviar o email de ativação da conta do usuário %s: %s", userId, err.Error())
		return context.JSON(http.StatusInternalServerError, response.NewInternalServerError())
	}

	return context.JSON(http.StatusOK, response.NewUser(*userData))
}
