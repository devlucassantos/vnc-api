package filters

type PaginationFilter struct {
	Page         *int
	ItemsPerPage *int
}

const DefaultPageNumberFilter = 1
const DefaultNumberOfItemsPerPageFilter = 15

func (instance PaginationFilter) GetPage() int {
	if instance.Page == nil {
		return DefaultPageNumberFilter
	}

	return *instance.Page
}

func (instance PaginationFilter) GetItemsPerPage() int {
	if instance.ItemsPerPage == nil {
		return DefaultNumberOfItemsPerPageFilter
	}

	return *instance.ItemsPerPage
}

func (instance PaginationFilter) CalculateOffset() int {
	return (instance.GetPage() - 1) * instance.GetItemsPerPage()
}
