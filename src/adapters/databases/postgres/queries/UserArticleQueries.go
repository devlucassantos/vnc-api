package queries

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
	return `INSERT INTO user_article (user_id, article_id, view_later, view_later_set_at)
			VALUES ($1, $2, $3, CASE WHEN $3 = false THEN NULL ELSE TIMEZONE('America/Sao_Paulo'::TEXT, NOW()) END)`
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
			SET view_later = $1, view_later_set_at = CASE WHEN $1 = false THEN NULL
			    ELSE TIMEZONE('America/Sao_Paulo'::TEXT, NOW()) END,
			    updated_at = TIMEZONE('America/Sao_Paulo'::TEXT, NOW())
			WHERE active = true AND user_id = $2 AND article_id = $3`
}
