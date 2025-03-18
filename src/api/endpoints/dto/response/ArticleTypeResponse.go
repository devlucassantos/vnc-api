package response

import (
	"github.com/devlucassantos/vnc-domains/src/domains/article"
	"github.com/devlucassantos/vnc-domains/src/domains/articletype"
	"github.com/devlucassantos/vnc-domains/src/domains/eventtype"
	"github.com/devlucassantos/vnc-domains/src/domains/propositiontype"
	"github.com/google/uuid"
)

type ArticleType struct {
	Id            uuid.UUID     `json:"id"`
	Description   string        `json:"description"`
	Codes         string        `json:"codes,omitempty"`
	Color         string        `json:"color"`
	SpecificType  *ArticleType  `json:"specific_type,omitempty"`
	SpecificTypes []ArticleType `json:"specific_types,omitempty"`
	Articles      []Article     `json:"articles,omitempty"`
}

func NewArticleType(articleType articletype.ArticleType) *ArticleType {
	return &ArticleType{
		Id:          articleType.Id(),
		Description: articleType.Description(),
		Codes:       articleType.Codes(),
		Color:       articleType.Color(),
	}
}

func NewPropositionSpecificType(propositionType propositiontype.PropositionType) *ArticleType {
	return &ArticleType{
		Id:          propositionType.Id(),
		Description: propositionType.Description(),
		Color:       propositionType.Color(),
	}
}

func NewEventSpecificType(eventType eventtype.EventType) *ArticleType {
	return &ArticleType{
		Id:          eventType.Id(),
		Description: eventType.Description(),
		Color:       eventType.Color(),
	}
}

func SortingArticleTypeWithArticles(articles []article.Article, currentArticleTypes []ArticleType) []ArticleType {
	for _, articleDomain := range articles {
		articleType := articleDomain.Type()
		var articleTypeFound bool
		for index, currentArticleType := range currentArticleTypes {
			if currentArticleType.Id == articleType.Id() {
				currentArticleTypes[index].Articles = append(currentArticleType.Articles, *NewArticle(articleDomain))
				articleTypeFound = true
				break
			}
		}

		if !articleTypeFound {
			articleTypeResponse := *NewArticleType(articleType)
			articleTypeResponse.Articles = append(articleTypeResponse.Articles, *NewArticle(articleDomain))
			currentArticleTypes = append(currentArticleTypes, articleTypeResponse)
		}
	}

	return currentArticleTypes
}

func SortingArticleTypesWithSpecificTypesAndArticles(articles []article.Article,
	currentArticleTypes []ArticleType) []ArticleType {
	for _, articleDomain := range articles {
		articleType := articleDomain.Type()
		var articleTypeResponse *ArticleType
		for index, currentArticleType := range currentArticleTypes {
			if currentArticleType.Id == articleType.Id() {
				articleTypeResponse = &currentArticleTypes[index]
				break
			}
		}

		if articleTypeResponse == nil {
			specificType := articleDomain.SpecificType()
			articleSpecificType := *NewArticleType(specificType)
			articleSpecificType.Articles = append(articleSpecificType.Articles, *NewArticle(articleDomain))
			articleTypeResponse = NewArticleType(articleType)
			articleTypeResponse.SpecificTypes = append([]ArticleType{}, articleSpecificType)
			currentArticleTypes = append(currentArticleTypes, *articleTypeResponse)
			continue
		}

		var articleSpecificTypeFound bool
		specificType := articleDomain.SpecificType()
		for index, specificTypeResponse := range articleTypeResponse.SpecificTypes {
			if specificTypeResponse.Id == specificType.Id() {
				articleTypeResponse.SpecificTypes[index].Articles = append(specificTypeResponse.Articles,
					*NewArticle(articleDomain))
				articleSpecificTypeFound = true
				break
			}
		}

		if !articleSpecificTypeFound {
			articleSpecificType := *NewArticleType(specificType)
			articleSpecificType.Articles = append(articleSpecificType.Articles, *NewArticle(articleDomain))
			articleTypeResponse.SpecificTypes = append(articleTypeResponse.SpecificTypes, articleSpecificType)
		}
	}

	return currentArticleTypes
}
