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
	return `SELECT article.id AS article_id, article.created_at AS article_created_at,
       			article.updated_at AS article_updated_at,
				COALESCE(AVG(user_article.rating), 0) AS article_average_rating,
				COUNT(user_article.rating) AS article_number_of_ratings,
				article_type.id AS article_type_id, article_type.description AS article_type_description,
				article_type.codes AS article_type_codes, article_type.color AS article_type_color,
				proposition.id AS proposition_id, proposition.original_text_url AS proposition_original_text_url,
				proposition.original_text_mime_type AS proposition_original_text_mime_type,
				proposition.title AS proposition_title, proposition.content AS proposition_content,
				proposition.submitted_at AS proposition_submitted_at,
				COALESCE(proposition.image_url, '') AS proposition_image_url,
				COALESCE(proposition.image_description, '') AS proposition_image_description,
				proposition_type.id AS proposition_type_id,
				proposition_type.description AS proposition_type_description,
				proposition_type.color AS proposition_type_color
			FROM article
			    INNER JOIN article_type ON article_type.id = article.article_type_id
				INNER JOIN proposition ON proposition.article_id = article.id
			    INNER JOIN proposition_type ON proposition_type.id = proposition.proposition_type_id
				LEFT JOIN user_article ON user_article.article_id = article.id
			WHERE article.active = true AND article_type.active = true AND proposition.active = true AND
				proposition_type.active = true AND user_article.active IS NOT false AND article.id = $1
			GROUP BY article.id, article_type.id, proposition.id, proposition_type.id`
}
