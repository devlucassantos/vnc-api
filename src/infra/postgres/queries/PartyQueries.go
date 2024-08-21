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
	return `SELECT id AS party_id, name AS party_name, acronym AS party_acronym, image_url AS party_image_url,
       			created_at AS party_created_at, updated_at AS party_updated_at
    		FROM party
    		WHERE party.active = true`
}
