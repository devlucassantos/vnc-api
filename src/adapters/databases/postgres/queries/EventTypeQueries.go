package queries

type eventTypeSqlManager struct{}

func EventType() *eventTypeSqlManager {
	return &eventTypeSqlManager{}
}

type eventTypeSelectSqlManager struct{}

func (eventTypeSqlManager) Select() *eventTypeSelectSqlManager {
	return &eventTypeSelectSqlManager{}
}

func (eventTypeSelectSqlManager) All() string {
	return `SELECT id AS event_type_id, description AS event_type_description, color AS event_type_color
			FROM event_type
			WHERE active = true
			ORDER BY sort_order`
}
