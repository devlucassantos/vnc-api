package swagger

type ArticlePagination struct {
	Page         int       `json:"page"           example:"1"`
	ItemsPerPage int       `json:"items_per_page" example:"15"`
	Total        int       `json:"total"          example:"4562"`
	Data         []Article `json:"data"`
}
