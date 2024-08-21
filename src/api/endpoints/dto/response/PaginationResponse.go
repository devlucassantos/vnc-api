package response

type Pagination struct {
	Page         int         `json:"page"`
	ItemsPerPage int         `json:"items_per_page"`
	Total        int         `json:"total"`
	Data         interface{} `json:"data"`
}
