package queries

type userSqlManager struct{}

func User() *userSqlManager {
	return &userSqlManager{}
}

func (userSqlManager) Insert() string {
	return `INSERT INTO "user"(first_name, last_name, email, password, activation_code)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING id, created_at`
}

func (userSqlManager) Update() string {
	return `UPDATE "user"
			SET first_name = $1, last_name = $2, email = $3, password = $4, activation_code = $5,
                  updated_at = TIMEZONE('America/Sao_Paulo'::TEXT, NOW())
			WHERE active = true AND id = $6
			RETURNING updated_at`
}

type userSelectSqlManager struct{}

func (userSqlManager) Select() *userSelectSqlManager {
	return &userSelectSqlManager{}
}

func (userSelectSqlManager) ById() string {
	return `SELECT id AS user_id, first_name AS user_first_name, last_name AS user_last_name, email AS user_email,
       			password AS user_password, activation_code AS user_activation_code, created_at AS user_created_at,
       			updated_at AS user_updated_at
			FROM "user"
			WHERE active = true AND id = $1`
}

func (userSelectSqlManager) ByEmail() string {
	return `SELECT id AS user_id, first_name AS user_first_name, last_name AS user_last_name, email AS user_email,
       			password AS user_password, activation_code AS user_activation_code, created_at AS user_created_at,
       			updated_at AS user_updated_at
			FROM "user"
			WHERE active = true AND email = $1`
}
