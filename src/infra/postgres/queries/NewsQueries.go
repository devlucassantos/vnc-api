package queries

type newsSqlManager struct{}

func News() *newsSqlManager {
	return &newsSqlManager{}
}

type newsSelectSqlManager struct{}

func (newsSqlManager) Select() *newsSelectSqlManager {
	return &newsSelectSqlManager{}
}

func (newsSelectSqlManager) TotalNumber() string {
	return `SELECT COUNT(*) FROM news`
}

func (newsSelectSqlManager) All() string {
	return `SELECT COALESCE(proposition.id, '00000000-0000-0000-0000-000000000000') AS proposition_id,
    			COALESCE(proposition.code, 0) AS proposition_code,
    			COALESCE(proposition.original_text_url, '') AS proposition_original_text_url,
       			COALESCE(proposition.title, '') AS proposition_title,
    			COALESCE(proposition.content, '') AS proposition_content,
    			COALESCE(proposition.submitted_at, '1970-01-01 00:00:00') AS proposition_submitted_at,
       			COALESCE(proposition.active, true) AS proposition_active,
    			COALESCE(proposition.created_at, '1970-01-01 00:00:00') AS proposition_created_at,
    			COALESCE(proposition.updated_at, '1970-01-01 00:00:00') AS proposition_updated_at,
    			
				COALESCE(newsletter.id, '00000000-0000-0000-0000-000000000000') AS newsletter_id,
       			COALESCE(newsletter.title, '') AS newsletter_title,
    			COALESCE(newsletter.content, '') AS newsletter_content,
    			COALESCE(newsletter.reference_date, '1970-01-01 00:00:00') AS newsletter_reference_date,
       			COALESCE(newsletter.active, true) AS newsletter_active,
    			COALESCE(newsletter.created_at, '1970-01-01 00:00:00') AS newsletter_created_at,
    			COALESCE(newsletter.updated_at, '1970-01-01 00:00:00') AS newsletter_updated_at
    		FROM news
    		    LEFT JOIN proposition ON proposition.id = news.proposition_id
    			LEFT JOIN newsletter ON newsletter.id = news.newsletter_id
    		WHERE news.active = true AND (proposition.active = true OR newsletter.active = true)
    		ORDER BY news.created_at DESC OFFSET $1 LIMIT $2`
}
