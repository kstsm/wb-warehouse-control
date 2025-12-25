package validator

type ActionType string

const (
	ActionCreate ActionType = "create"
	ActionUpdate ActionType = "update"
	ActionDelete ActionType = "delete"
)

type Role string

const (
	RoleAdmin   Role = "admin"
	RoleManager Role = "manager"
	RoleViewer  Role = "viewer"
)

//nolint:gochecknoglobals // These are constant maps used for validation
var AllowedActionTypes = map[ActionType]struct{}{
	ActionCreate: {},
	ActionUpdate: {},
	ActionDelete: {},
}

//nolint:gochecknoglobals // These are constant maps used for validation
var AllowedRoles = map[Role]struct{}{
	RoleAdmin:   {},
	RoleManager: {},
	RoleViewer:  {},
}
