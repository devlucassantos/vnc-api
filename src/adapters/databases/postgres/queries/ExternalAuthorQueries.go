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
	return `SELECT external_author.id AS external_author_id,
       			external_author.name AS external_author_name,
       			external_author_type.id AS external_author_type_id,
       			external_author_type.description AS external_author_type_description
    		FROM external_author
    			INNER JOIN external_author_type ON external_author_type.id = external_author.external_author_type_id
    		WHERE external_author.active = true AND external_author_type.active = true
    		ORDER BY external_author.name, external_author_type.description`
}
