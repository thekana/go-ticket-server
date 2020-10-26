package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"ticket-reservation/db/model"
)

type DBEventInterface interface {
	CreateEvent(ownerId int, eventName string, quota int) (int, error)
	ViewEventDetail(eventId int) (*model.EventDetail, error)
	ViewAllEvents() ([]*model.EventDetail, error)
	OrganizerViewAllEvents(uid int) ([]*model.EventDetail, error)
	EditEvent(eventId int, newName string, newQuota int, applicantID int) (*model.EventDetail, error)
	DeleteEvent(eventId int, applicantID int, admin bool) (string, error)
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
			return nil, errors.New("event not found")
		}
		return nil, err
	}
	err = tx.Commit(context.Background())
	if err != nil {
		return nil, errors.Wrap(err, "Unable to commit a transaction")
	}
	event.Quota = event.Quota - event.SoldAmount
	return event, nil
}

// Viewing all events shouldn't be a transaction due to performance concern
func (pgdb *PostgresqlDB) ViewAllEvents() ([]*model.EventDetail, error) {
	// For Admin and Customer
	var events []*model.EventDetail
	var sql string = `SELECT id, name, quota, owner, sold from events ORDER BY id ASC`
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
		event.Quota = event.Quota - event.SoldAmount
		events = append(events, &event)
	}
	return events, nil
}
func (pgdb *PostgresqlDB) OrganizerViewAllEvents(uid int) ([]*model.EventDetail, error) {
	// For Organizers
	var events []*model.EventDetail
	var sql string = `SELECT id, name, quota, owner, sold from events where owner=$1 ORDER BY id ASC`
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
		event.Quota = event.Quota - event.SoldAmount
		events = append(events, &event)
	}
	return events, nil
}
func (pgdb *PostgresqlDB) EditEvent(eventId int, newName string, newQuota int, applicantID int) (*model.EventDetail, error) {
	event := &model.EventDetail{}
	var trueOwnerID int
	tx, err := pgdb.DB.Begin(context.Background())
	if err != nil {
		return nil, errors.Wrap(err, "Unable to make a transaction")
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback(context.Background())
		}
	}()
	_, err = tx.Exec(context.Background(), "SET TRANSACTION ISOLATION LEVEL SERIALIZABLE")
	if err != nil {
		tx.Rollback(context.Background())
		return nil, errors.Wrap(err, "Unable to set transaction isolation level")
	}
	var sql string = `select owner from events where id=$1`
	err = tx.QueryRow(context.Background(), sql, eventId).Scan(&trueOwnerID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.New("event not found")
		}
		return nil, err
	}
	if trueOwnerID != applicantID {
		return nil, errors.New("Not Authorized")
	}
	sql = `UPDATE events SET quota=$1,name=$2,updated_at=now() where id=$3 returning id, name, quota, owner, sold`
	err = tx.QueryRow(context.Background(), sql, newQuota, newName, eventId).Scan(&event.EventID, &event.EventName, &event.Quota, &event.OrganizerID, &event.SoldAmount)
	if err != nil {
		return nil, err
	}
	err = tx.Commit(context.Background())
	if err != nil {
		return nil, errors.Wrap(err, "Unable to commit a transaction")
	}
	event.Quota -= event.SoldAmount
	return event, nil
}
func (pgdb *PostgresqlDB) DeleteEvent(eventId int, applicantID int, admin bool) (string, error) {
	tx, err := pgdb.DB.Begin(context.Background())
	if err != nil {
		return "", errors.Wrap(err, "Unable to make a transaction")
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback(context.Background())
		}
	}()
	_, err = tx.Exec(context.Background(), "SET TRANSACTION ISOLATION LEVEL SERIALIZABLE")
	if err != nil {
		tx.Rollback(context.Background())
		return "", errors.Wrap(err, "Unable to set transaction isolation level")
	}
	var sql string = `select owner from events where id=$1`
	var trueOwnerID int
	err = tx.QueryRow(context.Background(), sql, eventId).Scan(&trueOwnerID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", errors.New("event not found")
		}
		return "", err
	}
	if !admin && (trueOwnerID != applicantID) {
		return "", errors.New("Not Authorized")
	}
	sql = `DELETE from events where id=$1`
	_, err = tx.Exec(context.Background(), sql, eventId)
	if err != nil {
		return "", errors.Wrap(err, "Cannot delete event")
	}
	err = tx.Commit(context.Background())
	if err != nil {
		return "", errors.Wrap(err, "Unable to commit a transaction")
	}
	return fmt.Sprintf("Event id %d was deleted by user %d", eventId, applicantID), nil
}
