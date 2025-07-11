package services

import (
	"errors"
	"github.com/devlucassantos/vnc-domains/src/domains/role"
	"github.com/devlucassantos/vnc-domains/src/domains/user"
	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	"strings"
	"vnc-api/core/interfaces/postgres"
	"vnc-api/core/interfaces/redis"
	"vnc-api/core/interfaces/services"
	"vnc-api/core/services/utils"
)

type User struct {
	userRepository    postgres.User
	sessionRepository redis.Session
	emailService      services.Email
}

func NewUserService(userRepository postgres.User, sessionRepository redis.Session, emailService services.Email) *User {
	return &User{
		userRepository:    userRepository,
		sessionRepository: sessionRepository,
		emailService:      emailService,
	}
}

func (instance User) ResendUserAccountActivationEmail(userId uuid.UUID) error {
	userData, err := instance.userRepository.GetUserById(userId)
	if err != nil {
		log.Errorf("Error retrieving data for user %s from the database: %s", userId, err.Error())
		return err
	}

	if len(userData.Roles()) != 1 || userData.Roles()[0].Code() != role.InactiveUserRoleCode {
		errorMessage := "active account"
		log.Warnf("Error resending the account activation email for user %s: %s", userId, errorMessage)
		return errors.New(errorMessage)
	}

	userActivationCode, err := utils.GenerateUserActivationCode()
	if err != nil {
		log.Errorf("Error generating new account activation email for user %s: %s", userId, err.Error())
		return err
	}

	userData, err = userData.NewUpdater().ActivationCode(userActivationCode).Build()
	if err != nil {
		log.Errorf("Error setting new account activation code for user %s: %s", userId, err.Error())
		return nil
	}

	userData, err = instance.userRepository.UpdateUser(*userData)
	if err != nil {
		log.Errorf("Error updating account activation code for user %s: %s", userId, err.Error())
		return err
	}

	err = instance.emailService.SendUserAccountActivationEmail(*userData)
	if err != nil {
		log.Errorf("Error resending account activation email for user %s: %s", userId, err.Error())
	}

	return nil
}

func (instance User) ActivateUserAccount(userAccountActivationData user.User) (*user.User, error) {
	userData, err := instance.userRepository.GetUserById(userAccountActivationData.Id())
	if err != nil {
		log.Errorf("Error retrieving data for user %s from the database: %s",
			userAccountActivationData.Id(), err.Error())
		return nil, err
	}

	if len(userData.Roles()) != 1 || userData.Roles()[0].Code() != role.InactiveUserRoleCode {
		errorMessage := "active account"
		log.Warnf("Error activating account for user %s: %s", userAccountActivationData.Id(), errorMessage)
		return nil, errors.New(errorMessage)
	}

	if userData.ActivationCode() != strings.ToUpper(userAccountActivationData.ActivationCode()) {
		errorMessage := "invalid activation code"
		log.Warnf("The activation code provided for the account of user %s is invalid: %s", userData.Id(),
			errorMessage)
		return nil, errors.New(errorMessage)
	}

	var roles []role.Role
	userRole, err := role.NewBuilder().Code(role.UserRoleCode).Build()
	if err != nil {
		log.Errorf("Error building the account activation role for user %s: %s", userAccountActivationData.Id(),
			err.Error())
		return nil, err
	}
	roles = append(roles, *userRole)

	userData, err = userData.NewUpdater().Roles(roles).Build()
	if err != nil {
		log.Errorf("Error setting the account activation role for user %s: %s", userAccountActivationData.Id(),
			err.Error())
		return nil, err
	}

	userData, err = instance.userRepository.UpdateUser(*userData)
	if err != nil {
		log.Errorf("Error activating account for user %s in the database: %s", userAccountActivationData.Id(),
			err.Error())
		return nil, err
	}

	err = instance.sessionRepository.DeleteSessionsByUserId(userData.Id())
	if err != nil {
		log.Errorf("Error deleting session for user %s: %s", userAccountActivationData.Id(), err.Error())
		return nil, err
	}

	sessionId := uuid.New()
	userData, err = userData.NewUpdater().Tokens(sessionId).Build()
	if err != nil {
		log.Errorf("Error generating tokens for user %s: %s", userAccountActivationData.Id(), err.Error())
		return nil, err
	}

	err = instance.sessionRepository.CreateSession(userData.Id(), sessionId, userData.AccessToken(),
		userData.RefreshToken())
	if err != nil {
		log.Errorf("Error creating session for user %s: %s", userAccountActivationData.Id(), err.Error())
		return nil, err
	}

	return userData, nil
}
