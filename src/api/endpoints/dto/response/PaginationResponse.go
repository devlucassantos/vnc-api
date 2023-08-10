package response

type Pagination struct {
	Page         int         `json:"page"`
	ItensPerPage int         `json:"itens_per_page"`
	Total        int         `json:"total"`
	Data         interface{} `json:"data"`
}
