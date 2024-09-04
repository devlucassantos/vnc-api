package queries

import (
	"fmt"
	"strings"
)

type articleTypeSqlManager struct{}

func ArticleType() *articleTypeSqlManager {
	return &articleTypeSqlManager{}
}

type articleTypeSelectSqlManager struct{}

func (articleTypeSqlManager) Select() *articleTypeSelectSqlManager {
	return &articleTypeSelectSqlManager{}
}

func (articleTypeSelectSqlManager) In(numberOfTypes int) string {
	var parameters []string
	for i := 1; i <= numberOfTypes; i++ {
		parameters = append(parameters, fmt.Sprintf("$%d", i))
	}

	return fmt.Sprintf(`SELECT id AS article_type_id, description AS article_type_description,
       			color AS article_type_color, sort_order AS article_type_sort_order,
       			created_at AS article_type_created_at, updated_at AS article_type_updated_at
			FROM article_type
			WHERE active = true AND id IN (%s)
			ORDER BY article_type.sort_order`, strings.Join(parameters, ","))
}

func (articleTypeSelectSqlManager) All() string {
	return `SELECT id AS article_type_id, description AS article_type_description,
       			color AS article_type_color, sort_order AS article_type_sort_order,
       			created_at AS article_type_created_at, updated_at AS article_type_updated_at
			FROM article_type
			WHERE active = true
			ORDER BY sort_order`
}
