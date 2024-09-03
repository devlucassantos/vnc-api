package services

import (
	"encoding/hex"
	"errors"
	"github.com/devlucassantos/vnc-domains/src/domains/user"
	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	"golang.org/x/crypto/bcrypt"
	"vnc-api/core/interfaces/repositories"
	"vnc-api/core/services/utils"
)

type Authentication struct {
	userRepository    repositories.User
	sessionRepository repositories.Session
}

func NewAuthenticationService(userRepository repositories.User, sessionRepository repositories.Session) *Authentication {
	return &Authentication{
		userRepository:    userRepository,
		sessionRepository: sessionRepository,
	}
}

func (instance Authentication) SignUp(signUpData user.User) (*user.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(signUpData.Password()), bcrypt.DefaultCost)
	if err != nil {
		log.Errorf("Erro ao criptografar a senha do usuário %s: %s", signUpData.Email(), err.Error())
		return nil, err
	}

	userActivationCode, err := utils.GenerateUserActivationCode()
	if err != nil {
		log.Errorf("Erro ao gerar código de ativação durante a criação da conta do usuário %s: %s",
			signUpData.Email(), err.Error())
		return nil, err
	}

	userData, err := signUpData.NewUpdater().
		HashedPassword(hex.EncodeToString(hashedPassword)).
		ActivationCode(userActivationCode).
		Build()
	if err != nil {
		log.Errorf("Erro ao definir a senha com hash e o código de ativação da conta do usuário %s: %s",
			signUpData.Email(), err.Error())
		return nil, err
	}

	userData, err = instance.userRepository.CreateUser(*userData)
	if err != nil {
		log.Errorf("Erro ao cadastrar usuário %s no banco de dados: %s", signUpData.Email(), err.Error())
		return nil, err
	}

	sessionId := uuid.New()
	userData, err = userData.NewUpdater().Tokens(sessionId).Build()
	if err != nil {
		log.Errorf("Erro durante geração dos tokens do usuário %s: %s", userData.Email(), err.Error())
		return nil, err
	}

	err = instance.sessionRepository.CreateSession(userData.Id(), sessionId, userData.AccessToken(), userData.RefreshToken())
	if err != nil {
		log.Errorf("Erro durante a criação da sessão do usuário %s: %s", userData.Email(), err.Error())
		return nil, err
	}

	go func() {
		err = utils.SendActivationEmail(*userData)
		if err != nil {
			log.Error("Erro ao enviar email de ativação da conta do usuário %s durante a criação da conta: %s",
				userData.Email(), err)
		}
	}()

	return userData, nil
}

func (instance Authentication) SignIn(signInData user.User) (*user.User, error) {
	userData, err := instance.userRepository.GetUserByEmail(signInData.Email())
	if err != nil {
		log.Errorf("Erro ao obter os dados do usuário %s no banco de dados: %s", signInData.Email(), err.Error())
		return nil, err
	}

	decodedPassword, err := hex.DecodeString(userData.Password())
	if err != nil {
		log.Errorf("Erro ao decodificar a senha do usuário %s: %s", userData.Email(), err.Error())
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword(decodedPassword, []byte(signInData.Password()))
	if err != nil {
		log.Errorf("Erro ao comparar senha do usuário %s: %s", userData.Email(), err.Error())
		return nil, errors.New("senha incorreta")
	}

	sessionId := uuid.New()
	userData, err = userData.NewUpdater().Tokens(sessionId).Build()
	if err != nil {
		log.Errorf("Erro durante geração dos tokens do usuário %s: %s", userData.Email(), err.Error())
		return nil, err
	}

	err = instance.sessionRepository.CreateSession(userData.Id(), sessionId, userData.AccessToken(), userData.RefreshToken())
	if err != nil {
		log.Errorf("Erro durante a criação da sessão do usuário %s: %s", userData.Email(), err.Error())
		return nil, err
	}

	return userData, nil
}

func (instance Authentication) SignOut(userId uuid.UUID, sessionId uuid.UUID) error {
	err := instance.sessionRepository.DeleteSession(userId, sessionId)
	if err != nil {
		log.Errorf("Erro ao encerrar sessão do usuário %s: %s", userId, err.Error())
		return err
	}

	return nil
}

func (instance Authentication) RefreshTokens(userId uuid.UUID, sessionId uuid.UUID, refreshToken string) (*user.User, error) {
	userData, err := instance.userRepository.GetUserById(userId)
	if err != nil {
		log.Errorf("Erro ao obter os dados do usuário %s no banco de dados: %s", userId, err.Error())
		return nil, err
	}

	tokenExists, err := instance.sessionRepository.RefreshTokenExists(userData.Id(), sessionId, refreshToken)
	if err != nil {
		log.Errorf("Erro ao verificar se o token de atualização do usuário %s existe: %s", userData.Email(), err.Error())
		return nil, err
	}

	if !tokenExists {
		log.Errorf("O token de atualização do usuário %s é inválido", userData.Email())
		return nil, errors.New("token inválido")
	}

	err = instance.sessionRepository.DeleteSession(userData.Id(), sessionId)
	if err != nil {
		log.Errorf("Erro ao remover token de acesso do usuário %s no banco de dados: %s", userData.Email(),
			err.Error())
		return nil, err
	}

	userData, err = userData.NewUpdater().Tokens(sessionId).Build()
	if err != nil {
		log.Errorf("Erro durante geração dos novos tokens do usuário %s: %s", userData.Email(), err.Error())
		return nil, err
	}

	err = instance.sessionRepository.CreateSession(userData.Id(), sessionId, userData.AccessToken(), userData.RefreshToken())
	if err != nil {
		log.Errorf("Erro durante a atualização da sessão do usuário %s: %s", userData.Email(), err.Error())
		return nil, err
	}

	return userData, nil
}

func (instance Authentication) SessionExists(userId uuid.UUID, sessionId uuid.UUID, token string) (bool, error) {
	sessionExists, err := instance.sessionRepository.SessionExists(userId, sessionId, token)
	if err != nil {
		log.Errorf("Erro ao verificar se o token de acesso do usuário %s existe: %s", userId, err.Error())
		return false, err
	}

	return sessionExists, nil
}
