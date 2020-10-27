// Package model contains The structs in this folder are used to parse data from DB query
// Functions in app package will expect results from DB query to be in these structs
package model

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

type UserWithRoleList struct {
	ID       int      `json:"id"`
	Username string   `json:"username"`
	RoleList []string `json:"roleList"`
}

type EventDetail struct {
	EventID        int    `json:"eventID"`
	OrganizerID    int    `json:"orgID"`
	EventName      string `json:"eventName"`
	Quota          int    `json:"quota"`
	RemainingQuota int    `json:"remainingQuota"`
}

type ReservationDetail struct {
	ReservationID int
	EventID       int
	EventName     string
	OrganizerID   int
	UserID        int
	Tickets       int
}
