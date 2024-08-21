package queries

import (
	"fmt"
	"strings"
)

type userArticleSqlManager struct{}

func UserArticle() *userArticleSqlManager {
	return &userArticleSqlManager{}
}

type userArticleInsertSqlManager struct{}

func (userArticleSqlManager) Insert() *userArticleInsertSqlManager {
	return &userArticleInsertSqlManager{}
}

func (userArticleInsertSqlManager) Rating() string {
	return `INSERT INTO user_article (user_id, article_id, rating) VALUES ($1, $2, $3)`
}

func (userArticleInsertSqlManager) ViewLater() string {
	return `INSERT INTO user_article (user_id, article_id, view_later) VALUES ($1, $2, $3)`
}

type userArticleUpdateSqlManager struct{}

func (userArticleSqlManager) Update() *userArticleUpdateSqlManager {
	return &userArticleUpdateSqlManager{}
}

func (userArticleUpdateSqlManager) Rating() string {
	return `UPDATE user_article
			SET rating = $1, updated_at = TIMEZONE('America/Sao_Paulo'::TEXT, NOW())
    		WHERE active = true AND user_id = $2 AND article_id = $3`
}

func (userArticleUpdateSqlManager) ViewLater() string {
	return `UPDATE user_article
			SET view_later = $1, updated_at = TIMEZONE('America/Sao_Paulo'::TEXT, NOW())
			WHERE active = true AND user_id = $2 AND article_id = $3`
}

type userArticleSelectSqlManager struct{}

func (userArticleSqlManager) Select() *userArticleSelectSqlManager {
	return &userArticleSelectSqlManager{}
}

func (userArticleSelectSqlManager) RatingsAndArticlesSavedForLaterViewing(numberOfArticles int) string {
	var parameters []string
	for i := 1; i <= numberOfArticles; i++ {
		parameters = append(parameters, fmt.Sprintf("$%d", i+1))
	}

	return fmt.Sprintf(`SELECT article.id AS article_id,
							COALESCE(user_article.rating, 0) AS user_article_rating, user_article.view_later AS user_article_view_later
						FROM article
							LEFT JOIN proposition ON proposition.id = article.proposition_id
							LEFT JOIN newsletter ON newsletter.id = article.newsletter_id
							LEFT JOIN user_article ON user_article.article_id = article.id
						WHERE article.active = true AND (proposition.active = true OR newsletter.active = true)
							AND user_article.user_id = $1 AND article.id IN (%s)
						GROUP BY article.id, article.reference_date_time, user_article.rating, user_article.view_later
						ORDER BY article.reference_date_time DESC`, strings.Join(parameters, ","))
}

func (userArticleSelectSqlManager) NumberOfArticlesSavedToViewLater() string {
	return `SELECT COUNT(DISTINCT article.id)
			FROM article
			    LEFT JOIN proposition ON proposition.id = article.proposition_id
			    LEFT JOIN newsletter ON newsletter.id = article.newsletter_id
				LEFT JOIN user_article ON user_article.article_id = article.id
			WHERE article.active = true AND (proposition.active = true OR newsletter.active = true) AND
				user_article.active = true AND ((proposition.title ILIKE $1 OR proposition.content ILIKE $1) OR
				(newsletter.title ILIKE $1 OR newsletter.description ILIKE $1)) AND
				DATE(article.created_at) >= DATE(COALESCE($2, article.created_at)) AND
				DATE(article.created_at) <= DATE(COALESCE($3, article.created_at)) AND user_id = $4`
}

func (userArticleSelectSqlManager) NumberOfPropositionsSavedToViewLater() string {
	return `SELECT COUNT(DISTINCT article.id)
			FROM article
			    LEFT JOIN proposition prop ON prop.id = article.proposition_id
			    LEFT JOIN proposition_author ON proposition_author.proposition_id = prop.id
			    LEFT JOIN deputy ON deputy.id = proposition_author.deputy_id
				LEFT JOIN party party_in_the_proposal ON party_in_the_proposal.id = proposition_author.party_id
				LEFT JOIN external_author ON external_author.id = proposition_author.external_author_id
				LEFT JOIN user_article ON user_article.article_id = article.id
			WHERE article.active = true AND prop.active = true AND proposition_author.active = true AND
				user_article.active = true AND (prop.title ILIKE $1 OR prop.content ILIKE $1) AND
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
				DATE(article.created_at) <= DATE(COALESCE($6, article.created_at)) AND user_id = $7`
}

func (userArticleSelectSqlManager) NumberOfNewslettersSavedToViewLater() string {
	return `SELECT COUNT(*)
			FROM article
				LEFT JOIN newsletter ON newsletter.id = article.newsletter_id
				LEFT JOIN user_article ON user_article.article_id = article.id
    		WHERE article.active = true AND newsletter.active = true AND user_article.active = true AND
    			(newsletter.title ILIKE $1 OR newsletter.description ILIKE $1) AND
    			DATE(article.created_at) >= DATE(COALESCE($2, article.created_at)) AND
				DATE(article.created_at) <= DATE(COALESCE($3, article.created_at)) AND user_id = $4`
}

func (userArticleSelectSqlManager) ArticlesSavedToViewLater() string {
	return `SELECT article.id AS article_id,
    			COALESCE(user_article.rating, 0) AS user_article_rating, user_article.view_later AS user_article_view_later
			FROM article
			    LEFT JOIN proposition ON proposition.id = article.proposition_id
			    LEFT JOIN newsletter ON newsletter.id = article.newsletter_id
				LEFT JOIN user_article ON user_article.article_id = article.id
			WHERE article.active = true AND (proposition.active = true OR newsletter.active = true) AND
				user_article.active =true AND
				((proposition.title ILIKE $1 OR proposition.content ILIKE $1) OR
				(newsletter.title ILIKE $1 OR newsletter.description ILIKE $1)) AND
				DATE(article.created_at) >= DATE(COALESCE($2, article.created_at)) AND
				DATE(article.created_at) <= DATE(COALESCE($3, article.created_at)) AND
				user_article.user_id = $4 AND user_article.view_later = true
			ORDER BY article.reference_date_time DESC
			OFFSET $5 LIMIT $6`
}

func (userArticleSelectSqlManager) PropositionsSavedToViewLater() string {
	return `SELECT article.id AS article_id,
				COALESCE(user_article.rating, 0) AS user_article_rating, user_article.view_later AS user_article_view_later
			FROM article
			    LEFT JOIN proposition prop ON prop.id = article.proposition_id
			    LEFT JOIN proposition_author ON proposition_author.proposition_id = prop.id
			    LEFT JOIN deputy ON deputy.id = proposition_author.deputy_id
				LEFT JOIN party party_in_the_proposal ON party_in_the_proposal.id = proposition_author.party_id
				LEFT JOIN external_author ON external_author.id = proposition_author.external_author_id
				LEFT JOIN user_article ON user_article.article_id = article.id
			WHERE article.active = true AND prop.active = true AND proposition_author.active = true AND
				user_article.active = true AND (prop.title ILIKE $1 OR prop.content ILIKE $1) AND 
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
				DATE(article.created_at) <= DATE(COALESCE($6, article.created_at)) AND
				user_article.user_id = $7 AND user_article.view_later = true
			GROUP BY article.id, article.reference_date_time, prop.id, user_article.rating, user_article.view_later
			ORDER BY article.reference_date_time DESC OFFSET $8 LIMIT $9`
}

func (userArticleSelectSqlManager) NewslettersSavedToViewLater() string {
	return `SELECT article.id AS article_id,
				COALESCE(user_article.rating, 0) AS user_article_rating, user_article.view_later AS user_article_view_later
			FROM article
				LEFT JOIN newsletter ON newsletter.id = article.newsletter_id
    			LEFT JOIN user_article ON user_article.article_id = article.id
			WHERE article.active = true AND newsletter.active = true AND user_article.active = true AND 
				(newsletter.title ILIKE $1 OR newsletter.description ILIKE $1) AND
				DATE(article.created_at) >= DATE(COALESCE($2, article.created_at)) AND
				DATE(article.created_at) <= DATE(COALESCE($3, article.created_at)) AND
				user_article.user_id = $7 AND user_article.view_later = true
			ORDER BY article.reference_date_time DESC OFFSET $4 LIMIT $5`
}
