package postgres

import (
	"github.com/devlucassantos/vnc-domains/src/domains/role"
	"github.com/devlucassantos/vnc-domains/src/domains/user"
	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	"time"
	"vnc-api/infra/dto"
	"vnc-api/infra/postgres/queries"
)

type User struct {
	connectionManager connectionManagerInterface
}

func NewUserRepository(connectionManager connectionManagerInterface) *User {
	return &User{
		connectionManager: connectionManager,
	}
}

func (instance User) CreateUser(userData user.User) (*user.User, error) {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		log.Error("Erro ao tentar se conectar com o Postgres: ", err.Error())
		return nil, err
	}
	defer instance.connectionManager.closeConnection(postgresConnection)

	transaction, err := postgresConnection.Beginx()
	if err != nil {
		log.Errorf("Erro ao iniciar transação para o cadastro do usuário %s: %s", userData.Email(),
			err.Error())
		return nil, err
	}
	defer instance.connectionManager.rollbackTransaction(transaction)

	var userId uuid.UUID
	var userCreationDateAndTime time.Time
	err = transaction.QueryRow(queries.User().Insert(),
		userData.FirstName(), userData.LastName(), userData.Email(), userData.Password()).
		Scan(&userId, &userCreationDateAndTime)
	if err != nil {
		log.Errorf("Erro ao cadastrar usuário %s no banco de dados: %s", userData.Email(), err.Error())
		return nil, err
	}

	var roleList []dto.Role
	var userRoles []interface{}
	for _, userRole := range userData.Roles() {
		userRoles = append(userRoles, userRole.Code())
	}

	err = postgresConnection.Select(&roleList, queries.Role().Select().ByDescriptions(len(userRoles)), userRoles...)
	if err != nil {
		log.Errorf("Erro ao obter os dados dos papeis do usuário %s no banco de dados: %s", userData.Email(),
			err.Error())
		return nil, err
	}

	for _, roleData := range roleList {
		_, err = transaction.Exec(queries.UserRole().Insert(), userId, roleData.Id)
		if err != nil {
			log.Errorf("Erro ao cadastrar papel %s para o usuário %s no banco de dados: %s", roleData.Id,
				userData.Email(), err.Error())
			return nil, err
		}
	}

	var roles []role.Role
	for _, roleData := range roleList {
		roleDomain, err := role.NewBuilder().
			Id(roleData.Id).
			Code(roleData.Code).
			CreatedAt(roleData.CreatedAt).
			UpdatedAt(roleData.UpdatedAt).
			Build()
		if err != nil {
			log.Errorf("Erro ao validar os dados do tipo de proposição %s: %s", roleData.Id, err.Error())
		}

		roles = append(roles, *roleDomain)
	}

	userDomain, err := user.NewBuilder().
		Id(userId).
		FirstName(userData.FirstName()).
		LastName(userData.LastName()).
		Email(userData.Email()).
		HashedPassword(userData.Password()).
		CreatedAt(userCreationDateAndTime).
		UpdatedAt(userCreationDateAndTime).
		Roles(roles).
		Build()
	if err != nil {
		log.Errorf("Erro ao validar os dados do usuário %s: %s", userData.Email(), err.Error())
		return nil, err
	}

	err = transaction.Commit()
	if err != nil {
		log.Errorf("Erro ao confirmar transação para o cadastro do usuário %s: %s", userData.Email(), err.Error())
		return nil, err
	}

	return userDomain, nil
}

func (instance User) GetUserById(id uuid.UUID) (*user.User, error) {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		log.Error("Erro ao tentar se conectar com o Postgres: ", err.Error())
		return nil, err
	}
	defer instance.connectionManager.closeConnection(postgresConnection)

	var userData dto.User
	err = postgresConnection.Get(&userData, queries.User().Select().ById(), id)
	if err != nil {
		log.Errorf("Erro ao obter os dados do usuário %s no banco de dados: %s", id, err.Error())
		return nil, err
	}

	var roleList []dto.Role
	err = postgresConnection.Select(&roleList, queries.Role().Select().ByUserId(), userData.Id)
	if err != nil {
		log.Errorf("Erro ao obter os dados dos papeis do usuário %s no banco de dados: %s", userData.Id, err.Error())
		return nil, err
	}

	var roles []role.Role
	for _, roleData := range roleList {
		roleDomain, err := role.NewBuilder().
			Id(roleData.Id).
			Code(roleData.Code).
			CreatedAt(roleData.CreatedAt).
			UpdatedAt(roleData.UpdatedAt).
			Build()
		if err != nil {
			log.Errorf("Erro ao validar os dados do papel %s: %s", roleData.Id, err.Error())
		}

		roles = append(roles, *roleDomain)
	}

	userDomain, err := user.NewBuilder().
		Id(userData.Id).
		FirstName(userData.FirstName).
		LastName(userData.LastName).
		Email(userData.Email).
		HashedPassword(userData.Password).
		CreatedAt(userData.CreatedAt).
		UpdatedAt(userData.UpdatedAt).
		Roles(roles).
		Build()
	if err != nil {
		log.Errorf("Erro ao validar os dados do usuário %s: %s", userData.Id, err.Error())
		return nil, err
	}

	return userDomain, nil
}

func (instance User) GetUserByEmail(email string) (*user.User, error) {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		log.Error("Erro ao tentar se conectar com o Postgres: ", err.Error())
		return nil, err
	}
	defer instance.connectionManager.closeConnection(postgresConnection)

	var userData dto.User
	err = postgresConnection.Get(&userData, queries.User().Select().ByEmail(), email)
	if err != nil {
		log.Errorf("Erro ao obter os dados do usuário %s no banco de dados: %s", email, err.Error())
		return nil, err
	}

	var roleList []dto.Role
	err = postgresConnection.Select(&roleList, queries.Role().Select().ByUserId(), userData.Id)
	if err != nil {
		log.Errorf("Erro ao obter os dados dos papeis do usuário %s no banco de dados: %s", userData.Email,
			err.Error())
		return nil, err
	}

	var roles []role.Role
	for _, roleData := range roleList {
		roleDomain, err := role.NewBuilder().
			Id(roleData.Id).
			Code(roleData.Code).
			CreatedAt(roleData.CreatedAt).
			UpdatedAt(roleData.UpdatedAt).
			Build()
		if err != nil {
			log.Errorf("Erro ao validar os dados do papel %s: %s", roleData.Id, err.Error())
		}

		roles = append(roles, *roleDomain)
	}

	userDomain, err := user.NewBuilder().
		Id(userData.Id).
		FirstName(userData.FirstName).
		LastName(userData.LastName).
		Email(userData.Email).
		HashedPassword(userData.Password).
		CreatedAt(userData.CreatedAt).
		UpdatedAt(userData.UpdatedAt).
		Roles(roles).
		Build()
	if err != nil {
		log.Errorf("Erro ao validar os dados do usuário %s: %s", userData.Id, err.Error())
		return nil, err
	}

	return userDomain, nil
}
