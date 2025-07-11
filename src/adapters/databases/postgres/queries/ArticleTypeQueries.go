package queries

type articleTypeSqlManager struct{}

func ArticleType() *articleTypeSqlManager {
	return &articleTypeSqlManager{}
}

type articleTypeSelectSqlManager struct{}

func (articleTypeSqlManager) Select() *articleTypeSelectSqlManager {
	return &articleTypeSelectSqlManager{}
}

func (articleTypeSelectSqlManager) All() string {
	return `SELECT id AS article_type_id, description AS article_type_description, codes AS article_type_codes, 
       			color AS article_type_color
			FROM article_type
			WHERE active = true
			ORDER BY sort_order`
}
