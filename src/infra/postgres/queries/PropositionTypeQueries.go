package queries

type propositionTypeSqlManager struct{}

func PropositionType() *propositionTypeSqlManager {
	return &propositionTypeSqlManager{}
}

type propositionTypeSelectSqlManager struct{}

func (propositionTypeSqlManager) Select() *propositionTypeSelectSqlManager {
	return &propositionTypeSelectSqlManager{}
}

func (propositionTypeSelectSqlManager) All() string {
	return `SELECT id AS proposition_type_id, description AS proposition_type_description,
       			color AS proposition_type_color
			FROM proposition_type
			WHERE active = true
			ORDER BY sort_order`
}
