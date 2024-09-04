package dto

type UserArticle struct {
	Rating    int  `db:"user_article_rating"`
	ViewLater bool `db:"user_article_view_later"`
	*User
	*Article
}
