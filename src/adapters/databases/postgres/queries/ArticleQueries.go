package queries

import (
	"fmt"
	"strings"
)

type articleSqlManager struct{}

func Article() *articleSqlManager {
	return &articleSqlManager{}
}

type articleSelectSqlManager struct{}

func (articleSqlManager) Select() *articleSelectSqlManager {
	return &articleSelectSqlManager{}
}

func (articleSelectSqlManager) In(numberOfArticles int) string {
	var parameters []string
	for i := 1; i <= numberOfArticles; i++ {
		parameters = append(parameters, fmt.Sprintf("$%d", i))
	}

	return fmt.Sprintf(`SELECT article.id AS article_id, article.created_at AS article_created_at,
				article.updated_at AS article_updated_at,
				COALESCE(AVG(user_article.rating), 0) AS article_average_rating,
				COUNT(user_article.rating) AS article_number_of_ratings,
				article_type.id AS article_type_id, article_type.description AS article_type_description,
				article_type.codes AS article_type_codes, article_type.color AS article_type_color,
				COALESCE(proposition.id, '00000000-0000-0000-0000-000000000000') AS proposition_id,
				COALESCE(proposition.title, '') AS proposition_title,
				COALESCE(proposition.content, '') AS proposition_content,
				COALESCE(proposition.image_url, '') AS proposition_image_url,
				COALESCE(proposition.image_description, '') AS proposition_image_description,
				COALESCE(proposition_type.id, '00000000-0000-0000-0000-000000000000') AS proposition_type_id,
				COALESCE(proposition_type.description, '') AS proposition_type_description,
				COALESCE(proposition_type.color, '') AS proposition_type_color,
				COALESCE(voting.id, '00000000-0000-0000-0000-000000000000') AS voting_id,
				COALESCE(voting.code, '') AS voting_code, COALESCE(voting.description, '') AS voting_description,
				COALESCE(voting.result, '') AS voting_result,
				COALESCE(voting.result_announced_at, '0001-01-01 00:00:00') AS voting_result_announced_at,
				voting.is_approved AS voting_is_approved,
				COALESCE(event.id, '00000000-0000-0000-0000-000000000000') AS event_id,
				COALESCE(event.title, '') AS event_title, COALESCE(event.description, '') AS event_description,
				COALESCE(event.starts_at, '0001-01-01 00:00:00') AS event_starts_at,
				COALESCE(event.ends_at, '0001-01-01 00:00:00') AS event_ends_at,
				COALESCE(event.video_url, '') AS event_video_url,
				COALESCE(event_type.id, '00000000-0000-0000-0000-000000000000') AS event_type_id,
				COALESCE(event_type.description, '') AS event_type_description,
				COALESCE(event_type.color, '') AS event_type_color,
				COALESCE(event_situation.id, '00000000-0000-0000-0000-000000000000') AS event_situation_id,
				COALESCE(event_situation.description, '') AS event_situation_description,
				COALESCE(event_situation.color, '') AS event_situation_color,
				COALESCE(newsletter.id, '00000000-0000-0000-0000-000000000000') AS newsletter_id,
				COALESCE(newsletter.reference_date, '0001-01-01 00:00:00') AS newsletter_reference_date,
				COALESCE(newsletter.description, '') AS newsletter_description
			FROM article
				INNER JOIN article_type ON article_type.id = article.article_type_id
				LEFT JOIN proposition ON proposition.article_id = article.id
				LEFT JOIN proposition_type ON proposition_type.id = proposition.proposition_type_id
				LEFT JOIN voting ON voting.article_id = article.id
				LEFT JOIN event ON event.article_id = article.id
				LEFT JOIN event_type ON event_type.id = event.event_type_id
				LEFT JOIN event_situation ON event_situation.id = event.event_situation_id
				LEFT JOIN newsletter ON newsletter.article_id = article.id
				LEFT JOIN user_article ON user_article.article_id = article.id
			WHERE article.active = true AND article_type.active = true AND proposition.active IS NOT false AND
				proposition_type.active IS NOT false AND voting.active IS NOT false AND
				event.active IS NOT false AND event_type.active IS NOT false AND
				event_situation.active IS NOT false AND newsletter.active IS NOT false AND
				user_article.active IS NOT false AND article.id IN (%s)
			GROUP BY article.id, article.reference_date_time, article_type.id, proposition.id, proposition_type.id,
				voting.id, event.id, event_type.id, event_situation.id, newsletter.id
			ORDER BY article.reference_date_time DESC`, strings.Join(parameters, ","))
}

func (articleSelectSqlManager) TotalNumberOfArticles() string {
	return `SELECT COUNT(DISTINCT article.id)
			FROM article
				INNER JOIN article_type ON article_type.id = article.article_type_id
				LEFT JOIN proposition ON proposition.article_id = article.id
				LEFT JOIN proposition_type ON proposition_type.id = proposition.proposition_type_id
				LEFT JOIN voting ON voting.article_id = article.id
				LEFT JOIN event ON event.article_id = article.id
				LEFT JOIN event_type ON event_type.id = event.event_type_id
				LEFT JOIN event_situation ON event_situation.id = event.event_situation_id
				LEFT JOIN newsletter ON newsletter.article_id = article.id
			WHERE article.active = true AND article_type.active = true AND proposition.active IS NOT false AND
				proposition_type.active IS NOT false AND voting.active IS NOT false AND
				event.active IS NOT false AND event_type.active IS NOT false AND
				event_situation.active IS NOT false AND newsletter.active IS NOT false AND
				article_type.id = COALESCE($1, article_type.id) AND
				($2::uuid IS NULL OR proposition_type.id = $2 OR event_type.id = $2) AND
				((proposition.title ILIKE $3 OR proposition.content ILIKE $3) OR
				('Votação ' || voting.code ILIKE $3 OR voting.result ILIKE $3) OR
				(event.title ILIKE $3 OR event.description ILIKE $3) OR
				('Boletim do dia ' || TO_CHAR(newsletter.reference_date, 'DD/MM/YYYY') ILIKE $3 OR
				newsletter.description ILIKE $3)) AND
				DATE_TRUNC('day', article.created_at) >= DATE_TRUNC('day', COALESCE($4, article.created_at)) AND
				DATE_TRUNC('day', article.created_at) <= DATE_TRUNC('day', COALESCE($5, article.created_at))`
}

func (articleSelectSqlManager) TotalNumberOfPropositions() string {
	return `SELECT COUNT(DISTINCT article.id)
			FROM article
				INNER JOIN article_type ON article_type.id = article.article_type_id
				INNER JOIN proposition prop ON prop.article_id = article.id
				INNER JOIN proposition_type ON proposition_type.id = prop.proposition_type_id
				INNER JOIN proposition_author ON proposition_author.proposition_id = prop.id
				LEFT JOIN deputy ON deputy.id = proposition_author.deputy_id
				LEFT JOIN party previous_party ON previous_party.id = proposition_author.party_id
				LEFT JOIN external_author ON external_author.id = proposition_author.external_author_id
				LEFT JOIN user_article ON user_article.article_id = article.id
			WHERE article.active = true AND article_type.active = true AND prop.active = true AND
				proposition_type.active = true AND proposition_author.active = true AND deputy.active IS NOT false AND
				previous_party.active IS NOT false AND external_author.active IS NOT false AND
				user_article.active IS NOT false AND article_type.id = COALESCE($1, article_type.id) AND
				proposition_type.id = COALESCE($2, proposition_type.id) AND
				(prop.title ILIKE $3 OR prop.content ILIKE $3) AND
				DATE_TRUNC('day', article.created_at) >= DATE_TRUNC('day', COALESCE($4, article.created_at)) AND
				DATE_TRUNC('day', article.created_at) <= DATE_TRUNC('day', COALESCE($5, article.created_at)) AND
				((deputy.id IS NULL AND ($6::uuid IS NULL AND $8::uuid IS NOT NULL)) OR $6::uuid IS NULL OR
				(SELECT EXISTS(SELECT 1 FROM article
					INNER JOIN proposition prop2 ON prop2.article_id = article.id
					INNER JOIN proposition_author ON proposition_author.proposition_id = prop2.id
				WHERE proposition_author.deputy_id = $6 AND article.id = prop.article_id))) AND
				((previous_party.id IS NULL AND ($7::uuid IS NULL AND $8::uuid IS NOT NULL)) OR $7::uuid IS NULL OR
				(SELECT EXISTS(SELECT 1 FROM article
					INNER JOIN proposition prop2 ON prop2.article_id = article.id
					INNER JOIN proposition_author ON proposition_author.proposition_id = prop2.id
				WHERE proposition_author.party_id = $7 AND article.id = prop.article_id))) AND
				((external_author.id IS NULL AND (($6::uuid IS NOT NULL OR $7::uuid IS NOT NULL) AND $8::uuid IS NULL))
				OR $8::uuid IS NULL OR (SELECT EXISTS(SELECT 1 FROM article
					INNER JOIN proposition prop2 ON prop2.article_id = article.id
					INNER JOIN proposition_author ON proposition_author.proposition_id = prop2.id
				WHERE proposition_author.external_author_id = $8 AND article.id = prop.article_id)))`
}

func (articleSelectSqlManager) TotalNumberOfVotes() string {
	return `SELECT COUNT(DISTINCT article.id)
			FROM article
				INNER JOIN article_type ON article_type.id = article.article_type_id
				INNER JOIN voting ON voting.article_id = article.id
				LEFT JOIN user_article ON user_article.article_id = article.id
			WHERE article.active = true AND article_type.active = true AND voting.active = true AND
				user_article.active IS NOT false AND article_type.id = COALESCE($1, article_type.id) AND
				$2::uuid IS NULL AND ('Votação ' || voting.code ILIKE $3 OR voting.result ILIKE $3) AND
				DATE_TRUNC('day', article.created_at) >= DATE_TRUNC('day', COALESCE($4, article.created_at)) AND
				DATE_TRUNC('day', article.created_at) <= DATE_TRUNC('day', COALESCE($5, article.created_at)) AND
				DATE_TRUNC('day', voting.result_announced_at) >= DATE_TRUNC('day',
				COALESCE($6, voting.result_announced_at)) AND
				DATE_TRUNC('day', voting.result_announced_at) <= DATE_TRUNC('day',
				COALESCE($7, voting.result_announced_at)) AND voting.is_approved = COALESCE($8, voting.is_approved) AND
				voting.legislative_body_id = COALESCE($9, voting.legislative_body_id)`
}

func (articleSelectSqlManager) TotalNumberOfEvents() string {
	return `SELECT COUNT(DISTINCT article.id)
			FROM article
				INNER JOIN article_type ON article_type.id = article.article_type_id
				INNER JOIN event ON event.article_id = article.id
				INNER JOIN event_type ON event_type.id = event.event_type_id
				INNER JOIN event_situation ON event_situation.id = event.event_situation_id
				INNER JOIN event_legislative_body ON event_legislative_body.event_id = event.id
				LEFT JOIN event_agenda_item ON event_agenda_item.event_id = event.id
				LEFT JOIN user_article ON user_article.article_id = article.id
			WHERE article.active = true AND article_type.active = true AND event.active = true AND
				event_type.active = true AND event_situation.active = true AND event_legislative_body.active = true AND
				event_agenda_item.active IS NOT false AND user_article.active IS NOT false AND
				article_type.id = COALESCE($1, article_type.id) AND event_type.id = COALESCE($2, event_type.id) AND
				(event.title ILIKE $3 OR event.description ILIKE $3) AND
				DATE_TRUNC('day', article.created_at) >= DATE_TRUNC('day', COALESCE($4, article.created_at)) AND
				DATE_TRUNC('day', article.created_at) <= DATE_TRUNC('day', COALESCE($5, article.created_at)) AND
				DATE_TRUNC('day', event.ends_at) >= DATE_TRUNC('day', COALESCE($6, event.ends_at)) AND
				DATE_TRUNC('day', event.starts_at) <= DATE_TRUNC('day', COALESCE($7, event.starts_at)) AND
				event_situation.id = COALESCE($8, event_situation.id) AND
				event_legislative_body.legislative_body_id = COALESCE($9, event_legislative_body.legislative_body_id) AND
				($10::uuid IS NULL OR event_agenda_item.rapporteur_id = COALESCE($10, event_agenda_item.rapporteur_id))`
}

func (articleSelectSqlManager) All() string {
	return `SELECT article.id AS article_id, article.created_at AS article_created_at,
				article.updated_at AS article_updated_at,
				COALESCE(AVG(user_article.rating), 0) AS article_average_rating,
				COUNT(user_article.rating) AS article_number_of_ratings,
				article_type.id AS article_type_id, article_type.description AS article_type_description,
				article_type.codes AS article_type_codes, article_type.color AS article_type_color,
				COALESCE(proposition.id, '00000000-0000-0000-0000-000000000000') AS proposition_id,
				COALESCE(proposition.title, '') AS proposition_title,
				COALESCE(proposition.content, '') AS proposition_content,
				COALESCE(proposition.image_url, '') AS proposition_image_url,
				COALESCE(proposition.image_description, '') AS proposition_image_description,
				COALESCE(proposition_type.id, '00000000-0000-0000-0000-000000000000') AS proposition_type_id,
				COALESCE(proposition_type.description, '') AS proposition_type_description,
				COALESCE(proposition_type.color, '') AS proposition_type_color,
				COALESCE(voting.id, '00000000-0000-0000-0000-000000000000') AS voting_id,
				COALESCE(voting.code, '') AS voting_code, COALESCE(voting.description, '') AS voting_description,
				COALESCE(voting.result, '') AS voting_result,
				COALESCE(voting.result_announced_at, '0001-01-01 00:00:00') AS voting_result_announced_at,
				voting.is_approved AS voting_is_approved,
				COALESCE(event.id, '00000000-0000-0000-0000-000000000000') AS event_id,
				COALESCE(event.title, '') AS event_title, COALESCE(event.description, '') AS event_description,
				COALESCE(event.starts_at, '0001-01-01 00:00:00') AS event_starts_at,
				COALESCE(event.ends_at, '0001-01-01 00:00:00') AS event_ends_at,
				COALESCE(event.video_url, '') AS event_video_url,
				COALESCE(event_type.id, '00000000-0000-0000-0000-000000000000') AS event_type_id,
				COALESCE(event_type.description, '') AS event_type_description,
				COALESCE(event_type.color, '') AS event_type_color,
				COALESCE(event_situation.id, '00000000-0000-0000-0000-000000000000') AS event_situation_id,
				COALESCE(event_situation.description, '') AS event_situation_description,
				COALESCE(event_situation.color, '') AS event_situation_color,
				COALESCE(newsletter.id, '00000000-0000-0000-0000-000000000000') AS newsletter_id,
				COALESCE(newsletter.reference_date, '0001-01-01 00:00:00') AS newsletter_reference_date,
				COALESCE(newsletter.description, '') AS newsletter_description
			FROM article
				INNER JOIN article_type ON article_type.id = article.article_type_id
				LEFT JOIN proposition ON proposition.article_id = article.id
				LEFT JOIN proposition_type ON proposition_type.id = proposition.proposition_type_id
				LEFT JOIN voting ON voting.article_id = article.id
				LEFT JOIN event ON event.article_id = article.id
				LEFT JOIN event_type ON event_type.id = event.event_type_id
				LEFT JOIN event_situation ON event_situation.id = event.event_situation_id
				LEFT JOIN newsletter ON newsletter.article_id = article.id
				LEFT JOIN user_article ON user_article.article_id = article.id
			WHERE article.active = true AND article_type.active = true AND proposition.active IS NOT false AND
				proposition_type.active IS NOT false AND voting.active IS NOT false AND
				event.active IS NOT false AND event_type.active IS NOT false AND
				event_situation.active IS NOT false AND newsletter.active IS NOT false AND
				user_article.active IS NOT false AND article_type.id = COALESCE($1, article_type.id) AND
				($2::uuid IS NULL OR proposition_type.id = $2 OR event_type.id = $2) AND
				((proposition.title ILIKE $3 OR proposition.content ILIKE $3) OR
				('Votação ' || voting.code ILIKE $3 OR voting.result ILIKE $3) OR
				(event.title ILIKE $3 OR event.description ILIKE $3) OR
				('Boletim do dia ' || TO_CHAR(newsletter.reference_date, 'DD/MM/YYYY') ILIKE $3 OR
				newsletter.description ILIKE $3)) AND
				DATE_TRUNC('day', article.created_at) >= DATE_TRUNC('day', COALESCE($4, article.created_at)) AND
				DATE_TRUNC('day', article.created_at) <= DATE_TRUNC('day', COALESCE($5, article.created_at))
			GROUP BY article.id, article.reference_date_time, article_type.id, proposition.id, proposition_type.id,
				voting.id, event.id, event_type.id, event_situation.id, newsletter.id
			ORDER BY article.reference_date_time DESC
			OFFSET $6 LIMIT $7`
}

func (articleSelectSqlManager) Propositions() string {
	return `SELECT article.id AS article_id, article.created_at AS article_created_at,
				article.updated_at AS article_updated_at,
				COALESCE(AVG(user_article.rating), 0) AS article_average_rating,
				COUNT(user_article.rating) AS article_number_of_ratings,
				article_type.id AS article_type_id, article_type.description AS article_type_description,
				article_type.codes AS article_type_codes, article_type.color AS article_type_color,
				prop.id AS proposition_id, prop.title AS proposition_title, prop.content AS proposition_content,
				COALESCE(prop.image_url, '') AS proposition_image_url,
				COALESCE(prop.image_description, '') AS proposition_image_description,
				proposition_type.id AS proposition_type_id,
				proposition_type.description AS proposition_type_description,
				proposition_type.color AS proposition_type_color
			FROM article
				INNER JOIN article_type ON article_type.id = article.article_type_id
				INNER JOIN proposition prop ON prop.article_id = article.id
				INNER JOIN proposition_type ON proposition_type.id = prop.proposition_type_id
				INNER JOIN proposition_author ON proposition_author.proposition_id = prop.id
				LEFT JOIN deputy ON deputy.id = proposition_author.deputy_id
				LEFT JOIN party previous_party ON previous_party.id = proposition_author.party_id
				LEFT JOIN external_author ON external_author.id = proposition_author.external_author_id
				LEFT JOIN user_article ON user_article.article_id = article.id
			WHERE article.active = true AND article_type.active = true AND prop.active = true AND
				proposition_type.active = true AND proposition_author.active = true AND deputy.active IS NOT false AND
				previous_party.active IS NOT false AND external_author.active IS NOT false AND
				user_article.active IS NOT false AND article_type.id = COALESCE($1, article_type.id) AND
				proposition_type.id = COALESCE($2, proposition_type.id) AND
				(prop.title ILIKE $3 OR prop.content ILIKE $3) AND
				DATE_TRUNC('day', article.created_at) >= DATE_TRUNC('day', COALESCE($4, article.created_at)) AND
				DATE_TRUNC('day', article.created_at) <= DATE_TRUNC('day', COALESCE($5, article.created_at)) AND
				((deputy.id IS NULL AND ($6::uuid IS NULL AND $8::uuid IS NOT NULL)) OR $6::uuid IS NULL OR
				(SELECT EXISTS(SELECT 1 FROM article
					INNER JOIN proposition prop2 ON prop2.article_id = article.id
					INNER JOIN proposition_author ON proposition_author.proposition_id = prop2.id
				WHERE proposition_author.deputy_id = $6 AND article.id = prop.article_id))) AND
				((previous_party.id IS NULL AND ($7::uuid IS NULL AND $8::uuid IS NOT NULL)) OR $7::uuid IS NULL OR
				(SELECT EXISTS(SELECT 1 FROM article
					INNER JOIN proposition prop2 ON prop2.article_id = article.id
					INNER JOIN proposition_author ON proposition_author.proposition_id = prop2.id
				WHERE proposition_author.party_id = $7 AND article.id = prop.article_id))) AND
				((external_author.id IS NULL AND (($6::uuid IS NOT NULL OR $7::uuid IS NOT NULL) AND $8::uuid IS NULL))
				OR $8::uuid IS NULL OR (SELECT EXISTS(SELECT 1 FROM article
					INNER JOIN proposition prop2 ON prop2.article_id = article.id
					INNER JOIN proposition_author ON proposition_author.proposition_id = prop2.id
				WHERE proposition_author.external_author_id = $8 AND article.id = prop.article_id)))
			GROUP BY article.id, article.reference_date_time, article_type.id, prop.id, proposition_type.id
			ORDER BY article.reference_date_time DESC
			OFFSET $9 LIMIT $10`
}

func (articleSelectSqlManager) Votes() string {
	return `SELECT article.id AS article_id, article.created_at AS article_created_at,
				article.updated_at AS article_updated_at,
				COALESCE(AVG(user_article.rating), 0) AS article_average_rating,
				COUNT(user_article.rating) AS article_number_of_ratings,
				article_type.id AS article_type_id, article_type.description AS article_type_description,
				article_type.codes AS article_type_codes, article_type.color AS article_type_color,
				voting.id AS voting_id, voting.code AS voting_code, voting.result AS voting_result,
				voting.is_approved AS voting_is_approved
			FROM article
				INNER JOIN article_type ON article_type.id = article.article_type_id
				INNER JOIN voting ON voting.article_id = article.id
				LEFT JOIN user_article ON user_article.article_id = article.id
			WHERE article.active = true AND article_type.active = true AND voting.active = true AND
				user_article.active IS NOT false AND article_type.id = COALESCE($1, article_type.id) AND
				$2::uuid IS NULL AND ('Votação ' || voting.code ILIKE $3 OR voting.result ILIKE $3) AND
				DATE_TRUNC('day', article.created_at) >= DATE_TRUNC('day', COALESCE($4, article.created_at)) AND
				DATE_TRUNC('day', article.created_at) <= DATE_TRUNC('day', COALESCE($5, article.created_at)) AND
				DATE_TRUNC('day', voting.result_announced_at) >= DATE_TRUNC('day',
				COALESCE($6, voting.result_announced_at)) AND
				DATE_TRUNC('day', voting.result_announced_at) <= DATE_TRUNC('day',
				COALESCE($7, voting.result_announced_at)) AND voting.is_approved = COALESCE($8, voting.is_approved) AND
				voting.legislative_body_id = COALESCE($9, voting.legislative_body_id)
			GROUP BY article.id, article.reference_date_time, article_type.id, voting.id
			ORDER BY article.reference_date_time DESC
			OFFSET $10 LIMIT $11`
}

func (articleSelectSqlManager) Events() string {
	return `SELECT article.id AS article_id, article.created_at AS article_created_at,
				article.updated_at AS article_updated_at,
				COALESCE(AVG(user_article.rating), 0) AS article_average_rating,
				COUNT(user_article.rating) AS article_number_of_ratings,
				article_type.id AS article_type_id, article_type.description AS article_type_description,
				article_type.codes AS article_type_codes, article_type.color AS article_type_color,
				event.id AS event_id, event.title AS event_title, event.description AS event_description,
				event.starts_at AS event_starts_at,
				COALESCE(event.ends_at, '0001-01-01 00:00:00') AS event_ends_at,
				COALESCE(event.video_url, '') AS event_video_url,
				event_type.id AS event_type_id,
				event_type.description AS event_type_description,
				event_type.color AS event_type_color,
				event_situation.id AS event_situation_id,
				event_situation.description AS event_situation_description,
				event_situation.color AS event_situation_color
			FROM article
				INNER JOIN article_type ON article_type.id = article.article_type_id
				INNER JOIN event ON event.article_id = article.id
				INNER JOIN event_type ON event_type.id = event.event_type_id
				INNER JOIN event_situation ON event_situation.id = event.event_situation_id
				INNER JOIN event_legislative_body ON event_legislative_body.event_id = event.id
				LEFT JOIN event_agenda_item ON event_agenda_item.event_id = event.id
				LEFT JOIN user_article ON user_article.article_id = article.id
			WHERE article.active = true AND article_type.active = true AND event.active = true AND
				event_type.active = true AND event_situation.active = true AND event_legislative_body.active = true AND
				event_agenda_item.active IS NOT false AND user_article.active IS NOT false AND
				article_type.id = COALESCE($1, article_type.id) AND event_type.id = COALESCE($2, event_type.id) AND
				(event.title ILIKE $3 OR event.description ILIKE $3) AND
				DATE_TRUNC('day', article.created_at) >= DATE_TRUNC('day', COALESCE($4, article.created_at)) AND
				DATE_TRUNC('day', article.created_at) <= DATE_TRUNC('day', COALESCE($5, article.created_at)) AND
				DATE_TRUNC('day', event.ends_at) >= DATE_TRUNC('day', COALESCE($6, event.ends_at)) AND
				DATE_TRUNC('day', event.starts_at) <= DATE_TRUNC('day', COALESCE($7, event.starts_at)) AND
				event_situation.id = COALESCE($8, event_situation.id) AND
				event_legislative_body.legislative_body_id = COALESCE($9, event_legislative_body.legislative_body_id) AND
				($10::uuid IS NULL OR event_agenda_item.rapporteur_id = COALESCE($10, event_agenda_item.rapporteur_id))
			GROUP BY article.id, article.reference_date_time, article_type.id, event.id, event_type.id,
				event_situation.id
			ORDER BY article.reference_date_time DESC
			OFFSET $11 LIMIT $12`
}

func (articleSelectSqlManager) TrendingArticles() string {
	return `SELECT article.id AS article_id, COUNT(article_view.id) AS article_views,
				article.created_at AS article_created_at, article.updated_at AS article_updated_at,
				COALESCE(AVG(user_article.rating), 0) AS article_average_rating,
				COUNT(user_article.rating) AS article_number_of_ratings,
				article_type.id AS article_type_id, article_type.description AS article_type_description,
				article_type.codes AS article_type_codes, article_type.color AS article_type_color,
				COALESCE(proposition.id, '00000000-0000-0000-0000-000000000000') AS proposition_id,
				COALESCE(proposition.title, '') AS proposition_title,
				COALESCE(proposition.content, '') AS proposition_content,
				COALESCE(proposition.image_url, '') AS proposition_image_url,
				COALESCE(proposition.image_description, '') AS proposition_image_description,
				COALESCE(proposition_type.id, '00000000-0000-0000-0000-000000000000') AS proposition_type_id,
				COALESCE(proposition_type.description, '') AS proposition_type_description,
				COALESCE(proposition_type.color, '') AS proposition_type_color,
				COALESCE(voting.id, '00000000-0000-0000-0000-000000000000') AS voting_id,
				COALESCE(voting.code, '') AS voting_code, COALESCE(voting.description, '') AS voting_description,
				COALESCE(voting.result, '') AS voting_result,
				COALESCE(voting.result_announced_at, '0001-01-01 00:00:00') AS voting_result_announced_at,
				voting.is_approved AS voting_is_approved,
				COALESCE(event.id, '00000000-0000-0000-0000-000000000000') AS event_id,
				COALESCE(event.title, '') AS event_title, COALESCE(event.description, '') AS event_description,
				COALESCE(event.starts_at, '0001-01-01 00:00:00') AS event_starts_at,
				COALESCE(event.ends_at, '0001-01-01 00:00:00') AS event_ends_at,
				COALESCE(event.video_url, '') AS event_video_url,
				COALESCE(event_type.id, '00000000-0000-0000-0000-000000000000') AS event_type_id,
				COALESCE(event_type.description, '') AS event_type_description,
				COALESCE(event_type.color, '') AS event_type_color,
				COALESCE(event_situation.id, '00000000-0000-0000-0000-000000000000') AS event_situation_id,
				COALESCE(event_situation.description, '') AS event_situation_description,
				COALESCE(event_situation.color, '') AS event_situation_color,
				COALESCE(newsletter.id, '00000000-0000-0000-0000-000000000000') AS newsletter_id,
				COALESCE(newsletter.reference_date, '0001-01-01 00:00:00') AS newsletter_reference_date,
				COALESCE(newsletter.description, '') AS newsletter_description
			FROM article
				INNER JOIN article_type ON article_type.id = article.article_type_id
				LEFT JOIN proposition ON proposition.article_id = article.id
				LEFT JOIN proposition_type ON proposition_type.id = proposition.proposition_type_id
				LEFT JOIN voting ON voting.article_id = article.id
				LEFT JOIN event ON event.article_id = article.id
				LEFT JOIN event_type ON event_type.id = event.event_type_id
				LEFT JOIN event_situation ON event_situation.id = event.event_situation_id
				LEFT JOIN newsletter ON newsletter.article_id = article.id
				LEFT JOIN user_article ON user_article.article_id = article.id
				LEFT JOIN article_view ON article_view.article_id = article.id AND
					article_view.created_at >= CURRENT_DATE - INTERVAL '1 week'
			WHERE article.active = true AND article_type.active = true AND proposition.active IS NOT false AND
				proposition_type.active IS NOT false AND voting.active IS NOT false AND
				event.active IS NOT false AND event_type.active IS NOT false AND
				event_situation.active IS NOT false AND newsletter.active IS NOT false AND
				user_article.active IS NOT false AND article_type.id = COALESCE($1, article_type.id) AND
				($2::uuid IS NULL OR proposition_type.id = $2 OR event_type.id = $2) AND
				((proposition.title ILIKE $3 OR proposition.content ILIKE $3) OR
				('Votação ' || voting.code ILIKE $3 OR voting.result ILIKE $3) OR
				(event.title ILIKE $3 OR event.description ILIKE $3) OR
				('Boletim do dia ' || TO_CHAR(newsletter.reference_date, 'DD/MM/YYYY') ILIKE $3 OR
				newsletter.description ILIKE $3)) AND
				DATE_TRUNC('day', article.created_at) >= DATE_TRUNC('day', COALESCE($4, article.created_at)) AND
				DATE_TRUNC('day', article.created_at) <= DATE_TRUNC('day', COALESCE($5, article.created_at))
			GROUP BY article.id, article.reference_date_time, article_type.id, proposition.id, proposition_type.id,
				voting.id, event.id, event_type.id, event_situation.id, newsletter.id
			ORDER BY article_views DESC, article.reference_date_time DESC
			OFFSET $6 LIMIT $7`
}

func (articleSelectSqlManager) TrendingPropositions() string {
	return `SELECT article.id AS article_id, COUNT(article_view.id) AS article_views,
				article.created_at AS article_created_at, article.updated_at AS article_updated_at,
				COALESCE(AVG(user_article.rating), 0) AS article_average_rating,
				COUNT(user_article.rating) AS article_number_of_ratings,
				article_type.id AS article_type_id, article_type.description AS article_type_description,
				article_type.codes AS article_type_codes, article_type.color AS article_type_color,
				prop.id AS proposition_id, prop.title AS proposition_title, prop.content AS proposition_content,
				COALESCE(prop.image_url, '') AS proposition_image_url,
				COALESCE(prop.image_description, '') AS proposition_image_description,
				proposition_type.id AS proposition_type_id,
				proposition_type.description AS proposition_type_description,
				proposition_type.color AS proposition_type_color
			FROM article
				INNER JOIN article_type ON article_type.id = article.article_type_id
				INNER JOIN proposition prop ON prop.article_id = article.id
				INNER JOIN proposition_type ON proposition_type.id = prop.proposition_type_id
				INNER JOIN proposition_author ON proposition_author.proposition_id = prop.id
				LEFT JOIN deputy ON deputy.id = proposition_author.deputy_id
				LEFT JOIN party previous_party ON previous_party.id = proposition_author.party_id
				LEFT JOIN external_author ON external_author.id = proposition_author.external_author_id
				LEFT JOIN user_article ON user_article.article_id = article.id
				LEFT JOIN article_view ON article_view.article_id = article.id AND
					article_view.created_at >= CURRENT_DATE - INTERVAL '1 week'
			WHERE article.active = true AND article_type.active = true AND prop.active = true AND
				proposition_type.active = true AND proposition_author.active = true AND deputy.active IS NOT false AND
				previous_party.active IS NOT false AND external_author.active IS NOT false AND
				user_article.active IS NOT false AND article_type.id = COALESCE($1, article_type.id) AND
				proposition_type.id = COALESCE($2, proposition_type.id) AND
				(prop.title ILIKE $3 OR prop.content ILIKE $3) AND
				DATE_TRUNC('day', article.created_at) >= DATE_TRUNC('day', COALESCE($4, article.created_at)) AND
				DATE_TRUNC('day', article.created_at) <= DATE_TRUNC('day', COALESCE($5, article.created_at)) AND
				((deputy.id IS NULL AND ($6::uuid IS NULL AND $8::uuid IS NOT NULL)) OR $6::uuid IS NULL OR
				(SELECT EXISTS(SELECT 1 FROM article
					INNER JOIN proposition prop2 ON prop2.article_id = article.id
					INNER JOIN proposition_author ON proposition_author.proposition_id = prop2.id
				WHERE proposition_author.deputy_id = $6 AND article.id = prop.article_id))) AND
				((previous_party.id IS NULL AND ($7::uuid IS NULL AND $8::uuid IS NOT NULL)) OR $7::uuid IS NULL OR
				(SELECT EXISTS(SELECT 1 FROM article
					INNER JOIN proposition prop2 ON prop2.article_id = article.id
					INNER JOIN proposition_author ON proposition_author.proposition_id = prop2.id
				WHERE proposition_author.party_id = $7 AND article.id = prop.article_id))) AND
				((external_author.id IS NULL AND (($6::uuid IS NOT NULL OR $7::uuid IS NOT NULL) AND $8::uuid IS NULL))
				OR $8::uuid IS NULL OR (SELECT EXISTS(SELECT 1 FROM article
					INNER JOIN proposition prop2 ON prop2.article_id = article.id
					INNER JOIN proposition_author ON proposition_author.proposition_id = prop2.id
				WHERE proposition_author.external_author_id = $8 AND article.id = prop.article_id)))
			GROUP BY article.id, article.reference_date_time, article_type.id, prop.id, proposition_type.id
			ORDER BY article_views DESC, article.reference_date_time DESC
			OFFSET $9 LIMIT $10`
}

func (articleSelectSqlManager) TrendingVotes() string {
	return `SELECT article.id AS article_id, COUNT(article_view.id) AS article_views,
				article.created_at AS article_created_at, article.updated_at AS article_updated_at,
				COALESCE(AVG(user_article.rating), 0) AS article_average_rating,
				COUNT(user_article.rating) AS article_number_of_ratings,
				article_type.id AS article_type_id, article_type.description AS article_type_description,
				article_type.codes AS article_type_codes, article_type.color AS article_type_color,
				voting.id AS voting_id, voting.code AS voting_code, voting.result AS voting_result,
				voting.is_approved AS voting_is_approved
			FROM article
				INNER JOIN article_type ON article_type.id = article.article_type_id
				INNER JOIN voting ON voting.article_id = article.id
				LEFT JOIN user_article ON user_article.article_id = article.id
				LEFT JOIN article_view ON article_view.article_id = article.id AND
					article_view.created_at >= CURRENT_DATE - INTERVAL '1 week'
			WHERE article.active = true AND article_type.active = true AND voting.active = true AND
				user_article.active IS NOT false AND article_type.id = COALESCE($1, article_type.id) AND
				$2::uuid IS NULL AND ('Votação ' || voting.code ILIKE $3 OR voting.result ILIKE $3) AND
				DATE_TRUNC('day', article.created_at) >= DATE_TRUNC('day', COALESCE($4, article.created_at)) AND
				DATE_TRUNC('day', article.created_at) <= DATE_TRUNC('day', COALESCE($5, article.created_at)) AND
				DATE_TRUNC('day', voting.result_announced_at) >= DATE_TRUNC('day',
				COALESCE($6, voting.result_announced_at)) AND
				DATE_TRUNC('day', voting.result_announced_at) <= DATE_TRUNC('day',
				COALESCE($7, voting.result_announced_at)) AND voting.is_approved = COALESCE($8, voting.is_approved) AND
				voting.legislative_body_id = COALESCE($9, voting.legislative_body_id)
			GROUP BY article.id, article.reference_date_time, article_type.id, voting.id
			ORDER BY article_views DESC, article.reference_date_time DESC
			OFFSET $10 LIMIT $11`
}

func (articleSelectSqlManager) TrendingEvents() string {
	return `SELECT article.id AS article_id, COUNT(article_view.id) AS article_views,
				article.created_at AS article_created_at, article.updated_at AS article_updated_at,
				COALESCE(AVG(user_article.rating), 0) AS article_average_rating,
				COUNT(user_article.rating) AS article_number_of_ratings,
				article_type.id AS article_type_id, article_type.description AS article_type_description,
				article_type.codes AS article_type_codes, article_type.color AS article_type_color,
				event.id AS event_id, event.title AS event_title, event.description AS event_description,
				event.starts_at AS event_starts_at,
				COALESCE(event.ends_at, '0001-01-01 00:00:00') AS event_ends_at,
				COALESCE(event.video_url, '') AS event_video_url,
				event_type.id AS event_type_id,
				event_type.description AS event_type_description,
				event_type.color AS event_type_color,
				event_situation.id AS event_situation_id,
				event_situation.description AS event_situation_description,
				event_situation.color AS event_situation_color
			FROM article
				INNER JOIN article_type ON article_type.id = article.article_type_id
				INNER JOIN event ON event.article_id = article.id
				INNER JOIN event_type ON event_type.id = event.event_type_id
				INNER JOIN event_situation ON event_situation.id = event.event_situation_id
				INNER JOIN event_legislative_body ON event_legislative_body.event_id = event.id
				LEFT JOIN event_agenda_item ON event_agenda_item.event_id = event.id
				LEFT JOIN user_article ON user_article.article_id = article.id
				LEFT JOIN article_view ON article_view.article_id = article.id AND
					article_view.created_at >= CURRENT_DATE - INTERVAL '1 week'
			WHERE article.active = true AND article_type.active = true AND event.active = true AND
				event_type.active = true AND event_situation.active = true AND event_legislative_body.active = true AND
				event_agenda_item.active IS NOT false AND user_article.active IS NOT false AND
				article_type.id = COALESCE($1, article_type.id) AND event_type.id = COALESCE($2, event_type.id) AND
				(event.title ILIKE $3 OR event.description ILIKE $3) AND
				DATE_TRUNC('day', article.created_at) >= DATE_TRUNC('day', COALESCE($4, article.created_at)) AND
				DATE_TRUNC('day', article.created_at) <= DATE_TRUNC('day', COALESCE($5, article.created_at)) AND
				DATE_TRUNC('day', event.ends_at) >= DATE_TRUNC('day', COALESCE($6, event.ends_at)) AND
				DATE_TRUNC('day', event.starts_at) <= DATE_TRUNC('day', COALESCE($7, event.starts_at)) AND
				event_situation.id = COALESCE($8, event_situation.id) AND
				event_legislative_body.legislative_body_id = COALESCE($9, event_legislative_body.legislative_body_id) AND
				($10::uuid IS NULL OR event_agenda_item.rapporteur_id = COALESCE($10, event_agenda_item.rapporteur_id))
			GROUP BY article.id, article.reference_date_time, article_type.id, event.id, event_type.id,
				event_situation.id
			ORDER BY article_views DESC, article.reference_date_time DESC
			OFFSET $11 LIMIT $12`
}

func (articleSelectSqlManager) TrendingArticlesByTypeId() string {
	return `SELECT article.id AS article_id, COUNT(article_view.id) AS article_views,
				article.created_at AS article_created_at, article.updated_at AS article_updated_at,
				COALESCE(AVG(user_article.rating), 0) AS article_average_rating,
				COUNT(user_article.rating) AS article_number_of_ratings,
				article_type.id AS article_type_id, article_type.description AS article_type_description,
				article_type.codes AS article_type_codes, article_type.color AS article_type_color,
				COALESCE(proposition.id, '00000000-0000-0000-0000-000000000000') AS proposition_id,
				COALESCE(proposition.title, '') AS proposition_title,
				COALESCE(proposition.content, '') AS proposition_content,
				COALESCE(proposition.image_url, '') AS proposition_image_url,
				COALESCE(proposition.image_description, '') AS proposition_image_description,
				COALESCE(proposition_type.id, '00000000-0000-0000-0000-000000000000') AS proposition_type_id,
				COALESCE(proposition_type.description, '') AS proposition_type_description,
				COALESCE(proposition_type.color, '') AS proposition_type_color,
				COALESCE(voting.id, '00000000-0000-0000-0000-000000000000') AS voting_id,
				COALESCE(voting.code, '') AS voting_code, COALESCE(voting.description, '') AS voting_description,
				COALESCE(voting.result, '') AS voting_result,
				COALESCE(voting.result_announced_at, '0001-01-01 00:00:00') AS voting_result_announced_at,
				voting.is_approved AS voting_is_approved,
				COALESCE(event.id, '00000000-0000-0000-0000-000000000000') AS event_id,
				COALESCE(event.title, '') AS event_title, COALESCE(event.description, '') AS event_description,
				COALESCE(event.starts_at, '0001-01-01 00:00:00') AS event_starts_at,
				COALESCE(event.ends_at, '0001-01-01 00:00:00') AS event_ends_at,
				COALESCE(event.video_url, '') AS event_video_url,
				COALESCE(event_type.id, '00000000-0000-0000-0000-000000000000') AS event_type_id,
				COALESCE(event_type.description, '') AS event_type_description,
				COALESCE(event_type.color, '') AS event_type_color,
				COALESCE(event_situation.id, '00000000-0000-0000-0000-000000000000') AS event_situation_id,
				COALESCE(event_situation.description, '') AS event_situation_description,
				COALESCE(event_situation.color, '') AS event_situation_color,
				COALESCE(newsletter.id, '00000000-0000-0000-0000-000000000000') AS newsletter_id,
				COALESCE(newsletter.reference_date, '0001-01-01 00:00:00') AS newsletter_reference_date,
				COALESCE(newsletter.description, '') AS newsletter_description
			FROM article
				INNER JOIN article_type ON article_type.id = article.article_type_id
				LEFT JOIN proposition ON proposition.article_id = article.id
				LEFT JOIN proposition_type ON proposition_type.id = proposition.proposition_type_id
				LEFT JOIN voting ON voting.article_id = article.id
				LEFT JOIN event ON event.article_id = article.id
				LEFT JOIN event_type ON event_type.id = event.event_type_id
				LEFT JOIN event_situation ON event_situation.id = event.event_situation_id
				LEFT JOIN newsletter ON newsletter.article_id = article.id
				LEFT JOIN user_article ON user_article.article_id = article.id
				LEFT JOIN article_view ON article_view.article_id = article.id AND
					article_view.created_at >= CURRENT_DATE - INTERVAL '1 week'
			WHERE article.active = true AND article_type.active = true AND proposition.active IS NOT false AND
				proposition_type.active IS NOT false AND voting.active IS NOT false AND
				event.active IS NOT false AND event_type.active IS NOT false AND
				event_situation.active IS NOT false AND newsletter.active IS NOT false AND
				user_article.active IS NOT false AND article_type.id = $1
			GROUP BY article.id, article.reference_date_time, article_type.id, proposition.id, proposition_type.id,
				voting.id, event.id, event_type.id, event_situation.id, newsletter.id
			ORDER BY article_views DESC, article.reference_date_time DESC
			LIMIT $2`
}

func (articleSelectSqlManager) TrendingArticlesBySpecificTypeId() string {
	return `SELECT article.id AS article_id, COUNT(article_view.id) AS article_views,
				article.created_at AS article_created_at, article.updated_at AS article_updated_at,
				COALESCE(AVG(user_article.rating), 0) AS article_average_rating,
				COUNT(user_article.rating) AS article_number_of_ratings,
				article_type.id AS article_type_id, article_type.description AS article_type_description,
				article_type.codes AS article_type_codes, article_type.color AS article_type_color,
				COALESCE(proposition.id, '00000000-0000-0000-0000-000000000000') AS proposition_id,
				COALESCE(proposition.title, '') AS proposition_title,
				COALESCE(proposition.content, '') AS proposition_content,
				COALESCE(proposition.image_url, '') AS proposition_image_url,
				COALESCE(proposition.image_description, '') AS proposition_image_description,
				COALESCE(proposition_type.id, '00000000-0000-0000-0000-000000000000') AS proposition_type_id,
				COALESCE(proposition_type.description, '') AS proposition_type_description,
				COALESCE(proposition_type.color, '') AS proposition_type_color,
				COALESCE(event.id, '00000000-0000-0000-0000-000000000000') AS event_id,
				COALESCE(event.title, '') AS event_title, COALESCE(event.description, '') AS event_description,
				COALESCE(event.starts_at, '0001-01-01 00:00:00') AS event_starts_at,
				COALESCE(event.ends_at, '0001-01-01 00:00:00') AS event_ends_at,
				COALESCE(event.video_url, '') AS event_video_url,
				COALESCE(event_type.id, '00000000-0000-0000-0000-000000000000') AS event_type_id,
				COALESCE(event_type.description, '') AS event_type_description,
				COALESCE(event_type.color, '') AS event_type_color,
				COALESCE(event_situation.id, '00000000-0000-0000-0000-000000000000') AS event_situation_id,
				COALESCE(event_situation.description, '') AS event_situation_description,
				COALESCE(event_situation.color, '') AS event_situation_color
			FROM article
				INNER JOIN article_type ON article_type.id = article.article_type_id
				LEFT JOIN proposition ON proposition.article_id = article.id
				LEFT JOIN proposition_type ON proposition_type.id = proposition.proposition_type_id
				LEFT JOIN event ON event.article_id = article.id
				LEFT JOIN event_type ON event_type.id = event.event_type_id
				LEFT JOIN event_situation ON event_situation.id = event.event_situation_id
				LEFT JOIN user_article ON user_article.article_id = article.id
				LEFT JOIN article_view ON article_view.article_id = article.id AND
					article_view.created_at >= CURRENT_DATE - INTERVAL '1 week'
			WHERE article.active = true AND article_type.active = true AND proposition.active IS NOT false AND
				proposition_type.active IS NOT false AND event.active IS NOT false AND
				event_type.active IS NOT false AND event_situation.active IS NOT false AND
				user_article.active IS NOT false AND (proposition_type.id = $1 OR event_type.id = $1)
			GROUP BY article.id, article.reference_date_time, article_type.id, proposition.id, proposition_type.id,
				event.id, event_type.id, event_situation.id
			ORDER BY article_views DESC, article.reference_date_time DESC
			LIMIT $2`
}

func (articleSelectSqlManager) RelatedArticlesByPropositionId() string {
	return `WITH related_votes AS (
				SELECT voting_id FROM proposition_related_to_voting WHERE active = true AND proposition_id = $1
				UNION
				SELECT voting_id FROM proposition_affected_by_voting WHERE active = true AND proposition_id = $1),
			related_events AS (
				SELECT event_id FROM event_requirement WHERE active = true AND proposition_id = $1
				UNION
				SELECT event_id FROM event_agenda_item WHERE active = true AND (proposition_id = $1 OR
					related_proposition_id = $1))
			SELECT article.id AS article_id
			FROM article
				INNER JOIN voting ON voting.article_id = article.id
				LEFT JOIN related_votes ON voting.id = related_votes.voting_id
			WHERE voting.main_proposition_id = $1 OR related_votes.voting_id IS NOT NULL
			UNION
			SELECT article.id AS article_id
			FROM article
				INNER JOIN event ON event.article_id = article.id
				INNER JOIN related_events ON related_events.event_id = event.id
			UNION
			SELECT article.id AS article_id
			FROM article
				INNER JOIN newsletter ON newsletter.article_id = article.id
				INNER JOIN newsletter_article ON newsletter_article.newsletter_id = newsletter.id
				INNER JOIN proposition ON proposition.article_id = newsletter_article.article_id
			WHERE newsletter_article.active = true AND proposition.id = $1`
}

func (articleSelectSqlManager) RelatedArticlesByVotingId() string {
	return `SELECT article.id AS article_id
			FROM article
				INNER JOIN event ON event.article_id = article.id
				INNER JOIN event_agenda_item ON event_agenda_item.event_id = event.id
			WHERE event_agenda_item.active = true AND event_agenda_item.voting_id = $1
			UNION
			SELECT article.id AS article_id
			FROM article
				INNER JOIN newsletter ON newsletter.article_id = article.id
				INNER JOIN newsletter_article ON newsletter_article.newsletter_id = newsletter.id
				INNER JOIN voting ON voting.article_id = newsletter_article.article_id
			WHERE newsletter_article.active = true AND voting.id = $1`
}

func (articleSelectSqlManager) NewsletterArticleByArticleId() string {
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
				INNER JOIN newsletter_article ON newsletter_article.newsletter_id = newsletter.id
				LEFT JOIN user_article ON user_article.article_id = article.id
			WHERE article.active = true AND article_type.active = true AND newsletter.active = true AND
				newsletter_article.active = true AND user_article.active IS NOT false AND
				newsletter_article.article_id = $1
			GROUP BY article.id, article_type.id, newsletter.id`
}

func (articleSelectSqlManager) ArticlesByNewsletterId() string {
	return `SELECT article.id AS article_id, article.created_at AS article_created_at,
				article.updated_at AS article_updated_at,
				COALESCE(AVG(user_article.rating), 0) AS article_average_rating,
				COUNT(user_article.rating) AS article_number_of_ratings,
				article_type.id AS article_type_id, article_type.description AS article_type_description,
				article_type.codes AS article_type_codes, article_type.color AS article_type_color,
				COALESCE(proposition.id, '00000000-0000-0000-0000-000000000000') AS proposition_id,
				COALESCE(proposition.title, '') AS proposition_title,
				COALESCE(proposition.content, '') AS proposition_content,
				COALESCE(proposition.image_url, '') AS proposition_image_url,
				COALESCE(proposition.image_description, '') AS proposition_image_description,
				COALESCE(proposition_type.id, '00000000-0000-0000-0000-000000000000') AS proposition_type_id,
				COALESCE(proposition_type.description, '') AS proposition_type_description,
				COALESCE(proposition_type.color, '') AS proposition_type_color,
				COALESCE(voting.id, '00000000-0000-0000-0000-000000000000') AS voting_id,
				COALESCE(voting.code, '') AS voting_code, COALESCE(voting.description, '') AS voting_description,
				COALESCE(voting.result, '') AS voting_result,
				COALESCE(voting.result_announced_at, '0001-01-01 00:00:00') AS voting_result_announced_at,
				voting.is_approved AS voting_is_approved,
				COALESCE(event.id, '00000000-0000-0000-0000-000000000000') AS event_id,
				COALESCE(event.title, '') AS event_title, COALESCE(event.description, '') AS event_description,
				COALESCE(event.starts_at, '0001-01-01 00:00:00') AS event_starts_at,
				COALESCE(event.ends_at, '0001-01-01 00:00:00') AS event_ends_at,
				COALESCE(event.video_url, '') AS event_video_url,
				COALESCE(event_type.id, '00000000-0000-0000-0000-000000000000') AS event_type_id,
				COALESCE(event_type.description, '') AS event_type_description,
				COALESCE(event_type.color, '') AS event_type_color,
				COALESCE(event_situation.id, '00000000-0000-0000-0000-000000000000') AS event_situation_id,
				COALESCE(event_situation.description, '') AS event_situation_description,
				COALESCE(event_situation.color, '') AS event_situation_color
			FROM article
				INNER JOIN newsletter_article ON newsletter_article.article_id = article.id
				INNER JOIN article_type ON article_type.id = article.article_type_id
				LEFT JOIN proposition ON proposition.article_id = article.id
				LEFT JOIN proposition_type ON proposition_type.id = proposition.proposition_type_id
				LEFT JOIN voting ON voting.article_id = article.id
				LEFT JOIN event ON event.article_id = article.id
				LEFT JOIN event_type ON event_type.id = event.event_type_id
				LEFT JOIN event_situation ON event_situation.id = event.event_situation_id
				LEFT JOIN user_article ON user_article.article_id = article.id
			WHERE article.active = true AND article_type.active = true AND proposition.active IS NOT false AND
				proposition_type.active IS NOT false AND voting.active IS NOT false AND
				event.active IS NOT false AND event_type.active IS NOT false AND
				event_situation.active IS NOT false AND user_article.active IS NOT false AND
				newsletter_article.newsletter_id = $1
			GROUP BY article.id, article.reference_date_time, article_type.id, proposition.id, proposition_type.id,
				voting.id, event.id, event_type.id, event_situation.id
			ORDER BY article.reference_date_time`
}

func (articleSelectSqlManager) MainPropositionByVotingId() string {
	return `SELECT article.id AS article_id, article.created_at AS article_created_at,
				article.updated_at AS article_updated_at,
				COALESCE(AVG(user_article.rating), 0) AS article_average_rating,
				COUNT(user_article.rating) AS article_number_of_ratings,
				article_type.id AS article_type_id, article_type.description AS article_type_description,
				article_type.codes AS article_type_codes, article_type.color AS article_type_color,
				proposition.id AS proposition_id,
				proposition.title AS proposition_title,
				proposition.content AS proposition_content,
				COALESCE(proposition.image_url, '') AS proposition_image_url,
				COALESCE(proposition.image_description, '') AS proposition_image_description,
				proposition_type.id AS proposition_type_id,
				proposition_type.description AS proposition_type_description,
				proposition_type.color AS proposition_type_color
			FROM article
				INNER JOIN article_type ON article_type.id = article.article_type_id
				INNER JOIN proposition ON proposition.article_id = article.id
				INNER JOIN proposition_type ON proposition_type.id = proposition.proposition_type_id
				INNER JOIN voting ON voting.main_proposition_id = proposition.id
				LEFT JOIN user_article ON user_article.article_id = article.id
			WHERE article.active = true AND article_type.active = true AND proposition.active IS NOT false AND
			      proposition_type.active IS NOT false AND voting.active IS NOT false AND
			      user_article.active IS NOT false AND voting.id = $1
			GROUP BY article.id, article_type.id, proposition.id, proposition_type.id, voting.id`
}

func (articleSelectSqlManager) PropositionsRelatedByVotingId() string {
	return `SELECT article.id AS article_id, article.created_at AS article_created_at,
				article.updated_at AS article_updated_at,
				COALESCE(AVG(user_article.rating), 0) AS article_average_rating,
				COUNT(user_article.rating) AS article_number_of_ratings,
				article_type.id AS article_type_id, article_type.description AS article_type_description,
				article_type.codes AS article_type_codes, article_type.color AS article_type_color,
				proposition.id AS proposition_id,
				proposition.title AS proposition_title,
				proposition.content AS proposition_content,
				COALESCE(proposition.image_url, '') AS proposition_image_url,
				COALESCE(proposition.image_description, '') AS proposition_image_description,
				proposition_type.id AS proposition_type_id,
				proposition_type.description AS proposition_type_description,
				proposition_type.color AS proposition_type_color
			FROM article
				INNER JOIN article_type ON article_type.id = article.article_type_id
				INNER JOIN proposition ON proposition.article_id = article.id
				INNER JOIN proposition_type ON proposition_type.id = proposition.proposition_type_id
				INNER JOIN proposition_related_to_voting ON proposition_related_to_voting.proposition_id = proposition.id
				LEFT JOIN user_article ON user_article.article_id = article.id
			WHERE article.active = true AND article_type.active = true AND proposition.active = true AND
				proposition_type.active = true AND proposition_related_to_voting.active = true AND
				user_article.active IS NOT false AND proposition_related_to_voting.voting_id = $1
			GROUP BY article.id, article.reference_date_time, article_type.id, proposition.id, proposition_type.id
			ORDER BY article.reference_date_time`
}

func (articleSelectSqlManager) PropositionsAffectedByVotingId() string {
	return `SELECT article.id AS article_id, article.created_at AS article_created_at,
				article.updated_at AS article_updated_at,
				COALESCE(AVG(user_article.rating), 0) AS article_average_rating,
				COUNT(user_article.rating) AS article_number_of_ratings,
				article_type.id AS article_type_id, article_type.description AS article_type_description,
				article_type.codes AS article_type_codes, article_type.color AS article_type_color,
				proposition.id AS proposition_id,
				proposition.title AS proposition_title,
				proposition.content AS proposition_content,
				COALESCE(proposition.image_url, '') AS proposition_image_url,
				COALESCE(proposition.image_description, '') AS proposition_image_description,
				proposition_type.id AS proposition_type_id,
				proposition_type.description AS proposition_type_description,
				proposition_type.color AS proposition_type_color
			FROM article
				INNER JOIN article_type ON article_type.id = article.article_type_id
				INNER JOIN proposition ON proposition.article_id = article.id
				INNER JOIN proposition_type ON proposition_type.id = proposition.proposition_type_id
				INNER JOIN proposition_affected_by_voting ON proposition_affected_by_voting.proposition_id = proposition.id
				LEFT JOIN user_article ON user_article.article_id = article.id
			WHERE article.active = true AND article_type.active = true AND proposition.active = true AND
				proposition_type.active = true AND proposition_affected_by_voting.active = true AND
				user_article.active IS NOT false AND proposition_affected_by_voting.voting_id = $1
			GROUP BY article.id, article.reference_date_time, article_type.id, proposition.id, proposition_type.id
			ORDER BY article.reference_date_time`
}

func (articleSelectSqlManager) PropositionsOfTheRequirementsByEventId() string {
	return `SELECT article.id AS article_id, article.created_at AS article_created_at,
				article.updated_at AS article_updated_at,
				COALESCE(AVG(user_article.rating), 0) AS article_average_rating,
				COUNT(user_article.rating) AS article_number_of_ratings,
				article_type.id AS article_type_id, article_type.description AS article_type_description,
				article_type.codes AS article_type_codes, article_type.color AS article_type_color,
				proposition.id AS proposition_id,
				proposition.title AS proposition_title,
				proposition.content AS proposition_content,
				COALESCE(proposition.image_url, '') AS proposition_image_url,
				COALESCE(proposition.image_description, '') AS proposition_image_description,
				proposition_type.id AS proposition_type_id,
				proposition_type.description AS proposition_type_description,
				proposition_type.color AS proposition_type_color
			FROM article
				INNER JOIN article_type ON article_type.id = article.article_type_id
				INNER JOIN proposition ON proposition.article_id = article.id
				INNER JOIN proposition_type ON proposition_type.id = proposition.proposition_type_id
				INNER JOIN event_requirement ON event_requirement.proposition_id = proposition.id
				LEFT JOIN user_article ON user_article.article_id = article.id
			WHERE article.active = true AND article_type.active = true AND proposition.active = true AND
				proposition_type.active = true AND event_requirement.active = true AND
				user_article.active IS NOT false AND event_requirement.event_id = $1
			GROUP BY article.id, article.reference_date_time, article_type.id, proposition.id, proposition_type.id
			ORDER BY article.reference_date_time`
}

func (articleSelectSqlManager) RatingsAndArticlesSavedForLaterViewing(numberOfArticles int) string {
	var parameters []string
	for i := 1; i <= numberOfArticles; i++ {
		parameters = append(parameters, fmt.Sprintf("$%d", i+1))
	}

	return fmt.Sprintf(`SELECT article.id AS article_id,
							COALESCE(user_article.rating, 0) AS user_article_rating,
							COALESCE(user_article.view_later, false) AS user_article_view_later
						FROM article
							LEFT JOIN user_article ON user_article.article_id = article.id
						WHERE article.active = true AND user_article.user_id = $1 AND article.id IN (%s)
						GROUP BY article.id, article.reference_date_time, user_article.rating, user_article.view_later
						ORDER BY article.reference_date_time DESC`, strings.Join(parameters, ","))
}

func (articleSelectSqlManager) NumberOfArticlesBookmarkedToViewLater() string {
	return `SELECT COUNT(DISTINCT article.id)
			FROM article
				INNER JOIN article_type ON article_type.id = article.article_type_id
				LEFT JOIN proposition ON proposition.article_id = article.id
				LEFT JOIN proposition_type ON proposition_type.id = proposition.proposition_type_id
				LEFT JOIN voting ON voting.article_id = article.id
				LEFT JOIN event ON event.article_id = article.id
				LEFT JOIN event_type ON event_type.id = event.event_type_id
				LEFT JOIN event_situation ON event_situation.id = event.event_situation_id
				LEFT JOIN newsletter ON newsletter.article_id = article.id
				INNER JOIN user_article ON user_article.article_id = article.id
			WHERE article.active = true AND article_type.active = true AND proposition.active IS NOT false AND
				proposition_type.active IS NOT false AND voting.active IS NOT false AND
				event.active IS NOT false AND event_type.active IS NOT false AND
				event_situation.active IS NOT false AND newsletter.active IS NOT false AND user_article.active = true AND
				article_type.id = COALESCE($1, article_type.id) AND ($2::uuid IS NULL OR proposition_type.id = $2 OR
				event_type.id = $2) AND ((proposition.title ILIKE $3 OR proposition.content ILIKE $3) OR
				('Votação ' || voting.code ILIKE $3 OR voting.result ILIKE $3) OR
				(event.title ILIKE $3 OR event.description ILIKE $3) OR
				('Boletim do dia ' || TO_CHAR(newsletter.reference_date, 'DD/MM/YYYY') ILIKE $3 OR
				newsletter.description ILIKE $3)) AND
				DATE_TRUNC('day', article.created_at) >= DATE_TRUNC('day', COALESCE($4, article.created_at)) AND
				DATE_TRUNC('day', article.created_at) <= DATE_TRUNC('day', COALESCE($5, article.created_at)) AND
				user_article.user_id = $6 AND user_article.view_later = true`
}

func (articleSelectSqlManager) NumberOfPropositionsBookmarkedToViewLater() string {
	return `SELECT COUNT(DISTINCT article.id)
			FROM article
				INNER JOIN article_type ON article_type.id = article.article_type_id
				INNER JOIN proposition prop ON prop.article_id = article.id
				INNER JOIN proposition_type ON proposition_type.id = prop.proposition_type_id
				INNER JOIN proposition_author ON proposition_author.proposition_id = prop.id
				LEFT JOIN deputy ON deputy.id = proposition_author.deputy_id
				LEFT JOIN party previous_party ON previous_party.id = proposition_author.party_id
				LEFT JOIN external_author ON external_author.id = proposition_author.external_author_id
				LEFT JOIN user_article ON user_article.article_id = article.id
			WHERE article.active = true AND article_type.active = true AND prop.active = true AND
				proposition_type.active = true AND proposition_author.active = true AND deputy.active IS NOT false AND
				previous_party.active IS NOT false AND external_author.active IS NOT false AND
				user_article.active IS NOT false AND article_type.id = COALESCE($1, article_type.id) AND
				proposition_type.id = COALESCE($2, proposition_type.id) AND
				(prop.title ILIKE $3 OR prop.content ILIKE $3) AND
				DATE_TRUNC('day', article.created_at) >= DATE_TRUNC('day', COALESCE($4, article.created_at)) AND
				DATE_TRUNC('day', article.created_at) <= DATE_TRUNC('day', COALESCE($5, article.created_at)) AND
				((deputy.id IS NULL AND ($6::uuid IS NULL AND $8::uuid IS NOT NULL)) OR $6::uuid IS NULL OR
				(SELECT EXISTS(SELECT 1 FROM article
					INNER JOIN proposition prop2 ON prop2.article_id = article.id
					INNER JOIN proposition_author ON proposition_author.proposition_id = prop2.id
				WHERE proposition_author.deputy_id = $6 AND article.id = prop.article_id))) AND
				((previous_party.id IS NULL AND ($7::uuid IS NULL AND $8::uuid IS NOT NULL)) OR $7::uuid IS NULL OR
				(SELECT EXISTS(SELECT 1 FROM article
					INNER JOIN proposition prop2 ON prop2.article_id = article.id
					INNER JOIN proposition_author ON proposition_author.proposition_id = prop2.id
				WHERE proposition_author.party_id = $7 AND article.id = prop.article_id))) AND
				((external_author.id IS NULL AND (($6::uuid IS NOT NULL OR $7::uuid IS NOT NULL) AND $8::uuid IS NULL))
				OR $8::uuid IS NULL OR (SELECT EXISTS(SELECT 1 FROM article
					INNER JOIN proposition prop2 ON prop2.article_id = article.id
					INNER JOIN proposition_author ON proposition_author.proposition_id = prop2.id
				WHERE proposition_author.external_author_id = $8 AND article.id = prop.article_id))) AND
				user_article.user_id = $9 AND user_article.view_later = true`
}

func (articleSelectSqlManager) NumberOfVotesBookmarkedToViewLater() string {
	return `SELECT COUNT(DISTINCT article.id)
			FROM article
				INNER JOIN article_type ON article_type.id = article.article_type_id
				INNER JOIN voting ON voting.article_id = article.id
				LEFT JOIN user_article ON user_article.article_id = article.id
			WHERE article.active = true AND article_type.active = true AND voting.active = true AND
				user_article.active IS NOT false AND article_type.id = COALESCE($1, article_type.id) AND
				$2::uuid IS NULL AND ('Votação ' || voting.code ILIKE $3 OR voting.result ILIKE $3) AND
				DATE_TRUNC('day', article.created_at) >= DATE_TRUNC('day', COALESCE($4, article.created_at)) AND
				DATE_TRUNC('day', article.created_at) <= DATE_TRUNC('day', COALESCE($5, article.created_at)) AND
				DATE_TRUNC('day', voting.result_announced_at) >= DATE_TRUNC('day',
				COALESCE($6, voting.result_announced_at)) AND
				DATE_TRUNC('day', voting.result_announced_at) <= DATE_TRUNC('day',
				COALESCE($7, voting.result_announced_at)) AND voting.is_approved = COALESCE($8, voting.is_approved) AND
				voting.legislative_body_id = COALESCE($9, voting.legislative_body_id) AND
				user_article.user_id = $10 AND user_article.view_later = true`
}

func (articleSelectSqlManager) NumberOfEventsBookmarkedToViewLater() string {
	return `SELECT COUNT(DISTINCT article.id)
			FROM article
				INNER JOIN article_type ON article_type.id = article.article_type_id
				INNER JOIN event ON event.article_id = article.id
				INNER JOIN event_type ON event_type.id = event.event_type_id
				INNER JOIN event_situation ON event_situation.id = event.event_situation_id
				INNER JOIN event_legislative_body ON event_legislative_body.event_id = event.id
				LEFT JOIN event_agenda_item ON event_agenda_item.event_id = event.id
				LEFT JOIN user_article ON user_article.article_id = article.id
			WHERE article.active = true AND article_type.active = true AND event.active = true AND
				event_type.active = true AND event_situation.active = true AND event_legislative_body.active = true AND
				event_agenda_item.active IS NOT false AND user_article.active IS NOT false AND
				article_type.id = COALESCE($1, article_type.id) AND event_type.id = COALESCE($2, event_type.id) AND
				(event.title ILIKE $3 OR event.description ILIKE $3) AND
				DATE_TRUNC('day', article.created_at) >= DATE_TRUNC('day', COALESCE($4, article.created_at)) AND
				DATE_TRUNC('day', article.created_at) <= DATE_TRUNC('day', COALESCE($5, article.created_at)) AND
				DATE_TRUNC('day', event.ends_at) >= DATE_TRUNC('day', COALESCE($6, event.ends_at)) AND
				DATE_TRUNC('day', event.starts_at) <= DATE_TRUNC('day', COALESCE($7, event.starts_at)) AND
				event_situation.id = COALESCE($8, event_situation.id) AND
				event_legislative_body.legislative_body_id = COALESCE($9, event_legislative_body.legislative_body_id) AND
				($10::uuid IS NULL OR event_agenda_item.rapporteur_id = COALESCE($10, event_agenda_item.rapporteur_id)) AND
				user_article.user_id = $11 AND user_article.view_later = true`
}

func (articleSelectSqlManager) ArticlesBookmarkedToViewLater() string {
	return `SELECT article.id AS article_id, COALESCE(user_article.rating, 0) AS user_article_rating,
				COALESCE(user_article.view_later, false) AS user_article_view_later
			FROM article
				INNER JOIN article_type ON article_type.id = article.article_type_id
				LEFT JOIN proposition ON proposition.article_id = article.id
				LEFT JOIN proposition_type ON proposition_type.id = proposition.proposition_type_id
				LEFT JOIN voting ON voting.article_id = article.id
				LEFT JOIN event ON event.article_id = article.id
				LEFT JOIN event_type ON event_type.id = event.event_type_id
				LEFT JOIN event_situation ON event_situation.id = event.event_situation_id
				LEFT JOIN newsletter ON newsletter.article_id = article.id
				INNER JOIN user_article ON user_article.article_id = article.id
			WHERE article.active = true AND article_type.active = true AND proposition.active IS NOT false AND
				proposition_type.active IS NOT false AND voting.active IS NOT false AND
				event.active IS NOT false AND event_type.active IS NOT false AND event_situation.active IS NOT false AND
				newsletter.active IS NOT false AND user_article.active = true AND
				article_type.id = COALESCE($1, article_type.id) AND ($2::uuid IS NULL OR proposition_type.id = $2 OR
				event_type.id = $2) AND ((proposition.title ILIKE $3 OR proposition.content ILIKE $3) OR
				('Votação ' || voting.code ILIKE $3 OR voting.result ILIKE $3) OR
				(event.title ILIKE $3 OR event.description ILIKE $3) OR
				('Boletim do dia ' || TO_CHAR(newsletter.reference_date, 'DD/MM/YYYY') ILIKE $3 OR
				newsletter.description ILIKE $3)) AND
				DATE_TRUNC('day', article.created_at) >= DATE_TRUNC('day', COALESCE($4, article.created_at)) AND
				DATE_TRUNC('day', article.created_at) <= DATE_TRUNC('day', COALESCE($5, article.created_at)) AND
				user_article.user_id = $6 AND user_article.view_later = true
			GROUP BY article.id, user_article.rating, user_article.view_later, user_article.view_later_set_at
			ORDER BY user_article.view_later_set_at DESC
			OFFSET $7 LIMIT $8`
}

func (articleSelectSqlManager) PropositionsBookmarkedToViewLater() string {
	return `SELECT article.id AS article_id, COALESCE(user_article.rating, 0) AS user_article_rating,
				COALESCE(user_article.view_later, false) AS user_article_view_later
			FROM article
				INNER JOIN article_type ON article_type.id = article.article_type_id
				INNER JOIN proposition prop ON prop.article_id = article.id
				INNER JOIN proposition_type ON proposition_type.id = prop.proposition_type_id
				INNER JOIN proposition_author ON proposition_author.proposition_id = prop.id
				LEFT JOIN deputy ON deputy.id = proposition_author.deputy_id
				LEFT JOIN party previous_party ON previous_party.id = proposition_author.party_id
				LEFT JOIN external_author ON external_author.id = proposition_author.external_author_id
				LEFT JOIN user_article ON user_article.article_id = article.id
			WHERE article.active = true AND article_type.active = true AND prop.active = true AND
				proposition_type.active = true AND proposition_author.active = true AND deputy.active IS NOT false AND
				previous_party.active IS NOT false AND external_author.active IS NOT false AND
				user_article.active IS NOT false AND article_type.id = COALESCE($1, article_type.id) AND
				proposition_type.id = COALESCE($2, proposition_type.id) AND
				(prop.title ILIKE $3 OR prop.content ILIKE $3) AND
				DATE_TRUNC('day', article.created_at) >= DATE_TRUNC('day', COALESCE($4, article.created_at)) AND
				DATE_TRUNC('day', article.created_at) <= DATE_TRUNC('day', COALESCE($5, article.created_at)) AND
				((deputy.id IS NULL AND ($6::uuid IS NULL AND $8::uuid IS NOT NULL)) OR $6::uuid IS NULL OR
				(SELECT EXISTS(SELECT 1 FROM article
					INNER JOIN proposition prop2 ON prop2.article_id = article.id
					INNER JOIN proposition_author ON proposition_author.proposition_id = prop2.id
				WHERE proposition_author.deputy_id = $6 AND article.id = prop.article_id))) AND
				((previous_party.id IS NULL AND ($7::uuid IS NULL AND $8::uuid IS NOT NULL)) OR $7::uuid IS NULL OR
				(SELECT EXISTS(SELECT 1 FROM article
					INNER JOIN proposition prop2 ON prop2.article_id = article.id
					INNER JOIN proposition_author ON proposition_author.proposition_id = prop2.id
				WHERE proposition_author.party_id = $7 AND article.id = prop.article_id))) AND
				((external_author.id IS NULL AND (($6::uuid IS NOT NULL OR $7::uuid IS NOT NULL) AND $8::uuid IS NULL))
				OR $8::uuid IS NULL OR (SELECT EXISTS(SELECT 1 FROM article
					INNER JOIN proposition prop2 ON prop2.article_id = article.id
					INNER JOIN proposition_author ON proposition_author.proposition_id = prop2.id
				WHERE proposition_author.external_author_id = $8 AND article.id = prop.article_id))) AND
				user_article.user_id = $9 AND user_article.view_later = true
			GROUP BY article.id, user_article.rating, user_article.view_later, user_article.view_later_set_at
			ORDER BY user_article.view_later_set_at DESC
			OFFSET $10 LIMIT $11`
}

func (articleSelectSqlManager) VotesBookmarkedToViewLater() string {
	return `SELECT article.id AS article_id, COALESCE(user_article.rating, 0) AS user_article_rating,
				COALESCE(user_article.view_later, false) AS user_article_view_later
			FROM article
				INNER JOIN article_type ON article_type.id = article.article_type_id
				INNER JOIN voting ON voting.article_id = article.id
				LEFT JOIN user_article ON user_article.article_id = article.id
			WHERE article.active = true AND article_type.active = true AND voting.active = true AND
				user_article.active IS NOT false AND article_type.id = COALESCE($1, article_type.id) AND
				$2::uuid IS NULL AND ('Votação ' || voting.code ILIKE $3 OR voting.result ILIKE $3) AND
				DATE_TRUNC('day', article.created_at) >= DATE_TRUNC('day', COALESCE($4, article.created_at)) AND
				DATE_TRUNC('day', article.created_at) <= DATE_TRUNC('day', COALESCE($5, article.created_at)) AND
				DATE_TRUNC('day', voting.result_announced_at) >= DATE_TRUNC('day',
				COALESCE($6, voting.result_announced_at)) AND
				DATE_TRUNC('day', voting.result_announced_at) <= DATE_TRUNC('day',
				COALESCE($7, voting.result_announced_at)) AND voting.is_approved = COALESCE($8, voting.is_approved) AND
				voting.legislative_body_id = COALESCE($9, voting.legislative_body_id) AND
				user_article.user_id = $10 AND user_article.view_later = true
			GROUP BY article.id, user_article.rating, user_article.view_later, user_article.view_later_set_at
			ORDER BY user_article.view_later_set_at DESC
			OFFSET $11 LIMIT $12`
}

func (articleSelectSqlManager) EventsBookmarkedToViewLater() string {
	return `SELECT article.id AS article_id, COALESCE(user_article.rating, 0) AS user_article_rating,
				COALESCE(user_article.view_later, false) AS user_article_view_later
			FROM article
				INNER JOIN article_type ON article_type.id = article.article_type_id
				INNER JOIN event ON event.article_id = article.id
				INNER JOIN event_type ON event_type.id = event.event_type_id
				INNER JOIN event_situation ON event_situation.id = event.event_situation_id
				INNER JOIN event_legislative_body ON event_legislative_body.event_id = event.id
				LEFT JOIN event_agenda_item ON event_agenda_item.event_id = event.id
				LEFT JOIN user_article ON user_article.article_id = article.id
			WHERE article.active = true AND article_type.active = true AND event.active = true AND
				event_type.active = true AND event_situation.active = true AND event_legislative_body.active = true AND
				event_agenda_item.active IS NOT false AND user_article.active IS NOT false AND
				article_type.id = COALESCE($1, article_type.id) AND event_type.id = COALESCE($2, event_type.id) AND
				(event.title ILIKE $3 OR event.description ILIKE $3) AND
				DATE_TRUNC('day', article.created_at) >= DATE_TRUNC('day', COALESCE($4, article.created_at)) AND
				DATE_TRUNC('day', article.created_at) <= DATE_TRUNC('day', COALESCE($5, article.created_at)) AND
				DATE_TRUNC('day', event.ends_at) >= DATE_TRUNC('day', COALESCE($6, event.ends_at)) AND
				DATE_TRUNC('day', event.starts_at) <= DATE_TRUNC('day', COALESCE($7, event.starts_at)) AND
				event_situation.id = COALESCE($8, event_situation.id) AND
				event_legislative_body.legislative_body_id = COALESCE($9, event_legislative_body.legislative_body_id) AND
				($10::uuid IS NULL OR event_agenda_item.rapporteur_id = COALESCE($10, event_agenda_item.rapporteur_id)) AND
				user_article.user_id = $11 AND user_article.view_later = true
			GROUP BY article.id, user_article.rating, user_article.view_later, user_article.view_later_set_at
			ORDER BY user_article.view_later_set_at DESC
			OFFSET $12 LIMIT $13`
}
