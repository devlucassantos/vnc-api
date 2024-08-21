package queries

type userRoleSqlManager struct{}

func UserRole() *userRoleSqlManager {
	return &userRoleSqlManager{}
}

func (userRoleSqlManager) Insert() string {
	return `INSERT INTO user_role(user_id, role_id) VALUES ($1, $2)`
}
