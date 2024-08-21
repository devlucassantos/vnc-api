package queries

import (
	"fmt"
	"strings"
)

type propositionTypeSqlManager struct{}

func PropositionType() *propositionTypeSqlManager {
	return &propositionTypeSqlManager{}
}

type propositionTypeSelectSqlManager struct{}

func (propositionTypeSqlManager) Select() *propositionTypeSelectSqlManager {
	return &propositionTypeSelectSqlManager{}
}

func (propositionTypeSelectSqlManager) In(numberOfTypes int) string {
	var parameters []string
	for i := 1; i <= numberOfTypes; i++ {
		parameters = append(parameters, fmt.Sprintf("$%d", i))
	}

	return fmt.Sprintf(`SELECT id AS proposition_type_id, description AS proposition_type_description,
       			color AS proposition_type_color, sort_order AS proposition_type_sort_order,
       			created_at AS proposition_type_created_at, updated_at AS proposition_type_updated_at
			FROM proposition_type
			WHERE active = true AND id IN (%s)
			ORDER BY proposition_type.sort_order`, strings.Join(parameters, ","))
}

func (propositionTypeSelectSqlManager) All() string {
	return `SELECT id AS proposition_type_id, description AS proposition_type_description,
       			color AS proposition_type_color, sort_order AS proposition_type_sort_order,
       			created_at AS proposition_type_created_at, updated_at AS proposition_type_updated_at
			FROM proposition_type
			WHERE active = true
			ORDER BY sort_order`
}
