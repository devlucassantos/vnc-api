package queries

type articleViewSqlManager struct{}

func ArticleView() *articleViewSqlManager {
	return &articleViewSqlManager{}
}

func (articleViewSqlManager) Insert() string {
	return `INSERT INTO article_view(article_id) VALUES ($1)`
}
