package queries

type votingSqlManager struct{}

func Voting() *votingSqlManager {
	return &votingSqlManager{}
}

type votingSelectSqlManager struct{}

func (votingSqlManager) Select() *votingSelectSqlManager {
	return &votingSelectSqlManager{}
}

func (votingSelectSqlManager) ByArticleId() string {
	return `SELECT article.id AS article_id, article.created_at AS article_created_at,
				article.updated_at AS article_updated_at,
				COALESCE(AVG(user_article.rating), 0) AS article_average_rating,
				COUNT(user_article.rating) AS article_number_of_ratings,
				article_type.id AS article_type_id, article_type.description AS article_type_description,
				article_type.codes AS article_type_codes, article_type.color AS article_type_color,
    			voting.id AS voting_id, voting.code AS voting_code, voting.description AS voting_description,
    			voting.result AS voting_result, voting.result_announced_at AS voting_result_announced_at,
    			voting.is_approved AS voting_is_approved, voting.main_proposition_id AS voting_main_proposition_id,
				legislative_body.id AS legislative_body_id,
				legislative_body.name AS legislative_body_name,
				legislative_body.acronym AS legislative_body_acronym,
				legislative_body_type.id AS legislative_body_type_id,
				legislative_body_type.description AS legislative_body_type_description
			FROM article
			    INNER JOIN article_type ON article_type.id = article.article_type_id
			    INNER JOIN voting ON voting.article_id = article.id
				INNER JOIN legislative_body ON legislative_body.id = voting.legislative_body_id
				INNER JOIN legislative_body_type ON legislative_body_type.id = legislative_body.legislative_body_type_id
				LEFT JOIN user_article ON user_article.article_id = article.id
			WHERE article.active = true AND article_type.active = true AND voting.active = true AND
				legislative_body.active = true AND legislative_body_type.active = true AND
				user_article.active IS NOT false AND article.id = $1
			GROUP BY article.id, article_type.id, voting.id, legislative_body.id, legislative_body_type.id`
}
