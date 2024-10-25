package models

type UserPermissionAPI struct {
	Result UserPermissionResult `json:"result"`
}

type UserPermissionResult struct {
	StatusCode int                      `json:"StatusCode" db:"-"`
	Message    string                   `json:"Message" db:"-"`
	Data       []*UserPermissionApiData `json:"Data" db:"-"`
}

type UserPermissionApiData struct {
	UserID       uint                `json:"UserID" db:"user_id"`
	PermissionID int                 `json:"PermissionID" db:"permission_id"`
	Permission   PermissionAPIResult `json:"Permission" db:"-"`
}
