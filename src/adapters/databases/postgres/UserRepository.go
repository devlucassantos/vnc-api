package postgres

import (
	"github.com/devlucassantos/vnc-domains/src/domains/role"
	"github.com/devlucassantos/vnc-domains/src/domains/user"
	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	"time"
	"vnc-api/adapters/databases/dto"
	"vnc-api/adapters/databases/postgres/queries"
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
		log.Error("Error creating a connection to the Postgres database: ", err.Error())
		return nil, err
	}
	defer instance.connectionManager.closeConnection(postgresConnection)

	transaction, err := postgresConnection.Beginx()
	if err != nil {
		log.Errorf("Error starting transaction to register user %s: %s", userData.Email(), err.Error())
		return nil, err
	}
	defer instance.connectionManager.rollbackTransaction(transaction)

	var userId uuid.UUID
	var userCreationDateAndTime time.Time
	err = transaction.QueryRow(queries.User().Insert(),
		userData.FirstName(), userData.LastName(), userData.Email(), userData.Password(), userData.ActivationCode()).
		Scan(&userId, &userCreationDateAndTime)
	if err != nil {
		log.Errorf("Error registering user %s in the database: %s", userData.Email(), err.Error())
		return nil, err
	}

	var roleSlice []dto.Role
	var userRoles []interface{}
	for _, userRole := range userData.Roles() {
		userRoles = append(userRoles, userRole.Code())
	}

	err = transaction.Select(&roleSlice, queries.Role().Select().ByCodes(len(userRoles)), userRoles...)
	if err != nil {
		log.Errorf("Error retrieving the roles data for user %s from the database: %s", userData.Email(),
			err.Error())
		return nil, err
	}

	for _, roleData := range roleSlice {
		_, err = transaction.Exec(queries.UserRole().Insert(), userId, roleData.Id)
		if err != nil {
			log.Errorf("Error registering role %s for user %s: %s", roleData.Id,
				userData.Email(), err.Error())
			return nil, err
		}
	}

	var roles []role.Role
	for _, roleData := range roleSlice {
		roleDomain, err := role.NewBuilder().
			Id(roleData.Id).
			Code(roleData.Code).
			Build()
		if err != nil {
			log.Errorf("Error validating data for role %s for user %s: %s", roleData.Id, userId, err.Error())
		}

		roles = append(roles, *roleDomain)
	}

	userDomain, err := user.NewBuilder().
		Id(userId).
		FirstName(userData.FirstName()).
		LastName(userData.LastName()).
		Email(userData.Email()).
		HashedPassword(userData.Password()).
		ActivationCode(userData.ActivationCode()).
		CreatedAt(userCreationDateAndTime).
		UpdatedAt(userCreationDateAndTime).
		Roles(roles).
		Build()
	if err != nil {
		log.Errorf("Error validating data for user %s: %s", userData.Email(), err.Error())
		return nil, err
	}

	err = transaction.Commit()
	if err != nil {
		log.Errorf("Error confirming transaction to register user %s: %s", userData.Email(), err.Error())
		return nil, err
	}

	return userDomain, nil
}

func (instance User) UpdateUser(userData user.User) (*user.User, error) {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		log.Error("Error creating a connection to the Postgres database: ", err.Error())
		return nil, err
	}
	defer instance.connectionManager.closeConnection(postgresConnection)

	transaction, err := postgresConnection.Beginx()
	if err != nil {
		log.Errorf("Error starting transaction to update user %s: %s", userData.Id(), err.Error())
		return nil, err
	}
	defer instance.connectionManager.rollbackTransaction(transaction)

	var userUpdateDateAndTime time.Time
	err = transaction.QueryRow(queries.User().Update(),
		userData.FirstName(), userData.LastName(), userData.Email(), userData.Password(), userData.ActivationCode(),
		userData.Id()).Scan(&userUpdateDateAndTime)
	if err != nil {
		log.Errorf("Error updating user %s in the database: %s", userData.Id(), err.Error())
		return nil, err
	}

	var userRoles []interface{}
	for _, userRole := range userData.Roles() {
		userRoles = append(userRoles, userRole.Code())
	}

	var roleSlice []dto.Role
	err = transaction.Select(&roleSlice, queries.Role().Select().ByCodes(len(userRoles)), userRoles...)
	if err != nil {
		log.Errorf("Error retrieving the data for the new roles of user %s from the database: %s", userData.Id(),
			err.Error())
		return nil, err
	}

	var userRolesSlice []dto.Role
	err = transaction.Select(&userRolesSlice, queries.UserRole().Select().ByUserId(), userData.Id())
	if err != nil {
		log.Errorf("Error retrieving the roles data for user %s from the database: %s", userData.Id(),
			err.Error())
		return nil, err
	}

	var rolesRemoved []dto.Role
	for _, userRoleData := range userRolesSlice {
		isAUserRoleRemoved := true
		for _, roleData := range roleSlice {
			if roleData.Id == userRoleData.Id {
				isAUserRoleRemoved = false
				break
			}
		}

		if isAUserRoleRemoved {
			rolesRemoved = append(rolesRemoved, userRoleData)
		}
	}

	var newRoles []dto.Role
	for _, roleData := range roleSlice {
		var isAUserRoleAlreadyRegistered bool
		for _, userRoleData := range userRolesSlice {
			if userRoleData.Id == roleData.Id {
				isAUserRoleAlreadyRegistered = true
				break
			}
		}

		if !isAUserRoleAlreadyRegistered {
			newRoles = append(newRoles, roleData)
		}
	}

	for _, roleData := range rolesRemoved {
		_, err = transaction.Exec(queries.UserRole().Delete(), userData.Id(), roleData.Id)
		if err != nil {
			log.Errorf("Error deleting role %s for user %s: %s", roleData.Id, userData.Id(), err.Error())
			return nil, err
		}
	}

	for _, roleData := range newRoles {
		sqlResult, err := transaction.Exec(queries.UserRole().Update(), userData.Id(), roleData.Id)
		if err != nil {
			log.Errorf("Error activating role %s for user %s: %s", roleData.Id, userData.Id(), err.Error())
			return nil, err
		}

		rowsAffected, err := sqlResult.RowsAffected()
		if err == nil && rowsAffected == 0 {
			_, err = transaction.Exec(queries.UserRole().Insert(), userData.Id(), roleData.Id)
			if err != nil {
				log.Errorf("Error registering role %s for user %s: %s", roleData.Id, userData.Id(), err.Error())
				return nil, err
			}
		} else if err != nil {
			log.Errorf("Error retrieving the number of rows affected by the activation of role %s for user %s "+
				"in the database: %s", roleData.Id, userData.Id(), err.Error())
			return nil, err
		}
	}

	var roles []role.Role
	for _, roleData := range roleSlice {
		roleDomain, err := role.NewBuilder().
			Id(roleData.Id).
			Code(roleData.Code).
			Build()
		if err != nil {
			log.Errorf("Error validating data for role %s for user %s: %s", roleData.Id, userData.Id(), err.Error())
		}

		roles = append(roles, *roleDomain)
	}

	userDomain, err := user.NewBuilder().
		Id(userData.Id()).
		FirstName(userData.FirstName()).
		LastName(userData.LastName()).
		Email(userData.Email()).
		HashedPassword(userData.Password()).
		ActivationCode(userData.ActivationCode()).
		CreatedAt(userData.CreatedAt()).
		UpdatedAt(userUpdateDateAndTime).
		Roles(roles).
		Build()
	if err != nil {
		log.Errorf("Error validating data for user %s: %s", userData.Id(), err.Error())
		return nil, err
	}

	err = transaction.Commit()
	if err != nil {
		log.Errorf("Error confirming transaction to update user %s: %s", userData.Id(), err.Error())
		return nil, err
	}

	return userDomain, nil
}

func (instance User) GetUserById(id uuid.UUID) (*user.User, error) {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		log.Error("Error creating a connection to the Postgres database: ", err.Error())
		return nil, err
	}
	defer instance.connectionManager.closeConnection(postgresConnection)

	var userData dto.User
	err = postgresConnection.Get(&userData, queries.User().Select().ById(), id)
	if err != nil {
		log.Errorf("Error retrieving data for user %s from the database: %s", id, err.Error())
		return nil, err
	}

	var roleSlice []dto.Role
	err = postgresConnection.Select(&roleSlice, queries.UserRole().Select().ByUserId(), userData.Id)
	if err != nil {
		log.Errorf("Error retrieving the roles data for user %s from the database: %s", userData.Id, err.Error())
		return nil, err
	}

	var roles []role.Role
	for _, roleData := range roleSlice {
		roleDomain, err := role.NewBuilder().
			Id(roleData.Id).
			Code(roleData.Code).
			Build()
		if err != nil {
			log.Errorf("Error validating data for role %s for user %s: %s", roleData.Id, userData.Id, err.Error())
		}

		roles = append(roles, *roleDomain)
	}

	userDomain, err := user.NewBuilder().
		Id(userData.Id).
		FirstName(userData.FirstName).
		LastName(userData.LastName).
		Email(userData.Email).
		HashedPassword(userData.Password).
		ActivationCode(userData.ActivationCode).
		CreatedAt(userData.CreatedAt).
		UpdatedAt(userData.UpdatedAt).
		Roles(roles).
		Build()
	if err != nil {
		log.Errorf("Error validating data for user %s: %s", userData.Id, err.Error())
		return nil, err
	}

	return userDomain, nil
}

func (instance User) GetUserByEmail(email string) (*user.User, error) {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		log.Error("Error creating a connection to the Postgres database: ", err.Error())
		return nil, err
	}
	defer instance.connectionManager.closeConnection(postgresConnection)

	var userData dto.User
	err = postgresConnection.Get(&userData, queries.User().Select().ByEmail(), email)
	if err != nil {
		log.Errorf("Error retrieving data for user %s from the database: %s", email, err.Error())
		return nil, err
	}

	var roleSlice []dto.Role
	err = postgresConnection.Select(&roleSlice, queries.UserRole().Select().ByUserId(), userData.Id)
	if err != nil {
		log.Errorf("Error retrieving the roles data for user %s from the database: %s", userData.Email,
			err.Error())
		return nil, err
	}

	var roles []role.Role
	for _, roleData := range roleSlice {
		roleDomain, err := role.NewBuilder().
			Id(roleData.Id).
			Code(roleData.Code).
			Build()
		if err != nil {
			log.Errorf("Error validating data for role %s for user %s: %s", roleData.Id, userData.Id, err.Error())
		}

		roles = append(roles, *roleDomain)
	}

	userDomain, err := user.NewBuilder().
		Id(userData.Id).
		FirstName(userData.FirstName).
		LastName(userData.LastName).
		Email(userData.Email).
		HashedPassword(userData.Password).
		ActivationCode(userData.ActivationCode).
		CreatedAt(userData.CreatedAt).
		UpdatedAt(userData.UpdatedAt).
		Roles(roles).
		Build()
	if err != nil {
		log.Errorf("Error validating data for user %s: %s", userData.Id, err.Error())
		return nil, err
	}

	return userDomain, nil
}
