package queries

import (
	"fmt"
	"strings"
)

type roleSqlManager struct{}

func Role() *roleSqlManager {
	return &roleSqlManager{}
}

type roleSelectSqlManager struct{}

func (roleSqlManager) Select() *roleSelectSqlManager {
	return &roleSelectSqlManager{}
}

func (roleSelectSqlManager) ByDescriptions(numberOfRoles int) string {
	var parameters []string
	for i := 1; i <= numberOfRoles; i++ {
		parameters = append(parameters, fmt.Sprintf("$%d", i))
	}

	return fmt.Sprintf(`SELECT id AS role_id, code AS role_code, created_at AS role_created_at,
       			updated_at AS role_updated_at
			FROM role
			WHERE role.active = true AND role.code IN (%s)`, strings.Join(parameters, ","))
}

func (roleSelectSqlManager) ByUserId() string {
	return `SELECT role.id AS role_id, role.code AS role_code, role.created_at AS role_created_at,
       			role.updated_at AS role_updated_at
			FROM role
				INNER JOIN user_role ON user_role.role_id = role.id
			WHERE role.active = true AND user_role.active AND user_id = $1`
}
