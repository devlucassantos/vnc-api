package queries

type propositionKeywordSqlManager struct{}

func PropositionKeyword() *propositionKeywordSqlManager {
	return &propositionKeywordSqlManager{}
}

type propositionKeywordSelectSqlManager struct{}

func (propositionKeywordSqlManager) Select() *propositionKeywordSelectSqlManager {
	return &propositionKeywordSelectSqlManager{}
}

func (propositionKeywordSelectSqlManager) ByPropositionId() string {
	return `SELECT COALESCE(keyword.id, '00000000-0000-0000-0000-000000000000') AS keyword_id,
       			COALESCE(keyword.keyword, '') AS keyword_keyword,
       			COALESCE(keyword.active, true) AS keyword_active,
       			COALESCE(keyword.created_at, '1970-01-01 00:00:00') AS keyword_created_at,
       			COALESCE(keyword.updated_at, '1970-01-01 00:00:00') AS keyword_updated_at
			FROM proposition
			INNER JOIN proposition_keyword ON proposition_keyword.proposition_id = proposition.id
			INNER JOIN keyword ON keyword.id = proposition_keyword.keyword_id
			WHERE proposition.active = true AND proposition_keyword.active = true AND keyword.active = true AND
				proposition.id = $1`
}
