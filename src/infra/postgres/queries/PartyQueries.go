package queries

type partySqlManager struct{}

func Party() *partySqlManager {
	return &partySqlManager{}
}

type partySelectSqlManager struct{}

func (partySqlManager) Select() *partySelectSqlManager {
	return &partySelectSqlManager{}
}

func (partySelectSqlManager) All() string {
	return `SELECT COALESCE(party.id, '00000000-0000-0000-0000-000000000000') AS party_id,
        		COALESCE(party.code, 0) AS party_code,
        		COALESCE(party.name, '') AS party_name,
        		COALESCE(party.acronym, '') AS party_acronym,
        		COALESCE(party.image_url, '') AS party_image_url,
        		COALESCE(party.active, true) AS party_active,
        		COALESCE(party.created_at, '1970-01-01 00:00:00') AS party_created_at,
        		COALESCE(party.updated_at, '1970-01-01 00:00:00') AS party_updated_at
    		FROM party WHERE party.active = true
    		ORDER BY party.acronym, party.name`
}
