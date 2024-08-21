package handlers

import (
	"github.com/devlucassantos/vnc-domains/src/domains/role"
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

type Authentication struct {
	service services.Authentication
}

func NewAuthenticationHandler(service services.Authentication) *Authentication {
	return &Authentication{service: service}
}

// SignUp
// @ID          SignUp
// @Summary     Criar conta
// @Tags        Autenticação
// @Description Esta requisição é responsável por permitir o cadastro do usuário na plataforma:
// @Accept      json
// @Produce     json
// @Param       body body request.SignUp true "JSON com todos os dados necessários para que a criação da conta seja realizada."
// @Success 201 {object} response.SwaggerUser      "Requisição realizada com sucesso."
// @Failure 400 {object} response.SwaggerHttpError "Requisição mal formulada."
// @Failure 409 {object} response.SwaggerHttpError "Requisição contém dados já cadastrados no banco de dados que deveriam ser únicos."
// @Failure 422 {object} response.SwaggerHttpError "Requisição não processada devido a algum dos dados enviados serem inválidos."
// @Failure 500 {object} response.SwaggerHttpError "Ocorreu um erro inesperado durante o processamento da requisição."
// @Failure 503 {object} response.SwaggerHttpError "Algum dos serviços/recursos está temporariamente indisponível."
// @Router /auth/sign-up [POST]
func (instance Authentication) SignUp(context echo.Context) error {
	var signUpDto request.SignUp
	err := context.Bind(&signUpDto)
	if err != nil {
		log.Error("Erro ao atribuir os dados da requisição ao DTO: ", err.Error())
		return context.JSON(http.StatusBadRequest, response.NewBadRequestError())
	}

	var roles []role.Role
	userRole, err := role.NewBuilder().Code("USER").Build()
	if err != nil {
		log.Errorf("Erro ao definir a role do usuário %s: %s", signUpDto.Email, err.Error())
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
		log.Errorf("Erro ao validar os dados do usuário %s: %s", signUpDto.Email, err.Error())
		return context.JSON(http.StatusUnprocessableEntity, response.NewHttpError(http.StatusUnprocessableEntity, err.Error()))
	}

	createdUserData, err := instance.service.SignUp(*userData)
	if err != nil {
		if strings.Contains(err.Error(), "user_email_key") {
			log.Errorf("Erro ao tentar cadastrar um novo usuário com um email já cadastrado (%s): %s",
				signUpDto.Email, err.Error())
			return context.JSON(http.StatusConflict, response.NewHttpError(http.StatusConflict, "O email "+
				"informado já pertence a uma conta cadastrada na plataforma."))
		} else if strings.Contains(err.Error(), "connection refused") {
			log.Error("Banco de dados indisponível: ", err.Error())
			return context.JSON(http.StatusServiceUnavailable, response.NewServiceUnavailableError())
		}

		log.Errorf("Erro ao cadastrar os dados do usuário %s: %s", userData.Email(), err.Error())
		return context.JSON(http.StatusInternalServerError, response.NewInternalServerError())
	}

	return context.JSON(http.StatusCreated, response.NewUser(*createdUserData))
}

// SignIn
// @ID          SignIn
// @Summary     Fazer login
// @Tags        Autenticação
// @Description Esta requisição é responsável por permitir a entrada do usuário em sua conta na plataforma:
// @Accept      json
// @Produce     json
// @Param       body body request.SignIn true "JSON com todos os dados necessários para que o login seja realizado."
// @Success 200 {object} response.SwaggerUser      "Requisição realizada com sucesso."
// @Failure 400 {object} response.SwaggerHttpError "Requisição mal formulada."
// @Failure 401 {object} response.SwaggerHttpError "Acesso não autorizado."
// @Failure 404 {object} response.SwaggerHttpError "Recurso solicitado não encontrado."
// @Failure 422 {object} response.SwaggerHttpError "Requisição não processada devido a algum dos dados enviados serem inválidos."
// @Failure 500 {object} response.SwaggerHttpError "Ocorreu um erro inesperado durante o processamento da requisição."
// @Failure 503 {object} response.SwaggerHttpError "Algum dos serviços/recursos está temporariamente indisponível."
// @Router /auth/sign-in [POST]
func (instance Authentication) SignIn(context echo.Context) error {
	var signInDto request.SignIn
	err := context.Bind(&signInDto)
	if err != nil {
		log.Error("Erro ao atribuir os dados da requisição ao DTO: ", err.Error())
		return context.JSON(http.StatusBadRequest, response.NewBadRequestError())
	}

	userData, err := user.NewBuilder().
		Email(signInDto.Email).
		Password(signInDto.Password).
		Build()
	if err != nil {
		log.Errorf("Erro ao validar os dados do usuário %s: %s", signInDto.Email, err.Error())
		return context.JSON(http.StatusUnprocessableEntity, response.NewHttpError(http.StatusUnprocessableEntity, err.Error()))
	}

	savedUserData, err := instance.service.SignIn(*userData)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			log.Error("Não foi possível encontrar os dados do usuário ", signInDto.Email)
			return context.JSON(http.StatusNotFound, response.NewHttpError(http.StatusNotFound,
				"Usuário não encontrado. Por favor, verifique se o email informado está correto e tente novamente."))
		} else if strings.Contains(err.Error(), "senha incorreta") {
			log.Errorf("A senha informada para o usuário %s é incorreta", signInDto.Email)
			return context.JSON(http.StatusUnauthorized, response.NewHttpError(http.StatusUnauthorized,
				"Senha incorreta. Por favor, verifique a senha informada e tente novamente."))
		} else if strings.Contains(err.Error(), "connection refused") {
			log.Error("Banco de dados indisponível: ", err.Error())
			return context.JSON(http.StatusServiceUnavailable, response.NewServiceUnavailableError())
		}

		log.Errorf("Erro ao buscas dados de acesso do usuário %s: %s", userData.Email(), err.Error())
		return context.JSON(http.StatusInternalServerError, response.NewInternalServerError())
	}

	return context.JSON(http.StatusOK, response.NewUser(*savedUserData))
}

// Refresh
// @ID          Refresh
// @Summary     Atualizar tokens de acesso
// @Tags        Autenticação
// @Description Esta requisição é responsável por realizar a atualização dos tokens do usuário na plataforma:
// @Accept      json
// @Produce     json
// @Param       body body request.RefreshTokens true "JSON com todos os dados necessários para que o login seja realizado."
// @Success 200 {object} response.SwaggerUser      "Requisição realizada com sucesso."
// @Failure 400 {object} response.SwaggerHttpError "Requisição mal formulada."
// @Failure 401 {object} response.SwaggerHttpError "Acesso não autorizado."
// @Failure 422 {object} response.SwaggerHttpError "Requisição não processada devido a algum dos dados enviados serem inválidos."
// @Failure 500 {object} response.SwaggerHttpError "Ocorreu um erro inesperado durante o processamento da requisição."
// @Failure 503 {object} response.SwaggerHttpError "Algum dos serviços/recursos está temporariamente indisponível."
// @Router /auth/refresh [POST]
func (instance Authentication) Refresh(context echo.Context) error {
	var refreshTokensDto request.RefreshTokens
	err := context.Bind(&refreshTokensDto)
	if err != nil {
		log.Error("Erro ao atribuir os dados da requisição ao DTO: ", err.Error())
		return context.JSON(http.StatusBadRequest, response.NewBadRequestError())
	}

	httpError := utils.ValidateRefreshToken(refreshTokensDto.RefreshToken)
	if httpError != nil {
		log.Error("Erro ao validar token de atualização: ", httpError.Message)
		return context.JSON(httpError.Code, httpError)
	}

	claims, httpError := utils.ExtractTokenClaims(refreshTokensDto.RefreshToken)
	if httpError != nil {
		log.Error("Erro ao extrair as claims do token de atualização: ", httpError.Message)
		return context.JSON(httpError.Code, httpError)
	}

	userId, httpError := utils.ConvertFromStringToUuid(claims.Subject, "ID do usuário")
	if httpError != nil {
		log.Error("Erro ao converter ID do usuário: ", httpError.Message)
		return context.JSON(httpError.Code, httpError)
	}

	sessionId, httpError := utils.ConvertFromStringToUuid(claims.SessionId, "ID da sessão do usuário")
	if httpError != nil {
		log.Error("Erro ao converter ID da sessão do usuário: ", httpError.Message)
		return context.JSON(httpError.Code, httpError)
	}

	savedUserData, err := instance.service.RefreshTokens(userId, sessionId, refreshTokensDto.RefreshToken)
	if err != nil {
		if strings.Contains(err.Error(), "connection refused") {
			log.Error("Banco de dados indisponível: ", err.Error())
			return context.JSON(http.StatusServiceUnavailable, response.NewServiceUnavailableError())
		}

		log.Error("Erro ao atualizar tokens do usuário: ", err.Error())
		return context.JSON(http.StatusUnauthorized, response.NewUnauthorizedError())
	}

	return context.JSON(http.StatusOK, response.NewUser(*savedUserData))
}

// SignOut
// @ID          SignOut
// @Summary     Fazer logout
// @Tags        Autenticação
// @Description Esta requisição é responsável por realizar o encerramento do acesso do usuário a plataforma:
// @Security    BearerAuth
// @Produce     json
// @Success 204 {object} nil                       "Requisição realizada com sucesso."
// @Failure 400 {object} response.SwaggerHttpError "Requisição mal formulada."
// @Failure 401 {object} response.SwaggerHttpError "Acesso não autorizado."
// @Failure 403 {object} response.SwaggerHttpError "Acesso negado."
// @Failure 500 {object} response.SwaggerHttpError "Ocorreu um erro inesperado durante o processamento da requisição."
// @Failure 503 {object} response.SwaggerHttpError "Algum dos serviços/recursos está temporariamente indisponível."
// @Router /auth/sign-out [POST]
func (instance Authentication) SignOut(context echo.Context) error {
	_, accessToken := utils.ExtractToken(context.Request().Header.Get("Authorization"))

	claims, httpError := utils.ExtractTokenClaims(accessToken)
	if httpError != nil {
		log.Error("Erro ao extrair as claims do token de atualização: ", httpError.Message)
		return context.JSON(httpError.Code, httpError)
	}

	userId, httpError := utils.ConvertFromStringToUuid(claims.Subject, "ID do usuário")
	if httpError != nil {
		log.Error("Erro ao converter ID do usuário: ", httpError.Message)
		return context.JSON(httpError.Code, httpError)
	}

	sessionId, httpError := utils.ConvertFromStringToUuid(claims.SessionId, "ID da sessão do usuário")
	if httpError != nil {
		log.Error("Erro ao converter ID da sessão do usuário: ", httpError.Message)
		return context.JSON(httpError.Code, httpError)
	}

	err := instance.service.SignOut(userId, sessionId)
	if err != nil {
		if strings.Contains(err.Error(), "connection refused") {
			log.Error("Banco de dados indisponível: ", err.Error())
			return context.JSON(http.StatusServiceUnavailable, response.NewServiceUnavailableError())
		}

		log.Error("Erro ao encerrar acesso do usuário ao sistema: ", err.Error())
		return context.JSON(http.StatusInternalServerError, response.NewInternalServerError())
	}

	return context.NoContent(http.StatusNoContent)
}
