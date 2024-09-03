package queries

type externalAuthorSqlManager struct{}

func ExternalAuthor() *externalAuthorSqlManager {
	return &externalAuthorSqlManager{}
}

type externalAuthorSelectSqlManager struct{}

func (externalAuthorSqlManager) Select() *externalAuthorSelectSqlManager {
	return &externalAuthorSelectSqlManager{}
}

func (externalAuthorSelectSqlManager) All() string {
	return `SELECT id AS external_author_id, name AS external_author_name, type AS external_author_type,
       			created_at AS external_author_created_at, updated_at AS external_author_updated_at
    		FROM external_author
    		WHERE external_author.active = true
    		ORDER BY name, type`
}
