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

var AllowedActionTypes = map[ActionType]struct{}{
	ActionCreate: {},
	ActionUpdate: {},
	ActionDelete: {},
}

var AllowedRoles = map[Role]struct{}{
	RoleAdmin:   {},
	RoleManager: {},
	RoleViewer:  {},
}
