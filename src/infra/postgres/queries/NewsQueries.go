package queries

type newsSqlManager struct{}

func News() *newsSqlManager {
	return &newsSqlManager{}
}

type newsSelectSqlManager struct{}

func (newsSqlManager) Select() *newsSelectSqlManager {
	return &newsSelectSqlManager{}
}

func (newsSelectSqlManager) TotalNumberOfNews() string {
	return `SELECT COUNT(DISTINCT news.id) FROM news
			    LEFT JOIN proposition ON proposition.id = news.proposition_id
			    LEFT JOIN proposition_author ON proposition_author.proposition_id = proposition.id
			    LEFT JOIN deputy ON deputy.id = proposition_author.deputy_id
				LEFT JOIN party party_in_the_proposal ON party_in_the_proposal.id = proposition_author.party_id
				LEFT JOIN organization ON organization.id = proposition_author.organization_id
				LEFT JOIN newsletter ON newsletter.id = news.newsletter_id
			WHERE news.active = true AND (proposition.active = true OR newsletter.active = true) AND
				((proposition.title ILIKE $1 OR proposition.content ILIKE $1) OR (newsletter.title ILIKE $1 OR newsletter.content ILIKE $1)) AND
				DATE(news.created_at) >= DATE(COALESCE($2, news.created_at)) AND
				DATE(news.created_at) <= DATE(COALESCE($3, news.created_at))`
}

func (newsSelectSqlManager) TotalNumberOfNewsletters() string {
	return `SELECT COUNT(*) FROM news
				LEFT JOIN newsletter ON newsletter.id = news.newsletter_id
    		WHERE news.active = true AND newsletter.active = true AND
    			(newsletter.title ILIKE $1 OR newsletter.content ILIKE $1) AND
    			DATE(news.created_at) >= DATE(COALESCE($2, news.created_at)) AND
				DATE(news.created_at) <= DATE(COALESCE($3, news.created_at))`
}

func (newsSelectSqlManager) TotalNumberOfPropositions() string {
	return `SELECT COUNT(DISTINCT news.id) FROM news
			    LEFT JOIN proposition prop ON prop.id = news.proposition_id
			    LEFT JOIN proposition_author ON proposition_author.proposition_id = prop.id
			    LEFT JOIN deputy ON deputy.id = proposition_author.deputy_id
				LEFT JOIN party party_in_the_proposal ON party_in_the_proposal.id = proposition_author.party_id
				LEFT JOIN organization ON organization.id = proposition_author.organization_id
			WHERE news.active = true AND prop.active = true AND
				(prop.title ILIKE $1 OR prop.content ILIKE $1) AND
				((deputy.id IS NULL AND ($2::uuid IS NULL AND $4::uuid IS NOT NULL))
	    			OR $2::uuid IS NULL
					OR (SELECT EXISTS(SELECT 1 FROM news
							LEFT JOIN proposition prop2 ON prop2.id = news.proposition_id
							LEFT JOIN proposition_author ON proposition_author.proposition_id = prop2.id
						WHERE proposition_author.deputy_id = $2 AND news.proposition_id = prop.id))) AND
				((party_in_the_proposal.id IS NULL AND ($3::uuid IS NULL AND $4::uuid IS NOT NULL))
					OR $3::uuid IS NULL
					OR (SELECT EXISTS(SELECT 1 FROM news
							LEFT JOIN proposition prop2 ON prop2.id = news.proposition_id
							LEFT JOIN proposition_author ON proposition_author.proposition_id = prop2.id
	                	WHERE proposition_author.party_id = $3 AND news.proposition_id = prop.id))) AND
				((organization.id IS NULL AND (($2::uuid IS NOT NULL OR $3::uuid IS NOT NULL) AND $4::uuid IS NULL))
					OR $4::uuid IS NULL
					OR (SELECT EXISTS(SELECT 1 FROM news
							LEFT JOIN proposition prop2 ON prop2.id = news.proposition_id
							LEFT JOIN proposition_author ON proposition_author.proposition_id = prop2.id
					WHERE proposition_author.organization_id = $4 AND news.proposition_id = prop.id))) AND
			    news.newsletter_id IS NULL AND
				DATE(news.created_at) >= DATE(COALESCE($5, news.created_at)) AND
				DATE(news.created_at) <= DATE(COALESCE($6, news.created_at))`
}

func (newsSelectSqlManager) All() string {
	return `SELECT COALESCE(proposition.id, '00000000-0000-0000-0000-000000000000') AS proposition_id,
				COALESCE(proposition.code, 0) AS proposition_code,
				COALESCE(proposition.original_text_url, '') AS proposition_original_text_url,
	   			COALESCE(proposition.title, '') AS proposition_title,
				COALESCE(proposition.content, '') AS proposition_content,
				COALESCE(proposition.submitted_at, '1970-01-01 00:00:00') AS proposition_submitted_at,
	   			COALESCE(proposition.active, true) AS proposition_active,
				COALESCE(proposition.created_at, '1970-01-01 00:00:00') AS proposition_created_at,
				COALESCE(proposition.updated_at, '1970-01-01 00:00:00') AS proposition_updated_at,
				
				COALESCE(newsletter.id, '00000000-0000-0000-0000-000000000000') AS newsletter_id,
	   			COALESCE(newsletter.title, '') AS newsletter_title,
				COALESCE(newsletter.content, '') AS newsletter_content,
				COALESCE(newsletter.reference_date, '1970-01-01 00:00:00') AS newsletter_reference_date,
	   			COALESCE(newsletter.active, true) AS newsletter_active,
				COALESCE(newsletter.created_at, '1970-01-01 00:00:00') AS newsletter_created_at,
				COALESCE(newsletter.updated_at, '1970-01-01 00:00:00') AS newsletter_updated_at
			FROM news
			    LEFT JOIN proposition ON proposition.id = news.proposition_id
			    LEFT JOIN proposition_author ON proposition_author.proposition_id = proposition.id
			    LEFT JOIN deputy ON deputy.id = proposition_author.deputy_id
				LEFT JOIN party party_in_the_proposal ON party_in_the_proposal.id = proposition_author.party_id
				LEFT JOIN organization ON organization.id = proposition_author.organization_id
				LEFT JOIN newsletter ON newsletter.id = news.newsletter_id
			WHERE news.active = true AND (proposition.active = true OR newsletter.active = true) AND
				((proposition.title ILIKE $1 OR proposition.content ILIKE $1) OR (newsletter.title ILIKE $1 OR newsletter.content ILIKE $1)) AND
				DATE(news.created_at) >= DATE(COALESCE($2, news.created_at)) AND
				DATE(news.created_at) <= DATE(COALESCE($3, news.created_at))
			GROUP BY proposition.id, proposition.code, proposition.original_text_url, proposition.title, proposition.content,
			         proposition.submitted_at, proposition.active, proposition.created_at, proposition.updated_at,
			         newsletter.id, newsletter.title, newsletter.content, newsletter.reference_date, newsletter.active,
			         newsletter.created_at, newsletter.updated_at, news.reference_date_time
			ORDER BY news.reference_date_time DESC OFFSET $4 LIMIT $5`
}

func (newsSelectSqlManager) Newsletters() string {
	return `SELECT COALESCE(newsletter.id, '00000000-0000-0000-0000-000000000000') AS newsletter_id,
				COALESCE(newsletter.title, '') AS newsletter_title,
				COALESCE(newsletter.content, '') AS newsletter_content,
				COALESCE(newsletter.reference_date, '1970-01-01 00:00:00') AS newsletter_reference_date,
				COALESCE(newsletter.active, true) AS newsletter_active,
				COALESCE(newsletter.created_at, '1970-01-01 00:00:00') AS newsletter_created_at,
				COALESCE(newsletter.updated_at, '1970-01-01 00:00:00') AS newsletter_updated_at
    		FROM news
				LEFT JOIN newsletter ON newsletter.id = news.newsletter_id
    		WHERE news.active = true AND newsletter.active = true AND
    			(newsletter.title ILIKE $1 OR newsletter.content ILIKE $1) AND
    			DATE(news.created_at) >= DATE(COALESCE($2, news.created_at)) AND
				DATE(news.created_at) <= DATE(COALESCE($3, news.created_at))
    		ORDER BY news.reference_date_time DESC OFFSET $4 LIMIT $5`
}

func (newsSelectSqlManager) Propositions() string {
	return `SELECT COALESCE(prop.id, '00000000-0000-0000-0000-000000000000') AS proposition_id,
				COALESCE(prop.code, 0) AS proposition_code,
				COALESCE(prop.original_text_url, '') AS proposition_original_text_url,
	   			COALESCE(prop.title, '') AS proposition_title,
				COALESCE(prop.content, '') AS proposition_content,
				COALESCE(prop.submitted_at, '1970-01-01 00:00:00') AS proposition_submitted_at,
	   			COALESCE(prop.active, true) AS proposition_active,
				COALESCE(prop.created_at, '1970-01-01 00:00:00') AS proposition_created_at,
				COALESCE(prop.updated_at, '1970-01-01 00:00:00') AS proposition_updated_at
			FROM news
			    LEFT JOIN proposition prop ON prop.id = news.proposition_id
			    LEFT JOIN proposition_author ON proposition_author.proposition_id = prop.id
			    LEFT JOIN deputy ON deputy.id = proposition_author.deputy_id
				LEFT JOIN party party_in_the_proposal ON party_in_the_proposal.id = proposition_author.party_id
				LEFT JOIN organization ON organization.id = proposition_author.organization_id
			WHERE news.active = true AND prop.active = true AND
				(prop.title ILIKE $1 OR prop.content ILIKE $1) AND
				((deputy.id IS NULL AND ($2::uuid IS NULL AND $4::uuid IS NOT NULL))
	    			OR $2::uuid IS NULL
					OR (SELECT EXISTS(SELECT 1 FROM news
							LEFT JOIN proposition prop2 ON prop2.id = news.proposition_id
							LEFT JOIN proposition_author ON proposition_author.proposition_id = prop2.id
						WHERE proposition_author.deputy_id = $2 AND news.proposition_id = prop.id))) AND
				((party_in_the_proposal.id IS NULL AND ($3::uuid IS NULL AND $4::uuid IS NOT NULL))
					OR $3::uuid IS NULL
					OR (SELECT EXISTS(SELECT 1 FROM news
							LEFT JOIN proposition prop2 ON prop2.id = news.proposition_id
							LEFT JOIN proposition_author ON proposition_author.proposition_id = prop2.id
	                	WHERE proposition_author.party_id = $3 AND news.proposition_id = prop.id))) AND
				((organization.id IS NULL AND (($2::uuid IS NOT NULL OR $3::uuid IS NOT NULL) AND $4::uuid IS NULL))
					OR $4::uuid IS NULL
					OR (SELECT EXISTS(SELECT 1 FROM news
							LEFT JOIN proposition prop2 ON prop2.id = news.proposition_id
							LEFT JOIN proposition_author ON proposition_author.proposition_id = prop2.id
					WHERE proposition_author.organization_id = $4 AND news.proposition_id = prop.id))) AND
			    news.newsletter_id IS NULL AND
				DATE(news.created_at) >= DATE(COALESCE($5, news.created_at)) AND
				DATE(news.created_at) <= DATE(COALESCE($6, news.created_at))
			GROUP BY prop.id, prop.code, prop.original_text_url, prop.title, prop.content, prop.submitted_at, prop.active,
			         prop.created_at, prop.updated_at, news.reference_date_time
			ORDER BY news.reference_date_time DESC OFFSET $7 LIMIT $8`
}

func (newsSelectSqlManager) TrendingNews() string {
	return `SELECT COALESCE(proposition.id, '00000000-0000-0000-0000-000000000000') AS proposition_id,
				COALESCE(proposition.code, 0) AS proposition_code,
				COALESCE(proposition.original_text_url, '') AS proposition_original_text_url,
				COALESCE(proposition.title, '') AS proposition_title,
				COALESCE(proposition.content, '') AS proposition_content,
				COALESCE(proposition.submitted_at, '1970-01-01 00:00:00') AS proposition_submitted_at,
				COALESCE(proposition.active, true) AS proposition_active,
				COALESCE(proposition.created_at, '1970-01-01 00:00:00') AS proposition_created_at,
				COALESCE(proposition.updated_at, '1970-01-01 00:00:00') AS proposition_updated_at,
				
				COALESCE(newsletter.id, '00000000-0000-0000-0000-000000000000') AS newsletter_id,
				COALESCE(newsletter.title, '') AS newsletter_title,
				COALESCE(newsletter.content, '') AS newsletter_content,
				COALESCE(newsletter.reference_date, '1970-01-01 00:00:00') AS newsletter_reference_date,
				COALESCE(newsletter.active, true) AS newsletter_active,
				COALESCE(newsletter.created_at, '1970-01-01 00:00:00') AS newsletter_created_at,
				COALESCE(newsletter.updated_at, '1970-01-01 00:00:00') AS newsletter_updated_at,
				
				COUNT(news_view.id) AS news_views
			FROM news
				LEFT JOIN proposition ON proposition.id = news.proposition_id
				LEFT JOIN newsletter ON newsletter.id = news.newsletter_id
				LEFT JOIN news_view ON news_view.news_id = news.id AND news_view.created_at >= CURRENT_DATE - INTERVAL '1 week'
			WHERE news.active = true AND (proposition.active = true OR newsletter.active = true) AND
    			((proposition.title ILIKE $1 OR proposition.content ILIKE $1) OR (newsletter.title ILIKE $1 OR newsletter.content ILIKE $1)) AND
				DATE(news.created_at) >= DATE(COALESCE($2, news.created_at)) AND
				DATE(news.created_at) <= DATE(COALESCE($3, news.created_at))
			GROUP BY proposition.id, proposition_code, proposition_original_text_url, proposition_title, proposition_content,
				proposition_submitted_at, proposition_active, proposition_created_at, proposition_updated_at,
				newsletter.id, newsletter_title, newsletter_content, newsletter_reference_date, newsletter_active,
				newsletter_created_at, newsletter_updated_at,
				news_view.news_id, news.reference_date_time
			ORDER BY news_views DESC, news.reference_date_time DESC OFFSET $4 LIMIT $5`
}

func (newsSelectSqlManager) TrendingNewsletters() string {
	return `SELECT COALESCE(newsletter.id, '00000000-0000-0000-0000-000000000000') AS newsletter_id,
				COALESCE(newsletter.title, '') AS newsletter_title,
				COALESCE(newsletter.content, '') AS newsletter_content,
				COALESCE(newsletter.reference_date, '1970-01-01 00:00:00') AS newsletter_reference_date,
				COALESCE(newsletter.active, true) AS newsletter_active,
				COALESCE(newsletter.created_at, '1970-01-01 00:00:00') AS newsletter_created_at,
				COALESCE(newsletter.updated_at, '1970-01-01 00:00:00') AS newsletter_updated_at,
				
				COUNT(news_view.id) AS news_views
			FROM news
				LEFT JOIN newsletter ON newsletter.id = news.newsletter_id
				LEFT JOIN news_view ON news_view.news_id = news.id AND news_view.created_at >= CURRENT_DATE - INTERVAL '1 week'
			WHERE news.active = true AND newsletter.active = true AND
    			(newsletter.title ILIKE $1 OR newsletter.content ILIKE $1) AND
    			DATE(news.created_at) >= DATE(COALESCE($2, news.created_at)) AND
				DATE(news.created_at) <= DATE(COALESCE($3, news.created_at))
			GROUP BY newsletter.id, newsletter_title, newsletter_content, newsletter_reference_date, newsletter_active,
				newsletter_created_at, newsletter_updated_at, news_view.news_id, news.reference_date_time
			ORDER BY news_views DESC, news.reference_date_time DESC OFFSET $4 LIMIT $5`
}

func (newsSelectSqlManager) TrendingPropositions() string {
	return `SELECT COALESCE(prop.id, '00000000-0000-0000-0000-000000000000') AS proposition_id,
				COALESCE(prop.code, 0) AS proposition_code,
				COALESCE(prop.original_text_url, '') AS proposition_original_text_url,
				COALESCE(prop.title, '') AS proposition_title,
				COALESCE(prop.content, '') AS proposition_content,
				COALESCE(prop.submitted_at, '1970-01-01 00:00:00') AS proposition_submitted_at,
				COALESCE(prop.active, true) AS proposition_active,
				COALESCE(prop.created_at, '1970-01-01 00:00:00') AS proposition_created_at,
				COALESCE(prop.updated_at, '1970-01-01 00:00:00') AS proposition_updated_at,
				
				COUNT(news_view.id) AS news_views
			FROM news
				LEFT JOIN proposition prop ON prop.id = news.proposition_id
			    LEFT JOIN proposition_author ON proposition_author.proposition_id = prop.id
			    LEFT JOIN deputy ON deputy.id = proposition_author.deputy_id
				LEFT JOIN party party_in_the_proposal ON party_in_the_proposal.id = proposition_author.party_id
				LEFT JOIN organization ON organization.id = proposition_author.organization_id
				LEFT JOIN news_view ON news_view.news_id = news.id AND news_view.created_at >= CURRENT_DATE - INTERVAL '1 week'
			WHERE news.active = true AND prop.active = true AND
				(prop.title ILIKE $1 OR prop.content ILIKE $1) AND
				((deputy.id IS NULL AND ($2::uuid IS NULL AND $4::uuid IS NOT NULL))
	    			OR $2::uuid IS NULL
					OR (SELECT EXISTS(SELECT 1 FROM news
							LEFT JOIN proposition prop2 ON prop2.id = news.proposition_id
							LEFT JOIN proposition_author ON proposition_author.proposition_id = prop2.id
						WHERE proposition_author.deputy_id = $2 AND news.proposition_id = prop.id))) AND
				((party_in_the_proposal.id IS NULL AND ($3::uuid IS NULL AND $4::uuid IS NOT NULL))
					OR $3::uuid IS NULL
					OR (SELECT EXISTS(SELECT 1 FROM news
							LEFT JOIN proposition prop2 ON prop2.id = news.proposition_id
							LEFT JOIN proposition_author ON proposition_author.proposition_id = prop2.id
	                	WHERE proposition_author.party_id = $3 AND news.proposition_id = prop.id))) AND
				((organization.id IS NULL AND (($2::uuid IS NOT NULL OR $3::uuid IS NOT NULL) AND $4::uuid IS NULL))
					OR $4::uuid IS NULL
					OR (SELECT EXISTS(SELECT 1 FROM news
							LEFT JOIN proposition prop2 ON prop2.id = news.proposition_id
							LEFT JOIN proposition_author ON proposition_author.proposition_id = prop2.id
					WHERE proposition_author.organization_id = $4 AND news.proposition_id = prop.id))) AND
			    news.newsletter_id IS NULL AND
				DATE(news.created_at) >= DATE(COALESCE($5, news.created_at)) AND
				DATE(news.created_at) <= DATE(COALESCE($6, news.created_at))
			GROUP BY prop.id, proposition_code, proposition_original_text_url, proposition_title, proposition_content,
				proposition_submitted_at, proposition_active, proposition_created_at, proposition_updated_at, news_view.news_id,
				news.reference_date_time
			ORDER BY news_views DESC, news.reference_date_time DESC OFFSET $7 LIMIT $8`
}
