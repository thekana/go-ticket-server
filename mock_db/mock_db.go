package mock_db

import (
	"ticket-reservation/db"
	"ticket-reservation/db/model"
)

// Must implement the functions we want to mock

type MockDB struct {
	db.DB
	StubViewEventDetail func(eventId int) (*model.EventDetail, error)
	StubEditEvent       func(eventID int, newName string, newQuota int, applicantID int) (*model.EventDetail, error)
	StubDeleteEvent     func(eventId int, applicantID int, admin bool) (string, error)
}

func (d *MockDB) ViewEventDetail(eventId int) (*model.EventDetail, error) {
	return d.StubViewEventDetail(eventId)
}

func (d *MockDB) EditEvent(eventID int, newName string, newQuota int, applicantID int) (*model.EventDetail, error) {
	return d.StubEditEvent(eventID, newName, newQuota, applicantID)
}

func (d *MockDB) DeleteEvent(eventId int, applicantID int, admin bool) (string, error) {
	return d.StubDeleteEvent(eventId, applicantID, admin)
}
