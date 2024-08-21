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

	return fmt.Sprintf(`SELECT article.id AS article_id, article.reference_date_time AS article_reference_date_time,
							article.created_at AS article_created_at, article.updated_at AS article_updated_at,
							COALESCE(AVG(user_article.rating), 0) AS article_average_rating,
							COUNT(user_article.rating) AS article_number_of_ratings,
							COALESCE(proposition.id, '00000000-0000-0000-0000-000000000000') AS proposition_id,
							COALESCE(proposition.original_text_url, '') AS proposition_original_text_url,
							COALESCE(proposition.title, '') AS proposition_title,
							COALESCE(proposition.content, '') AS proposition_content,
							COALESCE(proposition.submitted_at, '1970-01-01 00:00:00') AS proposition_submitted_at,
							COALESCE(proposition.created_at, '1970-01-01 00:00:00') AS proposition_created_at,
							COALESCE(proposition.updated_at, '1970-01-01 00:00:00') AS proposition_updated_at,
							COALESCE(newsletter.id, '00000000-0000-0000-0000-000000000000') AS newsletter_id,
							COALESCE(newsletter.reference_date, '1970-01-01 00:00:00') AS newsletter_reference_date,
							COALESCE(newsletter.title, '') AS newsletter_title,
							COALESCE(newsletter.description, '') AS newsletter_description,
							COALESCE(newsletter.created_at, '1970-01-01 00:00:00') AS newsletter_created_at,
							COALESCE(newsletter.updated_at, '1970-01-01 00:00:00') AS newsletter_updated_at
						FROM article
							LEFT JOIN proposition ON proposition.id = article.proposition_id
							LEFT JOIN newsletter ON newsletter.id = article.newsletter_id
							LEFT JOIN user_article ON user_article.article_id = article.id
						WHERE article.id IN (%s)
						GROUP BY article.id, article.reference_date_time, proposition.id, newsletter.id
						ORDER BY article.reference_date_time DESC`, strings.Join(parameters, ","))
}

func (articleSelectSqlManager) TotalNumberOfArticles() string {
	return `SELECT COUNT(DISTINCT article.id)
			FROM article
			    LEFT JOIN proposition ON proposition.id = article.proposition_id
			    LEFT JOIN newsletter ON newsletter.id = article.newsletter_id
			WHERE article.active = true AND (proposition.active = true OR newsletter.active = true) AND
				((proposition.title ILIKE $1 OR proposition.content ILIKE $1) OR
				(newsletter.title ILIKE $1 OR newsletter.description ILIKE $1)) AND
				DATE(article.created_at) >= DATE(COALESCE($2, article.created_at)) AND
				DATE(article.created_at) <= DATE(COALESCE($3, article.created_at))`
}

func (articleSelectSqlManager) TotalNumberOfPropositions() string {
	return `SELECT COUNT(DISTINCT article.id)
			FROM article
			    LEFT JOIN proposition prop ON prop.id = article.proposition_id
			    LEFT JOIN proposition_author ON proposition_author.proposition_id = prop.id
			    LEFT JOIN deputy ON deputy.id = proposition_author.deputy_id
				LEFT JOIN party party_in_the_proposal ON party_in_the_proposal.id = proposition_author.party_id
				LEFT JOIN external_author ON external_author.id = proposition_author.external_author_id
			WHERE article.active = true AND prop.active = true AND proposition_author.active = true AND
				(prop.title ILIKE $1 OR prop.content ILIKE $1) AND
				((deputy.id IS NULL AND ($2::uuid IS NULL AND $4::uuid IS NOT NULL))
	    			OR $2::uuid IS NULL
					OR (SELECT EXISTS(SELECT 1 FROM article
							LEFT JOIN proposition prop2 ON prop2.id = article.proposition_id
							LEFT JOIN proposition_author ON proposition_author.proposition_id = prop2.id
						WHERE proposition_author.deputy_id = $2 AND article.proposition_id = prop.id))) AND
				((party_in_the_proposal.id IS NULL AND ($3::uuid IS NULL AND $4::uuid IS NOT NULL))
					OR $3::uuid IS NULL
					OR (SELECT EXISTS(SELECT 1 FROM article
							LEFT JOIN proposition prop2 ON prop2.id = article.proposition_id
							LEFT JOIN proposition_author ON proposition_author.proposition_id = prop2.id
	                	WHERE proposition_author.party_id = $3 AND article.proposition_id = prop.id))) AND
				((external_author.id IS NULL AND (($2::uuid IS NOT NULL OR $3::uuid IS NOT NULL) AND $4::uuid IS NULL))
					OR $4::uuid IS NULL
					OR (SELECT EXISTS(SELECT 1 FROM article
							LEFT JOIN proposition prop2 ON prop2.id = article.proposition_id
							LEFT JOIN proposition_author ON proposition_author.proposition_id = prop2.id
					WHERE proposition_author.external_author_id = $4 AND article.proposition_id = prop.id))) AND
			    article.newsletter_id IS NULL AND
				DATE(article.created_at) >= DATE(COALESCE($5, article.created_at)) AND
				DATE(article.created_at) <= DATE(COALESCE($6, article.created_at))`
}

func (articleSelectSqlManager) TotalNumberOfNewsletters() string {
	return `SELECT COUNT(*)
			FROM article
				LEFT JOIN newsletter ON newsletter.id = article.newsletter_id
    		WHERE article.active = true AND newsletter.active = true AND
    			(newsletter.title ILIKE $1 OR newsletter.description ILIKE $1) AND
    			DATE(article.created_at) >= DATE(COALESCE($2, article.created_at)) AND
				DATE(article.created_at) <= DATE(COALESCE($3, article.created_at))`
}

func (articleSelectSqlManager) All() string {
	return `SELECT article.id AS article_id, article.reference_date_time AS article_reference_date_time,
				article.created_at AS article_created_at, article.updated_at AS article_updated_at,
				COALESCE(AVG(user_article.rating), 0) AS article_average_rating,
				COUNT(user_article.rating) AS article_number_of_ratings,
    			COALESCE(proposition.id, '00000000-0000-0000-0000-000000000000') AS proposition_id,
				COALESCE(proposition.original_text_url, '') AS proposition_original_text_url,
	   			COALESCE(proposition.title, '') AS proposition_title,
				COALESCE(proposition.content, '') AS proposition_content,
				COALESCE(proposition.submitted_at, '1970-01-01 00:00:00') AS proposition_submitted_at,
				COALESCE(proposition.created_at, '1970-01-01 00:00:00') AS proposition_created_at,
				COALESCE(proposition.updated_at, '1970-01-01 00:00:00') AS proposition_updated_at,
				COALESCE(newsletter.id, '00000000-0000-0000-0000-000000000000') AS newsletter_id,
				COALESCE(newsletter.reference_date, '1970-01-01 00:00:00') AS newsletter_reference_date,
				COALESCE(newsletter.title, '') AS newsletter_title,
				COALESCE(newsletter.description, '') AS newsletter_description,
				COALESCE(newsletter.created_at, '1970-01-01 00:00:00') AS newsletter_created_at,
				COALESCE(newsletter.updated_at, '1970-01-01 00:00:00') AS newsletter_updated_at
			FROM article
			    LEFT JOIN proposition ON proposition.id = article.proposition_id
			    LEFT JOIN newsletter ON newsletter.id = article.newsletter_id
				LEFT JOIN user_article ON user_article.article_id = article.id
			WHERE article.active = true AND (proposition.active = true OR newsletter.active = true) AND
				((proposition.title ILIKE $1 OR proposition.content ILIKE $1) OR
				(newsletter.title ILIKE $1 OR newsletter.description ILIKE $1)) AND
				DATE(article.created_at) >= DATE(COALESCE($2, article.created_at)) AND
				DATE(article.created_at) <= DATE(COALESCE($3, article.created_at))
			GROUP BY article.id, article.reference_date_time, proposition.id, newsletter.id
			ORDER BY article.reference_date_time DESC OFFSET $4 LIMIT $5`
}

func (articleSelectSqlManager) Propositions() string {
	return `SELECT article.id AS article_id, article.reference_date_time AS article_reference_date_time,
				article.created_at AS article_created_at, article.updated_at AS article_updated_at,
				COALESCE(AVG(user_article.rating), 0) AS article_average_rating,
				COUNT(user_article.rating) AS article_number_of_ratings,
				prop.id AS proposition_id, prop.original_text_url AS proposition_original_text_url,
				prop.title AS proposition_title, prop.content AS proposition_content,
				prop.submitted_at AS proposition_submitted_at, prop.created_at AS proposition_created_at,
				prop.updated_at AS proposition_updated_at
			FROM article
			    LEFT JOIN proposition prop ON prop.id = article.proposition_id
			    LEFT JOIN proposition_author ON proposition_author.proposition_id = prop.id
			    LEFT JOIN deputy ON deputy.id = proposition_author.deputy_id
				LEFT JOIN party party_in_the_proposal ON party_in_the_proposal.id = proposition_author.party_id
				LEFT JOIN external_author ON external_author.id = proposition_author.external_author_id
				LEFT JOIN user_article ON user_article.article_id = article.id
			WHERE article.active = true AND prop.active = true AND proposition_author.active = true AND
				(prop.title ILIKE $1 OR prop.content ILIKE $1) AND
				((deputy.id IS NULL AND ($2::uuid IS NULL AND $4::uuid IS NOT NULL))
	    			OR $2::uuid IS NULL
					OR (SELECT EXISTS(SELECT 1 FROM article
							LEFT JOIN proposition prop2 ON prop2.id = article.proposition_id
							LEFT JOIN proposition_author ON proposition_author.proposition_id = prop2.id
						WHERE proposition_author.deputy_id = $2 AND article.proposition_id = prop.id))) AND
				((party_in_the_proposal.id IS NULL AND ($3::uuid IS NULL AND $4::uuid IS NOT NULL))
					OR $3::uuid IS NULL
					OR (SELECT EXISTS(SELECT 1 FROM article
							LEFT JOIN proposition prop2 ON prop2.id = article.proposition_id
							LEFT JOIN proposition_author ON proposition_author.proposition_id = prop2.id
	                	WHERE proposition_author.party_id = $3 AND article.proposition_id = prop.id))) AND
				((external_author.id IS NULL AND (($2::uuid IS NOT NULL OR $3::uuid IS NOT NULL) AND $4::uuid IS NULL))
					OR $4::uuid IS NULL
					OR (SELECT EXISTS(SELECT 1 FROM article
							LEFT JOIN proposition prop2 ON prop2.id = article.proposition_id
							LEFT JOIN proposition_author ON proposition_author.proposition_id = prop2.id
					WHERE proposition_author.external_author_id = $4 AND article.proposition_id = prop.id))) AND
			    article.newsletter_id IS NULL AND
				DATE(article.created_at) >= DATE(COALESCE($5, article.created_at)) AND
				DATE(article.created_at) <= DATE(COALESCE($6, article.created_at))
			GROUP BY article.id, article.reference_date_time, prop.id
			ORDER BY article.reference_date_time DESC OFFSET $7 LIMIT $8`
}

func (articleSelectSqlManager) Newsletters() string {
	return `SELECT article.id AS article_id, article.reference_date_time AS article_reference_date_time,
				article.created_at AS article_created_at, article.updated_at AS article_updated_at,
				COALESCE(AVG(user_article.rating), 0) AS article_average_rating,
				COUNT(user_article.rating) AS article_number_of_ratings,
				newsletter.id AS newsletter_id, newsletter.reference_date AS newsletter_reference_date,
				newsletter.title AS newsletter_title, newsletter.description AS newsletter_description,
				newsletter.created_at AS newsletter_created_at, newsletter.updated_at AS newsletter_updated_at
			FROM article
				LEFT JOIN newsletter ON newsletter.id = article.newsletter_id
    			LEFT JOIN user_article ON user_article.article_id = article.id
			WHERE article.active = true AND newsletter.active = true AND
				(newsletter.title ILIKE $1 OR newsletter.description ILIKE $1) AND
				DATE(article.created_at) >= DATE(COALESCE($2, article.created_at)) AND
				DATE(article.created_at) <= DATE(COALESCE($3, article.created_at))
			GROUP BY article.id, article.reference_date_time, newsletter.id
			ORDER BY article.reference_date_time DESC OFFSET $4 LIMIT $5`
}

func (articleSelectSqlManager) TrendingArticles() string {
	return `SELECT article.id AS article_id, article.reference_date_time AS article_reference_date_time,
				article.created_at AS article_created_at, article.updated_at AS article_updated_at,
				COALESCE(AVG(user_article.rating), 0) AS article_average_rating,
				COUNT(user_article.rating) AS article_number_of_ratings,
				COUNT(article_view.id) AS article_views,
       			COALESCE(proposition.id, '00000000-0000-0000-0000-000000000000') AS proposition_id,
				COALESCE(proposition.original_text_url, '') AS proposition_original_text_url,
				COALESCE(proposition.title, '') AS proposition_title,
				COALESCE(proposition.content, '') AS proposition_content,
				COALESCE(proposition.submitted_at, '1970-01-01 00:00:00') AS proposition_submitted_at,
				COALESCE(proposition.created_at, '1970-01-01 00:00:00') AS proposition_created_at,
				COALESCE(proposition.updated_at, '1970-01-01 00:00:00') AS proposition_updated_at,
				COALESCE(newsletter.id, '00000000-0000-0000-0000-000000000000') AS newsletter_id,
				COALESCE(newsletter.reference_date, '1970-01-01 00:00:00') AS newsletter_reference_date,
				COALESCE(newsletter.title, '') AS newsletter_title,
				COALESCE(newsletter.description, '') AS newsletter_description,
				COALESCE(newsletter.created_at, '1970-01-01 00:00:00') AS newsletter_created_at,
				COALESCE(newsletter.updated_at, '1970-01-01 00:00:00') AS newsletter_updated_at
			FROM article
				LEFT JOIN proposition ON proposition.id = article.proposition_id
				LEFT JOIN newsletter ON newsletter.id = article.newsletter_id
			    LEFT JOIN user_article ON user_article.article_id = article.id
				LEFT JOIN article_view ON article_view.article_id = article.id AND article_view.created_at >= CURRENT_DATE - INTERVAL '1 week'
			WHERE article.active = true AND (proposition.active = true OR newsletter.active = true) AND
    			((proposition.title ILIKE $1 OR proposition.content ILIKE $1) OR
    			(newsletter.title ILIKE $1 OR newsletter.description ILIKE $1)) AND
				DATE(article.created_at) >= DATE(COALESCE($2, article.created_at)) AND
				DATE(article.created_at) <= DATE(COALESCE($3, article.created_at))
			GROUP BY article.id, article.reference_date_time, proposition.id, newsletter.id
			ORDER BY article_views DESC, article.reference_date_time DESC OFFSET $4 LIMIT $5`
}

func (articleSelectSqlManager) TrendingPropositions() string {
	return `SELECT article.id AS article_id, article.reference_date_time AS article_reference_date_time,
				article.created_at AS article_created_at, article.updated_at AS article_updated_at,
				COALESCE(AVG(user_article.rating), 0) AS article_average_rating,
				COUNT(user_article.rating) AS article_number_of_ratings,
				COUNT(article_view.id) AS article_views,
				prop.id AS proposition_id, prop.original_text_url AS proposition_original_text_url,
				prop.title AS proposition_title, prop.content AS proposition_content,
				prop.submitted_at AS proposition_submitted_at, prop.created_at AS proposition_created_at,
				prop.updated_at AS proposition_updated_at
			FROM article
				LEFT JOIN proposition prop ON prop.id = article.proposition_id
			    LEFT JOIN proposition_author ON proposition_author.proposition_id = prop.id
			    LEFT JOIN deputy ON deputy.id = proposition_author.deputy_id
				LEFT JOIN party party_in_the_proposal ON party_in_the_proposal.id = proposition_author.party_id
				LEFT JOIN external_author ON external_author.id = proposition_author.external_author_id
			    LEFT JOIN user_article ON user_article.article_id = article.id
				LEFT JOIN article_view ON article_view.article_id = article.id AND article_view.created_at >= CURRENT_DATE - INTERVAL '1 week'
			WHERE article.active = true AND prop.active = true AND proposition_author.active = true AND
				(prop.title ILIKE $1 OR prop.content ILIKE $1) AND
				((deputy.id IS NULL AND ($2::uuid IS NULL AND $4::uuid IS NOT NULL))
	    			OR $2::uuid IS NULL
					OR (SELECT EXISTS(SELECT 1 FROM article
							LEFT JOIN proposition prop2 ON prop2.id = article.proposition_id
							LEFT JOIN proposition_author ON proposition_author.proposition_id = prop2.id
						WHERE proposition_author.deputy_id = $2 AND article.proposition_id = prop.id))) AND
				((party_in_the_proposal.id IS NULL AND ($3::uuid IS NULL AND $4::uuid IS NOT NULL))
					OR $3::uuid IS NULL
					OR (SELECT EXISTS(SELECT 1 FROM article
							LEFT JOIN proposition prop2 ON prop2.id = article.proposition_id
							LEFT JOIN proposition_author ON proposition_author.proposition_id = prop2.id
	                	WHERE proposition_author.party_id = $3 AND article.proposition_id = prop.id))) AND
				((external_author.id IS NULL AND (($2::uuid IS NOT NULL OR $3::uuid IS NOT NULL) AND $4::uuid IS NULL))
					OR $4::uuid IS NULL
					OR (SELECT EXISTS(SELECT 1 FROM article
							LEFT JOIN proposition prop2 ON prop2.id = article.proposition_id
							LEFT JOIN proposition_author ON proposition_author.proposition_id = prop2.id
					WHERE proposition_author.external_author_id = $4 AND article.proposition_id = prop.id))) AND
			    article.newsletter_id IS NULL AND
				DATE(article.created_at) >= DATE(COALESCE($5, article.created_at)) AND
				DATE(article.created_at) <= DATE(COALESCE($6, article.created_at))
			GROUP BY article.id, article.reference_date_time, prop.id
			ORDER BY article_views DESC, article.reference_date_time DESC OFFSET $7 LIMIT $8`
}

func (articleSelectSqlManager) TrendingNewsletters() string {
	return `SELECT article.id AS article_id, article.reference_date_time AS article_reference_date_time,
				article.created_at AS article_created_at, article.updated_at AS article_updated_at,
				COALESCE(AVG(user_article.rating), 0) AS article_average_rating,
				COUNT(user_article.rating) AS article_number_of_ratings,
				COUNT(article_view.id) AS article_views,
				newsletter.id AS newsletter_id, newsletter.reference_date AS newsletter_reference_date,
				newsletter.title AS newsletter_title, newsletter.description AS newsletter_description,
				newsletter.created_at AS newsletter_created_at, newsletter.updated_at AS newsletter_updated_at
			FROM article
				LEFT JOIN newsletter ON newsletter.id = article.newsletter_id
			    LEFT JOIN user_article ON user_article.article_id = article.id
				LEFT JOIN article_view ON article_view.article_id = article.id AND article_view.created_at >= CURRENT_DATE - INTERVAL '1 week'
			WHERE article.active = true AND newsletter.active = true AND
				(newsletter.title ILIKE $1 OR newsletter.description ILIKE $1) AND
    			DATE(article.created_at) >= DATE(COALESCE($2, article.created_at)) AND
				DATE(article.created_at) <= DATE(COALESCE($3, article.created_at))
			GROUP BY article.id, article.reference_date_time, newsletter.id
			ORDER BY article_views DESC, article.reference_date_time DESC OFFSET $4 LIMIT $5`
}

func (articleSelectSqlManager) TrendingArticlesByPropositionType() string {
	return `SELECT article.id AS article_id, article.reference_date_time AS article_reference_date_time,
				article.created_at AS article_created_at, article.updated_at AS article_updated_at,
				COALESCE(AVG(user_article.rating), 0) AS article_average_rating,
				COUNT(user_article.rating) AS article_number_of_ratings,
				COUNT(article_view.id) AS article_views,
				proposition.id AS proposition_id, proposition.original_text_url AS proposition_original_text_url,
				proposition.title AS proposition_title, proposition.content AS proposition_content,
				proposition.submitted_at AS proposition_submitted_at, proposition.created_at AS proposition_created_at,
				proposition.updated_at AS proposition_updated_at
			FROM article
				LEFT JOIN proposition ON proposition.id = article.proposition_id
				LEFT JOIN user_article ON user_article.article_id = article.id
				LEFT JOIN proposition_type ON proposition_type.id = proposition.proposition_type_id
				LEFT JOIN article_view ON article_view.article_id = article.id AND article_view.created_at >= CURRENT_DATE - INTERVAL '1 week'
			WHERE article.active = true AND proposition.active = true AND proposition_type.active = true AND proposition_type.id = $1
			GROUP BY article.id, article.reference_date_time, proposition.id
			ORDER BY article_views DESC, article.reference_date_time DESC LIMIT $2`
}

func (articleSelectSqlManager) NewsletterByPropositionId() string {
	return `SELECT article.id AS article_id, article.reference_date_time AS article_reference_date_time,
				article.created_at AS article_created_at, article.updated_at AS article_updated_at,
				COALESCE(AVG(user_article.rating), 0) AS article_average_rating,
				COUNT(user_article.rating) AS article_number_of_ratings,
				newsletter.id AS newsletter_id, newsletter.reference_date AS newsletter_reference_date,
				newsletter.title AS newsletter_title, newsletter.description AS newsletter_description,
				newsletter.created_at AS newsletter_created_at, newsletter.updated_at AS newsletter_updated_at
			FROM article
				LEFT JOIN newsletter ON newsletter.id = article.newsletter_id
				LEFT JOIN newsletter_proposition ON newsletter.id = newsletter_proposition.newsletter_id
				LEFT JOIN proposition ON proposition.id = newsletter_proposition.proposition_id
				LEFT JOIN user_article ON user_article.article_id = article.id
			WHERE article.active = true AND newsletter.active = true AND newsletter_proposition.active = true AND
			proposition.active = true AND proposition.id = $1
			GROUP BY article.id, newsletter.id`
}

func (articleSelectSqlManager) PropositionsByNewsletterId() string {
	return `SELECT article.id AS article_id, article.reference_date_time AS article_reference_date_time,
				article.created_at AS article_created_at, article.updated_at AS article_updated_at,
				COALESCE(AVG(user_article.rating), 0) AS article_average_rating,
				COUNT(user_article.rating) AS article_number_of_ratings,
				proposition.id AS proposition_id, proposition.original_text_url AS proposition_original_text_url,
				proposition.title AS proposition_title, proposition.content AS proposition_content,
				proposition.submitted_at AS proposition_submitted_at, proposition.created_at AS proposition_created_at,
				proposition.updated_at AS proposition_updated_at
			FROM article
				LEFT JOIN proposition ON proposition.id = article.proposition_id
				LEFT JOIN newsletter_proposition ON proposition.id = newsletter_proposition.proposition_id
			    LEFT JOIN newsletter ON newsletter.id = newsletter_proposition.newsletter_id
				LEFT JOIN user_article ON user_article.article_id = article.id
			WHERE article.active = true AND proposition.active = true AND newsletter_proposition.active = true AND
			      newsletter.active = true AND newsletter_proposition.newsletter_id = $1
			GROUP BY article.id, article.reference_date_time, proposition.id
			ORDER BY article.reference_date_time`
}
