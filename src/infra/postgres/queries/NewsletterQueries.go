package queries

type newsletterSqlManager struct{}

func Newsletter() *newsletterSqlManager {
	return &newsletterSqlManager{}
}

type newsletterSelectSqlManager struct{}

func (newsletterSqlManager) Select() *newsletterSelectSqlManager {
	return &newsletterSelectSqlManager{}
}

func (newsletterSelectSqlManager) ById() string {
	return `SELECT COALESCE(newsletter.id, '00000000-0000-0000-0000-000000000000') AS newsletter_id,
       			COALESCE(newsletter.title, '') AS newsletter_title,
    			COALESCE(newsletter.content, '') AS newsletter_content,
    			COALESCE(newsletter.reference_date, '1970-01-01 00:00:00') AS newsletter_reference_date,
       			COALESCE(newsletter.active, true) AS newsletter_active,
    			COALESCE(newsletter.created_at, '1970-01-01 00:00:00') AS newsletter_created_at,
    			COALESCE(newsletter.updated_at, '1970-01-01 00:00:00') AS newsletter_updated_at,
    			COALESCE(news.id, '00000000-0000-0000-0000-000000000000') AS newsletter_news_id
    		FROM newsletter
    			INNER JOIN news ON news.newsletter_id = newsletter.id
    		WHERE newsletter.active = true AND news.active = true AND newsletter.id = $1`
}
