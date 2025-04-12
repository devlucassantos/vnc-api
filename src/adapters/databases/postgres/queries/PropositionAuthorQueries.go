package queries

type propositionAuthorSqlManager struct{}

func PropositionAuthor() *propositionAuthorSqlManager {
	return &propositionAuthorSqlManager{}
}

type propositionAuthorSelectSqlManager struct{}

func (propositionAuthorSqlManager) Select() *propositionAuthorSelectSqlManager {
	return &propositionAuthorSelectSqlManager{}
}

func (propositionAuthorSelectSqlManager) ByPropositionId() string {
	return `SELECT COALESCE(deputy.id, '00000000-0000-0000-0000-000000000000') AS deputy_id,
       			COALESCE(deputy.name, '') AS deputy_name,
       			COALESCE(deputy.electoral_name, '') AS deputy_electoral_name,
       			COALESCE(deputy.image_url, '') AS deputy_image_url,
				COALESCE(deputy.federated_unit, '') AS deputy_federated_unit,
				COALESCE(proposition_author.federated_unit, '') AS deputy_previous_federated_unit,
        		COALESCE(party.id, '00000000-0000-0000-0000-000000000000') AS party_id,
        		COALESCE(party.name, '') AS party_name,
        		COALESCE(party.acronym, '') AS party_acronym,
        		COALESCE(party.image_url, '') AS party_image_url,
        		COALESCE(previous_party.id, '00000000-0000-0000-0000-000000000000') AS previous_party_id,
        		COALESCE(previous_party.name, '') AS previous_party_name,
        		COALESCE(previous_party.acronym, '') AS previous_party_acronym,
        		COALESCE(previous_party.image_url, '') AS previous_party_image_url,
        		COALESCE(external_author.id, '00000000-0000-0000-0000-000000000000') AS external_author_id,
        		COALESCE(external_author.name, '') AS external_author_name,
				COALESCE(external_author_type.id, '00000000-0000-0000-0000-000000000000') AS external_author_type_id,
       			COALESCE(external_author_type.description, '') AS external_author_type_description
			FROM proposition_author
				LEFT JOIN deputy ON deputy.id = proposition_author.deputy_id
				LEFT JOIN party ON party.id = deputy.party_id
				LEFT JOIN party previous_party ON previous_party.id = proposition_author.party_id
				LEFT JOIN external_author ON external_author.id = proposition_author.external_author_id
				LEFT JOIN external_author_type ON external_author_type.id = external_author.external_author_type_id
			WHERE proposition_author.active = true AND ((deputy.active = true AND party.active = true AND
				previous_party.active = true) OR (external_author.active = true AND
				external_author_type.active = true)) AND proposition_author.proposition_id = $1`
}
