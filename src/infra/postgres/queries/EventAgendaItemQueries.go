package queries

type eventAgendaItemSqlManager struct{}

func EventAgendaItem() *eventAgendaItemSqlManager {
	return &eventAgendaItemSqlManager{}
}

type eventAgendaItemSelectSqlManager struct{}

func (eventAgendaItemSqlManager) Select() *eventAgendaItemSelectSqlManager {
	return &eventAgendaItemSelectSqlManager{}
}

func (eventAgendaItemSelectSqlManager) ByEventId() string {
	return `SELECT event_agenda_item.id AS event_agenda_item_id, event_agenda_item.title AS event_agenda_item_title,
				event_agenda_item.topic AS event_agenda_item_topic,
				event_agenda_item.situation AS event_agenda_item_situation,
				COALESCE(event_agenda_item.rapporteur_federated_unit, '') AS deputy_previous_federated_unit,
				proposition.article_id AS event_agenda_item_proposition_article_id,
				COALESCE(related_proposition.article_id, '00000000-0000-0000-0000-000000000000')
					AS event_agenda_item_related_proposition_article_id,
				COALESCE(voting.article_id, '00000000-0000-0000-0000-000000000000')
					AS event_agenda_item_voting_article_id,
				agenda_item_regime.id AS agenda_item_regime_id,
				agenda_item_regime.description AS agenda_item_regime_description,
				COALESCE(deputy.id, '00000000-0000-0000-0000-000000000000') AS deputy_id,
				COALESCE(deputy.name, '') AS deputy_name,
				COALESCE(deputy.electoral_name, '') AS deputy_electoral_name,
				COALESCE(deputy.image_url, '') AS deputy_image_url,
				COALESCE(deputy.federated_unit, '') AS deputy_federated_unit,
				COALESCE(party.id, '00000000-0000-0000-0000-000000000000') AS party_id,
				COALESCE(party.name, '') AS party_name,
				COALESCE(party.acronym, '') AS party_acronym,
				COALESCE(party.image_url, '') AS party_image_url,
				COALESCE(previous_party.id, '00000000-0000-0000-0000-000000000000') AS previous_party_id,
				COALESCE(previous_party.name, '') AS previous_party_name,
				COALESCE(previous_party.acronym, '') AS previous_party_acronym,
				COALESCE(previous_party.image_url, '') AS previous_party_image_url
			FROM event_agenda_item
				INNER JOIN agenda_item_regime ON agenda_item_regime.id = event_agenda_item.agenda_item_regime_id
				LEFT JOIN deputy ON deputy.id = event_agenda_item.rapporteur_id
				LEFT JOIN party ON party.id = deputy.party_id
				LEFT JOIN party previous_party ON previous_party.id = event_agenda_item.rapporteur_party_id
				INNER JOIN proposition ON proposition.id = event_agenda_item.proposition_id
				LEFT JOIN proposition related_proposition ON related_proposition.id =
					event_agenda_item.related_proposition_id
				LEFT JOIN voting ON voting.id = event_agenda_item.voting_id
			WHERE event_agenda_item.active = true AND agenda_item_regime.active = true AND
				deputy.active IS NOT false AND party.active IS NOT false AND proposition.active = true AND
				related_proposition.active IS NOT false AND voting.active IS NOT false AND
				event_agenda_item.event_id = $1`
}
