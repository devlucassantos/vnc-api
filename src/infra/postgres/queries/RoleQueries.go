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

func (roleSelectSqlManager) ByCodes(numberOfRoles int) string {
	var parameters []string
	for i := 1; i <= numberOfRoles; i++ {
		parameters = append(parameters, fmt.Sprintf("$%d", i))
	}

	return fmt.Sprintf(`SELECT id AS role_id, code AS role_code
			FROM role
			WHERE role.active = true AND role.code IN (%s)`, strings.Join(parameters, ","))
}
