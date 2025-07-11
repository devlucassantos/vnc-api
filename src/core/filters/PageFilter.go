package filters

type Pagination struct {
	Page         *int
	ItemsPerPage *int
}

const DefaultPageNumberFilter = 1
const DefaultNumberOfItemsPerPageFilter = 15

func (instance Pagination) GetPage() int {
	if instance.Page == nil {
		return DefaultPageNumberFilter
	}

	return *instance.Page
}

func (instance Pagination) GetItemsPerPage() int {
	if instance.ItemsPerPage == nil {
		return DefaultNumberOfItemsPerPageFilter
	}

	return *instance.ItemsPerPage
}

func (instance Pagination) CalculateOffset() int {
	return (instance.GetPage() - 1) * instance.GetItemsPerPage()
}
