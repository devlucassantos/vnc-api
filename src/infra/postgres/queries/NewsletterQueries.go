package queries

type newsletterSqlManager struct{}

func Newsletter() *newsletterSqlManager {
	return &newsletterSqlManager{}
}

type newsletterSelectSqlManager struct{}

func (newsletterSqlManager) Select() *newsletterSelectSqlManager {
	return &newsletterSelectSqlManager{}
}

func (newsletterSelectSqlManager) ByArticleId() string {
	return `SELECT article.id AS article_id, article.created_at AS article_created_at,
       			article.updated_at AS article_updated_at,
				COALESCE(AVG(user_article.rating), 0) AS article_average_rating,
				COUNT(user_article.rating) AS article_number_of_ratings,
				article_type.id AS article_type_id, article_type.description AS article_type_description,
				article_type.codes AS article_type_codes, article_type.color AS article_type_color,
				newsletter.id AS newsletter_id, newsletter.reference_date AS newsletter_reference_date,
				newsletter.description AS newsletter_description
			FROM article
			    INNER JOIN article_type ON article_type.id = article.article_type_id
				INNER JOIN newsletter ON newsletter.article_id = article.id
				LEFT JOIN user_article ON user_article.article_id = article.id
			WHERE article.active = true AND article_type.active = true AND newsletter.active = true AND
				user_article.active IS NOT false AND article.id = $1
			GROUP BY article.id, article_type.id, newsletter.id`
}
