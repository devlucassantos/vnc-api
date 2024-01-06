package queries

type newsViewSqlManager struct{}

func NewsView() *newsViewSqlManager {
	return &newsViewSqlManager{}
}

func (newsViewSqlManager) Insert() string {
	return `INSERT INTO news_view(news_id) VALUES ($1)`
}
