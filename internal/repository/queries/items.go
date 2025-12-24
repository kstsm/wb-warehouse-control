package queries

const (
	CreateItemQuery = `
		INSERT INTO items (id,
		                   name,
		                   description,
		                   quantity,
		                   price,
		                   created_at,
		                   updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
`

	GetItemByIDQuery = `
		SELECT id,
		       name,
		       description,
		       quantity,
		       price,
		       created_at,
		       updated_at
		FROM items
		WHERE id = $1
`

	GetItemsQuery = `
		SELECT id,
		       name,
		       description,
		       quantity,
		       price,
		       created_at,
		       updated_at
		FROM items
		ORDER BY created_at DESC
`

	UpdateItemQuery = `
		UPDATE items
		SET 
			name = COALESCE($2, name),
			description = COALESCE($3, description),
			quantity = COALESCE($4, quantity),
			price = COALESCE($5, price),
			updated_at = NOW()
		WHERE id = $1
		RETURNING id, name, description, quantity, price, created_at, updated_at
`

	DeleteItemQuery = `
		DELETE FROM items
		WHERE id = $1
		RETURNING id
`

	GetHistoryQuery = `
		SELECT id,
		       item_id,
		       action,
		       user_id,
		       changed_at,
		       old_data,
		       new_data
		FROM items_history
		%s
		%s
`

	GetHistoryCountQuery = `
		SELECT COUNT(*)
		FROM items_history
		%s
`

	GetHistoryByItemIDQuery = `
		SELECT id, item_id, action, user_id, changed_at, old_data, new_data
		FROM items_history
		WHERE item_id = $1
		ORDER BY changed_at DESC
`
)
