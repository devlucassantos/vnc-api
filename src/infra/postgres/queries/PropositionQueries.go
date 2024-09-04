package queries

type propositionSqlManager struct{}

func Proposition() *propositionSqlManager {
	return &propositionSqlManager{}
}

type propositionSelectSqlManager struct{}

func (propositionSqlManager) Select() *propositionSelectSqlManager {
	return &propositionSelectSqlManager{}
}

func (propositionSelectSqlManager) ByArticleId() string {
	return `SELECT article.id AS article_id, article.reference_date_time AS article_reference_date_time,
				article.created_at AS article_created_at, article.updated_at AS article_updated_at,
				COALESCE(AVG(user_article.rating), 0) AS article_average_rating,
				COUNT(user_article.rating) AS article_number_of_ratings,
				article_type.id AS article_type_id, article_type.description AS article_type_description,
       			article_type.color AS article_type_color, article_type.sort_order AS article_type_sort_order,
       			article_type.created_at AS article_type_created_at, article_type.updated_at AS article_type_updated_at,
				proposition.id AS proposition_id, proposition.original_text_url AS proposition_original_text_url,
				proposition.title AS proposition_title, proposition.content AS proposition_content,
				proposition.submitted_at AS proposition_submitted_at, proposition.image_url AS proposition_image_url,
				proposition.created_at AS proposition_created_at, proposition.updated_at AS proposition_updated_at
			FROM article
			    INNER JOIN article_type ON article_type.id = article.article_type_id
				LEFT JOIN proposition ON proposition.id = article.proposition_id
				LEFT JOIN user_article ON user_article.article_id = article.id
			WHERE article.active = true AND proposition.active = true AND article.id = $1
			GROUP BY article.id, article.reference_date_time, article_type.id, proposition.id
			ORDER BY article.reference_date_time`
}
