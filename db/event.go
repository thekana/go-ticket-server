package db

import (
	"github.com/pkg/errors"
	"ticket-reservation/db/model"
)

type DBEventInterface interface {
	CreateEvent(ownerId int, eventName string, quota int) (string, error)
	ViewEventDetail(eventId string) (*model.EventDetail, error)
	ViewAllEvents() ([]*model.EventDetail, error)
	OrganizerViewAllEvents(uid int) ([]*model.EventDetail, error)
	EditEvent(eventId string, newName string, newQuota int, applicantID int) (*model.EventDetail, error)
	DeleteEvent(eventId string) (interface{}, error)
}

// TODO: Errors handling
func (pgdb *PostgresqlDB) CreateEvent(ownerId int, eventName string, quota int) (string, error) {
	event := NewEvent(ownerId, eventName, quota)
	// Add event to memory
	pgdb.MemoryDB.AddEventToSystem(event)
	// Also add event to user
	pgdb.MemoryDB.UserMap[ownerId].AddEvent(event)
	return event.Id, nil
}

// TODO: Errors handling
func (pgdb *PostgresqlDB) ViewEventDetail(eventId string) (*model.EventDetail, error) {
	// Event detail is open anyone can view it
	thisEvent := pgdb.MemoryDB.EventMap[eventId]
	return &model.EventDetail{
		EventID:     thisEvent.Id,
		OrganizerID: thisEvent.OrganizerID,
		EventName:   thisEvent.Name,
		Quota:       thisEvent.Quota,
		SoldAmount:  thisEvent.SoldAmount,
	}, nil
}
func (pgdb *PostgresqlDB) ViewAllEvents() ([]*model.EventDetail, error) {
	// For Admin and Customer
	var event []*model.EventDetail
	for _, e := range pgdb.MemoryDB.EventList {
		event = append(event, &model.EventDetail{
			EventID:     e.Id,
			OrganizerID: e.OrganizerID,
			EventName:   e.Name,
			Quota:       e.Quota,
			SoldAmount:  e.SoldAmount,
		})
	}
	return event, nil
}

func (pgdb *PostgresqlDB) OrganizerViewAllEvents(uid int) ([]*model.EventDetail, error) {
	// For Admin and Customer
	var event []*model.EventDetail
	for _, e := range pgdb.MemoryDB.UserMap[uid].Events {
		event = append(event, &model.EventDetail{
			EventID:     e.Id,
			OrganizerID: e.OrganizerID,
			EventName:   e.Name,
			Quota:       e.Quota,
			SoldAmount:  e.SoldAmount,
		})
	}
	return event, nil
}

func (pgdb *PostgresqlDB) EditEvent(eventId string, newName string, newQuota int, applicantID int) (*model.EventDetail, error) {
	eventToEdit, exist := pgdb.MemoryDB.EventMap[eventId]
	if !exist {
		return nil, errors.New("Event not in system")
	}
	if eventToEdit.OrganizerID != applicantID {
		return nil, errors.New("Not Authorized")
	}
	eventToEdit.Name = newName
	if newQuota < eventToEdit.SoldAmount {
		return nil, errors.New("New Quota must be more than sold quota")
	}
	eventToEdit.Quota = newQuota

	return &model.EventDetail{
		EventID:     eventToEdit.Id,
		OrganizerID: eventToEdit.OrganizerID,
		EventName:   eventToEdit.Name,
		Quota:       eventToEdit.Quota,
		SoldAmount:  eventToEdit.SoldAmount,
	}, nil
}
func (pgdb *PostgresqlDB) DeleteEvent(eventId string) (interface{}, error) {
	return nil, nil
}
