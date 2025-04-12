package queries

type legislativeBodySqlManager struct{}

func LegislativeBody() *legislativeBodySqlManager {
	return &legislativeBodySqlManager{}
}

type legislativeBodySelectSqlManager struct{}

func (legislativeBodySqlManager) Select() *legislativeBodySelectSqlManager {
	return &legislativeBodySelectSqlManager{}
}

func (legislativeBodySelectSqlManager) All() string {
	return `SELECT legislative_body.id AS legislative_body_id,
       			legislative_body.name AS legislative_body_name,
				legislative_body.acronym AS legislative_body_acronym,
				legislative_body_type.id AS legislative_body_type_id,
				legislative_body_type.description AS legislative_body_type_description
			FROM legislative_body
				INNER JOIN legislative_body_type ON legislative_body_type.id = legislative_body.legislative_body_type_id
			WHERE legislative_body.active = true AND legislative_body_type.active = true
			ORDER BY legislative_body.name, legislative_body_type.description`
}

func (legislativeBodySelectSqlManager) LegislativeBodiesByEventId() string {
	return `SELECT legislative_body.id AS legislative_body_id,
				legislative_body.name AS legislative_body_name,
				legislative_body.acronym AS legislative_body_acronym,
				legislative_body_type.id AS legislative_body_type_id,
				legislative_body_type.description AS legislative_body_type_description
			FROM legislative_body
				INNER JOIN legislative_body_type ON legislative_body_type.id = legislative_body.legislative_body_type_id
				INNER JOIN event_legislative_body ON event_legislative_body.legislative_body_id = legislative_body.id
			WHERE legislative_body.active = true AND legislative_body_type.active = true AND
				event_legislative_body.active = true AND event_legislative_body.event_id = $1`
}
