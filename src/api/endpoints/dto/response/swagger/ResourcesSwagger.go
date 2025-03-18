package swagger

type Resources struct {
	ArticleTypes      []ArticleType     `json:"article_types"`
	PropositionTypes  []PropositionType `json:"proposition_types"`
	Parties           []Party           `json:"parties"`
	Deputies          []DeputyResource  `json:"deputies"`
	ExternalAuthors   []ExternalAuthor  `json:"external_authors"`
	LegislativeBodies []LegislativeBody `json:"legislative_bodies"`
	EventTypes        []EventType       `json:"event_types"`
	EventSituations   []EventSituation  `json:"event_situations"`
}
