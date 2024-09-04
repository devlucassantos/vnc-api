package queries

type userRoleSqlManager struct{}

func UserRole() *userRoleSqlManager {
	return &userRoleSqlManager{}
}

func (userRoleSqlManager) Insert() string {
	return `INSERT INTO user_role(user_id, role_id) VALUES ($1, $2)`
}

func (userRoleSqlManager) Update() string {
	return `UPDATE user_role
			SET active = true, updated_at = TIMEZONE('America/Sao_Paulo'::TEXT, NOW())
            WHERE user_id = $1 AND role_id = $2`
}

func (userRoleSqlManager) Delete() string {
	return `UPDATE user_role
			SET active = false, updated_at = TIMEZONE('America/Sao_Paulo'::TEXT, NOW())
            WHERE user_id = $1 AND role_id = $2`
}

type userRoleSelectSqlManager struct{}

func (userRoleSqlManager) Select() *userRoleSelectSqlManager {
	return &userRoleSelectSqlManager{}
}

func (userRoleSelectSqlManager) ByUserId() string {
	return `SELECT role.id AS role_id, role.code AS role_code, role.created_at AS role_created_at,
       			role.updated_at AS role_updated_at
			FROM role
				INNER JOIN user_role ON user_role.role_id = role.id
			WHERE role.active = true AND user_role.active = true AND user_id = $1`
}
