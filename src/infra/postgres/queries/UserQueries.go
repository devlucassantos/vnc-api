package queries

type userSqlManager struct{}

func User() *userSqlManager {
	return &userSqlManager{}
}

func (userSqlManager) Insert() string {
	return `INSERT INTO "user"(first_name, last_name, email, password, active)
			VALUES ($1, $2, $3, $4, true)
			RETURNING id, created_at`
}

type userSelectSqlManager struct{}

func (userSqlManager) Select() *userSelectSqlManager {
	return &userSelectSqlManager{}
}

func (userSelectSqlManager) ById() string {
	return `SELECT id AS user_id, first_name AS user_first_name, last_name AS user_last_name, email AS user_email,
       			password AS user_password, created_at AS user_created_at, updated_at AS user_updated_at
			FROM "user"
			WHERE active = true AND id = $1`
}

func (userSelectSqlManager) ByEmail() string {
	return `SELECT id AS user_id, first_name AS user_first_name, last_name AS user_last_name, email AS user_email,
       			password AS user_password, created_at AS user_created_at, updated_at AS user_updated_at
			FROM "user"
			WHERE active = true AND email = $1`
}
