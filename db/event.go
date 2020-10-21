package db

import (
	"fmt"
	"net/http"
	customError "ticket-reservation/custom_error"
	"ticket-reservation/db/model"
)

type DBEventInterface interface {
	CreateEvent(ownerId int, eventName string, quota int) (string, error)
	ViewEventDetail(eventId string) (*model.EventDetail, error)
	ViewAllEvents() ([]*model.EventDetail, error)
	OrganizerViewAllEvents(uid int) ([]*model.EventDetail, error)
	EditEvent(eventId string, newName string, newQuota int, applicantID int) (*model.EventDetail, error)
	DeleteEvent(eventId string, applicantID int, admin bool) (string, error)
}

// TODO: Errors handling
func (pgdb *PostgresqlDB) CreateEvent(ownerId int, eventName string, quota int) (string, error) {
	event := NewEvent(ownerId, eventName, quota)
	// Add event to memory
	pgdb.MemoryDB.AddEventToSystem(event)
	return event.Id, nil
}

// TODO: Errors handling
func (pgdb *PostgresqlDB) ViewEventDetail(eventId string) (*model.EventDetail, error) {
	// Event detail is open anyone can view it
	thisEvent, found := pgdb.MemoryDB.GetEvent(eventId)
	if !found || thisEvent.Deleted {
		return nil, &customError.UserError{
			Code:           customError.UnknownError,
			Message:        "Event not found",
			HTTPStatusCode: http.StatusBadRequest,
		}
	}
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
	for _, e := range pgdb.MemoryDB.GetAllEvents() {
		if e.Deleted {
			continue
		}
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
	for _, e := range pgdb.MemoryDB.GetEventsOwnedByUser(uid) {
		if e.Deleted {
			continue
		}
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
	eventToEdit, exist := pgdb.MemoryDB.GetEvent(eventId)
	if !exist {
		return nil, &customError.UserError{
			Code:           customError.UnknownError,
			Message:        "Event not found",
			HTTPStatusCode: http.StatusBadRequest,
		}
	}
	if eventToEdit.OrganizerID != applicantID {
		return nil, &customError.AuthorizationError{
			Code:           customError.Unauthorized,
			Message:        "Not Allowed",
			HTTPStatusCode: http.StatusUnauthorized,
		}
	}
	eventToEdit.Name = newName
	if newQuota < eventToEdit.SoldAmount {
		return nil, &customError.UserError{
			Code:           customError.UnknownError,
			Message:        "New quota must not be less than sold tickets",
			HTTPStatusCode: http.StatusBadRequest,
		}
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
func (pgdb *PostgresqlDB) DeleteEvent(eventId string, applicantID int, admin bool) (string, error) {
	eventToDelete, exist := pgdb.MemoryDB.GetEvent(eventId)
	if !exist || eventToDelete.Deleted {
		return "Event not in system", &customError.UserError{
			Code:           customError.UnknownError,
			Message:        "Event not found",
			HTTPStatusCode: http.StatusBadRequest,
		}
	}
	if admin {
		pgdb.MemoryDB.DeleteEvent(eventId)
	} else {
		if eventToDelete.OrganizerID == applicantID {
			pgdb.MemoryDB.DeleteEvent(eventId)
		} else {
			return "Not allowed", &customError.AuthorizationError{
				Code:           customError.Unauthorized,
				Message:        "Not Allowed",
				HTTPStatusCode: http.StatusUnauthorized,
			}
		}
	}
	return fmt.Sprintf("%s deleted by user %d", eventId, applicantID), nil
}
