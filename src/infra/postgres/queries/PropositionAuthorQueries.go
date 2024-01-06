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
       			COALESCE(deputy.code, 0) AS deputy_code,
       			COALESCE(deputy.cpf, '') AS deputy_cpf,
       			COALESCE(deputy.name, '') AS deputy_name,
       			COALESCE(deputy.electoral_name, '') AS deputy_electoral_name,
       			COALESCE(deputy.image_url, '') AS deputy_image_url,
       			COALESCE(deputy.active, true) AS deputy_active,
       			COALESCE(deputy.created_at, '1970-01-01 00:00:00') AS deputy_created_at,
       			COALESCE(deputy.updated_at, '1970-01-01 00:00:00') AS deputy_updated_at,
       			
        		COALESCE(party.id, '00000000-0000-0000-0000-000000000000') AS party_id,
        		COALESCE(party.code, 0) AS party_code,
        		COALESCE(party.name, '') AS party_name,
        		COALESCE(party.acronym, '') AS party_acronym,
        		COALESCE(party.image_url, '') AS party_image_url,
        		COALESCE(party.active, true) AS party_active,
        		COALESCE(party.created_at, '1970-01-01 00:00:00') AS party_created_at,
        		COALESCE(party.updated_at, '1970-01-01 00:00:00') AS party_updated_at,
        		
        		COALESCE(party_in_the_proposal.id, '00000000-0000-0000-0000-000000000000') AS party_in_the_proposal_id,
        		COALESCE(party_in_the_proposal.code, 0) AS party_in_the_proposal_code,
        		COALESCE(party_in_the_proposal.name, '') AS party_in_the_proposal_name,
        		COALESCE(party_in_the_proposal.acronym, '') AS party_in_the_proposal_acronym,
        		COALESCE(party_in_the_proposal.image_url, '') AS party_in_the_proposal_image_url,
        		COALESCE(party_in_the_proposal.active, true) AS party_in_the_proposal_active,
        		COALESCE(party_in_the_proposal.created_at, '1970-01-01 00:00:00') AS party_in_the_proposal_created_at,
        		COALESCE(party_in_the_proposal.updated_at, '1970-01-01 00:00:00') AS party_in_the_proposal_updated_at,
        		
        		COALESCE(organization.id, '00000000-0000-0000-0000-000000000000') AS organization_id,
        		COALESCE(organization.code, 0) AS organization_code,
        		COALESCE(organization.name, '') AS organization_name,
        		COALESCE(organization.nickname, '') AS organization_nickname,
        		COALESCE(organization.acronym, '') AS organization_acronym,
        		COALESCE(organization.type, '') AS organization_type,
        		COALESCE(organization.active, true) AS organization_active,
        		COALESCE(organization.created_at, '1970-01-01 00:00:00') AS organization_created_at,
        		COALESCE(organization.updated_at, '1970-01-01 00:00:00') AS organization_updated_at
			FROM proposition_author
				LEFT JOIN deputy ON deputy.id = proposition_author.deputy_id
				LEFT JOIN party ON party.id = deputy.party_id
				LEFT JOIN party party_in_the_proposal ON party_in_the_proposal.id = proposition_author.party_id
				LEFT JOIN organization ON organization.id = proposition_author.organization_id
			WHERE proposition_author.active = true AND (deputy.active = true AND party.active = true AND
			   party_in_the_proposal.active = true) OR organization.active = true AND proposition_author.proposition_id = $1`
}
