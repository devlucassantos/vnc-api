package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"net/http"
	"vnc-read-api/api/endpoints/dto/filter"
	"vnc-read-api/api/endpoints/dto/response"
	"vnc-read-api/core/interfaces/services"
)

type News struct {
	service services.News
}

func NewNewsHandler(service services.News) *News {
	return &News{service: service}
}

// GetNews
// @ID 			GetNews
// @Summary 	Listagem das matérias mais recentes
// @Tags 		Matérias
// @Description Esta requisição é responsável por retornar as matérias mais recentes disponíveis na plataforma Você na Câmara.
// @Produce		json
// @Param page         query int false "Número da página. Por padrão é 1."
// @Param itemsPerPage query int false "Quantidade de matérias retornadas por página. Por padrão é 25."
// @Success 200 {array}  response.SwaggerNewsPagination "Requisição bem sucedida"
// @Failure 400 {object} response.SwaggerError          "Algum dado informado durante a requisição é inválido"
// @Failure 500 {object} response.SwaggerError          "Ocorreu um erro inesperado durante o processamento da requisição"
// @Router /news [get]
func (instance News) GetNews(context echo.Context) error {
	var newsFilter filter.NewsFilter
	pageParam := context.QueryParam("page")
	if pageParam != "" {
		page, err := convertToInt(pageParam, "Página")
		if err != nil {
			log.Error(err)
			return context.JSON(http.StatusBadRequest, response.NewError(err.Error()))
		}
		newsFilter.PaginationFilter.Page = &page
	}

	itemsPerPageParam := context.QueryParam("itemsPerPage")
	if itemsPerPageParam != "" {
		perPage, err := convertToInt(itemsPerPageParam, "Itens por página")
		if err != nil {
			log.Error(err)
			return context.JSON(http.StatusBadRequest, response.NewError(err.Error()))
		}
		newsFilter.PaginationFilter.ItemsPerPage = &perPage
	}

	newsList, totalNumberOfNews, err := instance.service.GetNews(newsFilter)
	if err != nil {
		return context.JSON(http.StatusInternalServerError, response.NewError(err.Error()))
	}

	var newsResponse []response.News
	for _, news := range newsList {
		newsResponse = append(newsResponse, *response.NewNews(news))
	}

	requestResult := response.Pagination{
		Page:         newsFilter.PaginationFilter.GetPage(),
		ItensPerPage: newsFilter.PaginationFilter.GetItemsPerPage(),
		Total:        totalNumberOfNews,
		Data:         newsResponse,
	}

	return context.JSON(http.StatusOK, requestResult)
}
