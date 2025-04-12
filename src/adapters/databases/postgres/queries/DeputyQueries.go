package queries

type deputySqlManager struct{}

func Deputy() *deputySqlManager {
	return &deputySqlManager{}
}

type deputySelectSqlManager struct{}

func (deputySqlManager) Select() *deputySelectSqlManager {
	return &deputySelectSqlManager{}
}

func (deputySelectSqlManager) All() string {
	return `SELECT deputy.id AS deputy_id, deputy.name AS deputy_name, deputy.electoral_name AS deputy_electoral_name,
       			deputy.image_url AS deputy_image_url, deputy.federated_unit AS deputy_federated_unit,
        		party.id AS party_id, party.name AS party_name, party.acronym AS party_acronym,
        		party.image_url AS party_image_url
    		FROM deputy
    			INNER JOIN party ON party.id = deputy.party_id
    		WHERE deputy.active = true AND party.active = true
    		ORDER BY deputy.electoral_name, party.acronym`
}
