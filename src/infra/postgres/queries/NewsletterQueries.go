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
	return `SELECT article.id AS article_id, article.reference_date_time AS article_reference_date_time,
				article.created_at AS article_created_at, article.updated_at AS article_updated_at,
				COALESCE(AVG(user_article.rating), 0) AS article_average_rating,
				COUNT(user_article.rating) AS article_number_of_ratings,
				article_type.id AS article_type_id, article_type.description AS article_type_description,
       			article_type.color AS article_type_color, article_type.sort_order AS article_type_sort_order,
       			article_type.created_at AS article_type_created_at, article_type.updated_at AS article_type_updated_at,
				newsletter.id AS newsletter_id, newsletter.reference_date AS newsletter_reference_date,
				newsletter.title AS newsletter_title, newsletter.description AS newsletter_description,
				newsletter.created_at AS newsletter_created_at, newsletter.updated_at AS newsletter_updated_at
			FROM article
			    INNER JOIN article_type ON article_type.id = article.article_type_id
				LEFT JOIN newsletter ON newsletter.id = article.newsletter_id
				LEFT JOIN user_article ON user_article.article_id = article.id
			WHERE article.active = true AND newsletter.active = true AND article.id = $1
			GROUP BY article.id, article_type.id, newsletter.id`
}
