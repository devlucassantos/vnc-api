package services

import (
	"errors"
	"github.com/devlucassantos/vnc-domains/src/domains/role"
	"github.com/devlucassantos/vnc-domains/src/domains/user"
	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	"strings"
	"vnc-api/core/interfaces/repositories"
	"vnc-api/core/interfaces/services"
	"vnc-api/core/services/utils"
)

type User struct {
	userRepository    repositories.User
	sessionRepository repositories.Session
	emailService      services.Email
}

func NewUserService(userRepository repositories.User, sessionRepository repositories.Session,
	emailService services.Email) *User {
	return &User{
		userRepository:    userRepository,
		sessionRepository: sessionRepository,
		emailService:      emailService,
	}
}

func (instance User) ResendUserAccountActivationEmail(userId uuid.UUID) error {
	userData, err := instance.userRepository.GetUserById(userId)
	if err != nil {
		log.Errorf("Erro ao obter os dados do usuário %s no banco de dados: %s", userId, err.Error())
		return err
	}

	if len(userData.Roles()) != 1 || userData.Roles()[0].Code() != role.InactiveUserRoleCode {
		errorMessage := "conta ativa"
		log.Errorf("Erro ao reenviar email de ativação da conta do usuário %s: %s", userId, errorMessage)
		return errors.New(errorMessage)
	}

	userActivationCode, err := utils.GenerateUserActivationCode()
	if err != nil {
		log.Errorf("Erro ao gerar o novo código de ativação da conta do usuário %s: %s",
			userData.Id(), err.Error())
		return err
	}

	userData, err = userData.NewUpdater().ActivationCode(userActivationCode).Build()
	if err != nil {
		log.Errorf("Erro ao definir o novo código de ativação da conta do usuário %s: %s", userData.Id(),
			err.Error())
		return nil
	}

	userData, err = instance.userRepository.UpdateUser(*userData)
	if err != nil {
		log.Errorf("Erro ao atualizar código de ativação da conta do usuário %s no banco de dados: %s",
			userData.Id(), err.Error())
		return err
	}

	err = instance.emailService.SendUserAccountActivationEmail(*userData)
	if err != nil {
		log.Error("Erro ao reenviar email de ativação da conta do usuário %s: %s", userData.Id(), err)
	}

	return nil
}

func (instance User) ActivateUserAccount(userAccountActivationData user.User) (*user.User, error) {
	userData, err := instance.userRepository.GetUserById(userAccountActivationData.Id())
	if err != nil {
		log.Errorf("Erro ao obter os dados do usuário %s no banco de dados: %s",
			userAccountActivationData.Id(), err.Error())
		return nil, err
	}

	if len(userData.Roles()) != 1 || userData.Roles()[0].Code() != role.InactiveUserRoleCode {
		errorMessage := "conta ativa"
		log.Errorf("Erro ao ativar conta do usuário %s: ",
			userAccountActivationData.Id(), errorMessage)
		return nil, errors.New(errorMessage)
	}

	if userData.ActivationCode() != strings.ToUpper(userAccountActivationData.ActivationCode()) {
		errorMessage := "código de ativação inválido"
		log.Errorf("O código de ativação informado para a conta do usuário %s é inválido: %s",
			userData.Id(), errorMessage)
		return nil, errors.New(errorMessage)
	}

	var roles []role.Role
	userRole, err := role.NewBuilder().Code(role.UserRoleCode).Build()
	if err != nil {
		log.Errorf("Erro ao definir a role de ativação da conta do usuário %s: %s", userData.Id(), err.Error())
		return nil, err
	}
	roles = append(roles, *userRole)

	userData, err = userData.NewUpdater().Roles(roles).Build()
	if err != nil {
		log.Errorf("Erro ao atribuir o papel ao usuário %s: %s", userData.Id(), err.Error())
		return nil, err
	}

	userData, err = instance.userRepository.UpdateUser(*userData)
	if err != nil {
		log.Errorf("Erro ao ativar conta do usuário %s no banco de dados: %s", userData.Id(), err.Error())
		return nil, err
	}

	err = instance.sessionRepository.DeleteSessionsByUserId(userData.Id())
	if err != nil {
		log.Errorf("Erro ao encerrar sessões do usuário %s: %s", userData.Id(), err.Error())
		return nil, err
	}

	sessionId := uuid.New()
	userData, err = userData.NewUpdater().Tokens(sessionId).Build()
	if err != nil {
		log.Errorf("Erro durante geração dos tokens do usuário %s: %s", userData.Id(), err.Error())
		return nil, err
	}

	err = instance.sessionRepository.CreateSession(userData.Id(), sessionId, userData.AccessToken(), userData.RefreshToken())
	if err != nil {
		log.Errorf("Erro durante a criação da sessão do usuário %s: %s", userData.Id(), err.Error())
		return nil, err
	}

	return userData, nil
}
