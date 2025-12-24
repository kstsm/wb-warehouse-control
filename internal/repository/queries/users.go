package queries

const (
	GetUserByNameQuery = `
		SELECT id,
		       name,
		       role,
		       created_at,
		       updated_at
		FROM users
		WHERE name = $1
`

	GetOrCreateUserQuery = `
	INSERT INTO users (id,
	                   name,
	                   role,
	                   created_at,
	                   updated_at)
	VALUES ($1, $2, $3, $4, $5)
	ON CONFLICT (name) DO UPDATE
	SET name = users.name
	RETURNING id, name, role, created_at, updated_at;
`
)
