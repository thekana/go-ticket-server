package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"net/http"
	customError "ticket-reservation/custom_error"
	"ticket-reservation/db/model"
)

type DBEventInterface interface {
	CreateEvent(ownerId int, eventName string, quota int) (int, error)
	ViewEventDetail(eventId int) (*model.EventDetail, error)
	ViewAllEvents() ([]*model.EventDetail, error)
	OrganizerViewAllEvents(uid int) ([]*model.EventDetail, error)
	EditEvent(eventId int, newName string, newQuota int, applicantID int) (*model.EventDetail, error)
	DeleteEvent(eventId string, applicantID int, admin bool) (string, error)
}

func (pgdb *PostgresqlDB) CreateEvent(ownerId int, eventName string, quota int) (int, error) {
	var eventID int
	tx, err := pgdb.DB.Begin(context.Background())
	if err != nil {
		return 0, errors.Wrap(err, "Unable to make a transaction")
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback(context.Background())
		}
	}()
	err = tx.QueryRow(context.Background(), `INSERT INTO events ("name","quota","owner") values ($1,$2,$3) returning id`, eventName, quota, ownerId).Scan(&eventID)
	if err != nil {
		return 0, errors.Wrap(err, "Unable to create event")
	}
	err = tx.Commit(context.Background())
	if err != nil {
		return 0, errors.Wrap(err, "Unable to commit a transaction")
	}
	return eventID, nil
}

// Viewing a particular event's detail should be a transaction
func (pgdb *PostgresqlDB) ViewEventDetail(eventId int) (*model.EventDetail, error) {
	event := &model.EventDetail{}
	tx, err := pgdb.DB.Begin(context.Background())
	if err != nil {
		return nil, errors.Wrap(err, "Unable to make a transaction")
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback(context.Background())
		}
	}()
	_, err = tx.Exec(context.Background(), "SET TRANSACTION ISOLATION LEVEL REPEATABLE READ")
	if err != nil {
		tx.Rollback(context.Background())
		return nil, errors.Wrap(err, "Unable to set transaction isolation level")
	}
	var sql string = `SELECT id, name, quota, owner, sold from events where id=$1`
	err = tx.QueryRow(context.Background(), sql, eventId).Scan(&event.EventID, &event.EventName, &event.Quota, &event.OrganizerID, &event.SoldAmount)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.Wrap(err, "event not found")
		}
		return nil, err
	}
	err = tx.Commit(context.Background())
	if err != nil {
		return nil, errors.Wrap(err, "Unable to commit a transaction")
	}
	return event, nil
}

// Viewing all events shouldn't be a transaction due to performance concern
func (pgdb *PostgresqlDB) ViewAllEvents() ([]*model.EventDetail, error) {
	// For Admin and Customer
	var events []*model.EventDetail
	var sql string = `SELECT id, name, quota, owner, sold from events`
	rows, err := pgdb.DB.Query(context.Background(), sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var event model.EventDetail
		err = rows.Scan(&event.EventID, &event.EventName, &event.Quota, &event.OrganizerID, &event.SoldAmount)
		if err != nil {
			return nil, err
		}
		events = append(events, &event)
	}
	return events, nil
}
func (pgdb *PostgresqlDB) OrganizerViewAllEvents(uid int) ([]*model.EventDetail, error) {
	// For Organizers
	var events []*model.EventDetail
	var sql string = `SELECT id, name, quota, owner, sold from events where owner=$1`
	rows, err := pgdb.DB.Query(context.Background(), sql, uid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var event model.EventDetail
		err = rows.Scan(&event.EventID, &event.EventName, &event.Quota, &event.OrganizerID, &event.SoldAmount)
		if err != nil {
			return nil, err
		}
		events = append(events, &event)
	}
	return events, nil
}
func (pgdb *PostgresqlDB) EditEvent(eventId int, newName string, newQuota int, applicantID int) (*model.EventDetail, error) {
	return nil, nil
	//eventToEdit, exist := pgdb.MemoryDB.GetEvent(eventId)
	//pgdb.MemoryDB.resourceLock.TryLock(eventToEdit.ID)
	//defer pgdb.MemoryDB.resourceLock.Unlock(eventToEdit.ID)
	//if !exist {
	//	return nil, &customError.UserError{
	//		Code:           customError.UnknownError,
	//		Message:        "Event not found",
	//		HTTPStatusCode: http.StatusBadRequest,
	//	}
	//}
	//if eventToEdit.OrganizerID != applicantID {
	//	return nil, &customError.AuthorizationError{
	//		Code:           customError.Unauthorized,
	//		Message:        "Not Allowed",
	//		HTTPStatusCode: http.StatusUnauthorized,
	//	}
	//}
	//eventToEdit.Name = newName
	//if newQuota < eventToEdit.SoldAmount {
	//	return nil, &customError.UserError{
	//		Code:           customError.UnknownError,
	//		Message:        "New quota must not be less than sold tickets",
	//		HTTPStatusCode: http.StatusBadRequest,
	//	}
	//}
	//eventToEdit.Quota = newQuota
	//
	//return &model.EventDetail{
	//	EventID:     0,
	//	OrganizerID: eventToEdit.OrganizerID,
	//	EventName:   eventToEdit.Name,
	//	Quota:       eventToEdit.Quota - eventToEdit.SoldAmount,
	//	SoldAmount:  eventToEdit.SoldAmount,
	//}, nil
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
