package queries

type eventSituationSqlManager struct{}

func EventSituation() *eventSituationSqlManager {
	return &eventSituationSqlManager{}
}

type eventSituationSelectSqlManager struct{}

func (eventSituationSqlManager) Select() *eventSituationSelectSqlManager {
	return &eventSituationSelectSqlManager{}
}

func (eventSituationSelectSqlManager) All() string {
	return `SELECT id AS event_situation_id, description AS event_situation_description, color AS event_situation_color
			FROM event_situation
			WHERE active = true
			ORDER BY description`
}
