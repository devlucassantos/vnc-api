package queries

type eventSqlManager struct{}

func Event() *eventSqlManager {
	return &eventSqlManager{}
}

type eventSelectSqlManager struct{}

func (eventSqlManager) Select() *eventSelectSqlManager {
	return &eventSelectSqlManager{}
}

func (eventSelectSqlManager) ByArticleId() string {
	return `SELECT article.id AS article_id, article.created_at AS article_created_at,
				article.updated_at AS article_updated_at,
				COALESCE(AVG(user_article.rating), 0) AS article_average_rating,
				COUNT(user_article.rating) AS article_number_of_ratings,
				article_type.id AS article_type_id, article_type.description AS article_type_description,
				article_type.codes AS article_type_codes, article_type.color AS article_type_color,
				event.id AS event_id, event.code AS event_code, event.title AS event_title,
       			event.description AS event_description, event.starts_at AS event_starts_at,
				COALESCE(event.ends_at, '1970-01-01 00:00:00') AS event_ends_at, event.location AS event_location,
				event.is_internal AS event_is_internal, COALESCE(event.video_url, '') AS event_video_url,
				event_type.id AS event_type_id, event_type.description AS event_type_description,
				event_type.color AS event_type_color,
				event_situation.id AS event_situation_id, event_situation.description AS event_situation_description,
				event_situation.color AS event_situation_color
			FROM article
			    INNER JOIN article_type ON article_type.id = article.article_type_id
			    INNER JOIN event ON event.article_id = article.id
				INNER JOIN event_type ON event_type.id = event.event_type_id
				INNER JOIN event_situation ON event_situation.id = event.event_situation_id
				LEFT JOIN user_article ON user_article.article_id = article.id
			WHERE article.active = true AND article_type.active = true AND event.active = true AND
				event_type.active = true AND event_situation.active = true AND user_article.active IS NOT false AND
				article.id = $1
			GROUP BY article.id, article_type.id, event.id, event_type.id, event_situation.id`
}
