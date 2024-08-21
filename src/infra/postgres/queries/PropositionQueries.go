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
				proposition.id AS proposition_id, proposition.original_text_url AS proposition_original_text_url,
				proposition.title AS proposition_title, proposition.content AS proposition_content,
				proposition.submitted_at AS proposition_submitted_at, proposition.created_at AS proposition_created_at,
				proposition.updated_at AS proposition_updated_at,
				proposition_type.id AS proposition_type_id, proposition_type.description AS proposition_type_description,
       			proposition_type.color AS proposition_type_color, proposition_type.sort_order AS proposition_type_sort_order,
       			proposition_type.created_at AS proposition_type_created_at, proposition_type.updated_at AS proposition_type_updated_at
			FROM article
				LEFT JOIN proposition ON proposition.id = article.proposition_id
			    LEFT JOIN proposition_type ON proposition_type.id = proposition.proposition_type_id
				LEFT JOIN user_article ON user_article.article_id = article.id
			WHERE article.active = true AND proposition.active = true AND article.id = $1
			GROUP BY article.id, article.reference_date_time, proposition.id, proposition_type.id
			ORDER BY article.reference_date_time`
}
