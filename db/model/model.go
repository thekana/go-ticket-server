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
	ReservationID int    `json:"reservationID"`
	EventID       int    `json:"eventID"`
	EventName     string `json:"eventName"`
	OrganizerID   int    `json:"organizerID"`
	UserID        int    `json:"userID"`
	Tickets       int    `json:"tickets"`
}

type ReservationTicket struct {
	ReservationID int `json:"reservationID"`
	EventID       int `json:"eventID"`
	UserID        int `json:"userID"`
	Tickets       int `json:"tickets"`
}

type DeletedTicket struct {
	ReservationID int `json:"reservationID"`
	EventID       int `json:"eventID"`
	Amount        int `json:"amount"`
}

type ReservationRequest struct {
	EventID int
	UserID  int
	Amount  int
}
