package services

import (
	"encoding/hex"
	"errors"
	"github.com/devlucassantos/vnc-domains/src/domains/user"
	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	"golang.org/x/crypto/bcrypt"
	"vnc-api/core/interfaces/postgres"
	"vnc-api/core/interfaces/redis"
	"vnc-api/core/interfaces/services"
	"vnc-api/core/services/utils"
)

type Authentication struct {
	userRepository    postgres.User
	sessionRepository redis.Session
	emailService      services.Email
}

func NewAuthenticationService(userRepository postgres.User, sessionRepository redis.Session,
	emailService services.Email) *Authentication {
	return &Authentication{
		userRepository:    userRepository,
		sessionRepository: sessionRepository,
		emailService:      emailService,
	}
}

func (instance Authentication) SignUp(signUpData user.User) (*user.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(signUpData.Password()), bcrypt.DefaultCost)
	if err != nil {
		log.Errorf("Error encrypting the password of user %s: %s", signUpData.Email(), err.Error())
		return nil, err
	}

	userActivationCode, err := utils.GenerateUserActivationCode()
	if err != nil {
		log.Errorf("Error generating activation code during the account creation of user %s: %s",
			signUpData.Email(), err.Error())
		return nil, err
	}

	userData, err := signUpData.NewUpdater().
		HashedPassword(hex.EncodeToString(hashedPassword)).
		ActivationCode(userActivationCode).
		Build()
	if err != nil {
		log.Errorf("Error setting the hashed password and activation code for the account of user %s: %s",
			signUpData.Email(), err.Error())
		return nil, err
	}

	userData, err = instance.userRepository.CreateUser(*userData)
	if err != nil {
		log.Errorf("Error registering user %s: %s", signUpData.Email(), err.Error())
		return nil, err
	}

	sessionId := uuid.New()
	userData, err = userData.NewUpdater().Tokens(sessionId).Build()
	if err != nil {
		log.Errorf("Error generating tokens for user %s: %s", signUpData.Email(), err.Error())
		return nil, err
	}

	err = instance.sessionRepository.CreateSession(userData.Id(), sessionId, userData.AccessToken(),
		userData.RefreshToken())
	if err != nil {
		log.Errorf("Error creating session for user %s: %s", signUpData.Email(), err.Error())
		return nil, err
	}

	go func() {
		err = instance.emailService.SendUserAccountActivationEmail(*userData)
		if err != nil {
			log.Errorf("Error sending account activation email for user %s when creating account: %s",
				signUpData.Email(), err.Error())
		}
	}()

	return userData, nil
}

func (instance Authentication) SignIn(signInData user.User) (*user.User, error) {
	userData, err := instance.userRepository.GetUserByEmail(signInData.Email())
	if err != nil {
		log.Errorf("Error retrieving data for user %s from the database: %s", signInData.Email(), err.Error())
		return nil, err
	}

	decodedPassword, err := hex.DecodeString(userData.Password())
	if err != nil {
		log.Errorf("Error decoding password for user %s: %s", signInData.Email(), err.Error())
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword(decodedPassword, []byte(signInData.Password()))
	if err != nil {
		log.Warnf("Error comparing password for user %s: %s", signInData.Email(), err.Error())
		return nil, errors.New("incorrect password")
	}

	sessionId := uuid.New()
	userData, err = userData.NewUpdater().Tokens(sessionId).Build()
	if err != nil {
		log.Errorf("Error generating tokens for user %s: %s", signInData.Email(), err.Error())
		return nil, err
	}

	err = instance.sessionRepository.CreateSession(userData.Id(), sessionId, userData.AccessToken(), userData.RefreshToken())
	if err != nil {
		log.Errorf("Error creating session for user %s: %s", signInData.Email(), err.Error())
		return nil, err
	}

	return userData, nil
}

func (instance Authentication) SignOut(userId uuid.UUID, sessionId uuid.UUID) error {
	err := instance.sessionRepository.DeleteSession(userId, sessionId)
	if err != nil {
		log.Errorf("Error signing out user %s: %s", userId, err.Error())
		return err
	}

	return nil
}

func (instance Authentication) RefreshTokens(userId uuid.UUID, sessionId uuid.UUID, refreshToken string) (*user.User, error) {
	userData, err := instance.userRepository.GetUserById(userId)
	if err != nil {
		log.Errorf("Error retrieving data for user %s from the database: %s", userId, err.Error())
		return nil, err
	}

	tokenExists, err := instance.sessionRepository.RefreshTokenExists(userData.Id(), sessionId, refreshToken)
	if err != nil {
		log.Errorf("Error checking if the refresh token for user %s exists: %s", userId, err.Error())
		return nil, err
	}

	if !tokenExists {
		log.Warnf("The refresh token for user %s is invalid", userId)
		return nil, errors.New("invalid token")
	}

	err = instance.sessionRepository.DeleteSession(userData.Id(), sessionId)
	if err != nil {
		log.Errorf("Error removing access token for user %s from the database: %s", userId, err.Error())
		return nil, err
	}

	userData, err = userData.NewUpdater().Tokens(sessionId).Build()
	if err != nil {
		log.Errorf("Error generating tokens for user %s: %s", userId, err.Error())
		return nil, err
	}

	err = instance.sessionRepository.CreateSession(userData.Id(), sessionId, userData.AccessToken(),
		userData.RefreshToken())
	if err != nil {
		log.Errorf("Error updating the session for user %s: %s", userId, err.Error())
		return nil, err
	}

	return userData, nil
}

func (instance Authentication) SessionExists(userId uuid.UUID, sessionId uuid.UUID, token string) (bool, error) {
	sessionExists, err := instance.sessionRepository.SessionExists(userId, sessionId, token)
	if err != nil {
		log.Errorf("Error checking if the access token for user %s exists: %s", userId, err.Error())
		return false, err
	}

	return sessionExists, nil
}
