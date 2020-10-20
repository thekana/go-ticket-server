// Package model contains The structs in this folder are used to parse data from DB query
// Functions in app package will expect results from DB query to be in these structs
package model

type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
}

type UserWithRoleList struct {
	ID       int64    `json:"id"`
	Username string   `json:"username"`
	RoleList []string `json:"roleList"`
}

type EventDetail struct {
	EventID     string `json:"event_id"`
	OrganizerID int    `json:"org_id"`
	EventName   string `json:"event_name"`
	Quota       int    `json:quota_remaining`
	SoldAmount  int    `json:sold`
}
