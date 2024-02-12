package queries

type organizationSqlManager struct{}

func Organization() *organizationSqlManager {
	return &organizationSqlManager{}
}

type organizationSelectSqlManager struct{}

func (organizationSqlManager) Select() *organizationSelectSqlManager {
	return &organizationSelectSqlManager{}
}

func (organizationSelectSqlManager) All() string {
	return `SELECT COALESCE(organization.id, '00000000-0000-0000-0000-000000000000') AS organization_id,
        		COALESCE(organization.code, 0) AS organization_code,
        		COALESCE(organization.name, '') AS organization_name,
        		COALESCE(organization.nickname, '') AS organization_nickname,
        		COALESCE(organization.acronym, '') AS organization_acronym,
        		COALESCE(organization.type, '') AS organization_type,
        		COALESCE(organization.active, true) AS organization_active,
        		COALESCE(organization.created_at, '1970-01-01 00:00:00') AS organization_created_at,
        		COALESCE(organization.updated_at, '1970-01-01 00:00:00') AS organization_updated_at
    		FROM organization WHERE organization.active = true
    		ORDER BY organization.name`
}
