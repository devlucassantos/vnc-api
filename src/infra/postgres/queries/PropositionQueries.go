package queries

type propositionSqlManager struct{}

func Proposition() *propositionSqlManager {
	return &propositionSqlManager{}
}

type propositionSelectSqlManager struct{}

func (propositionSqlManager) Select() *propositionSelectSqlManager {
	return &propositionSelectSqlManager{}
}

func (propositionSelectSqlManager) TotalNumber() string {
	return `SELECT COUNT(*) FROM proposition`
}

func (propositionSelectSqlManager) All() string {
	return `SELECT COALESCE(id, '00000000-0000-0000-0000-000000000000') AS proposition_id,
    			COALESCE(code, 0) AS proposition_code,
    			COALESCE(original_text_url, '') AS proposition_original_text_url,
       			COALESCE(title, '') AS proposition_title,
    			COALESCE(summary, '') AS proposition_summary,
    			COALESCE(submitted_at, '1970-01-01 00:00:00') AS proposition_submitted_at,
       			COALESCE(active, true) AS proposition_active,
    			COALESCE(created_at, '1970-01-01 00:00:00') AS proposition_created_at,
    			COALESCE(updated_at, '1970-01-01 00:00:00') AS proposition_updated_at
    		FROM proposition WHERE active = true OFFSET $1 LIMIT $2`
}

func (propositionSelectSqlManager) ById() string {
	return `SELECT COALESCE(id, '00000000-0000-0000-0000-000000000000') AS proposition_id,
    			COALESCE(code, 0) AS proposition_code,
    			COALESCE(original_text_url, '') AS proposition_original_text_url,
       			COALESCE(title, '') AS proposition_title,
    			COALESCE(summary, '') AS proposition_summary,
    			COALESCE(submitted_at, '1970-01-01 00:00:00') AS proposition_submitted_at,
       			COALESCE(active, true) AS proposition_active,
    			COALESCE(created_at, '1970-01-01 00:00:00') AS proposition_created_at,
    			COALESCE(updated_at, '1970-01-01 00:00:00') AS proposition_updated_at
    		FROM proposition WHERE active = true AND id = $1`
}
