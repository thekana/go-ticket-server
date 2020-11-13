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
	ViewAllEvents(isOrganizer bool, orgID int) ([]*model.EventDetail, error)
	EditEvent(eventID int, newName string, newQuota int, applicantID int) (*model.EventDetail, error)
	DeleteEvent(eventId int, applicantID int, admin bool) (string, error)
	RefreshEventQuotasFromEntryInReservationsTable() error
	ReclaimEventQuotas(cancelledTickets map[int]int) error
}

func (pgdb *PostgresqlDB) CreateEvent(ownerId int, eventName string, quota int) (int, error) {
	var eventID int
	var sql = `INSERT INTO events ("name","quota","remaining_quota","owner") values ($1,$2,$2,$3) returning id`
	err := pgdb.DB.QueryRow(context.Background(), sql, eventName, quota, ownerId).Scan(&eventID)
	if err != nil {
		return 0, errors.Wrap(err, "Unable to create event")
	}
	return eventID, nil
}
func (pgdb *PostgresqlDB) ViewEventDetail(eventId int) (*model.EventDetail, error) {
	event := model.EventDetail{}
	var sql = `SELECT * from events where id=$1`
	err := pgdb.DB.QueryRow(context.Background(),
		sql, eventId).Scan(&event.EventID,
		&event.EventName,
		&event.OrganizerID,
		&event.Quota,
		&event.RemainingQuota, nil, nil)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.New("Event not found")
		}
		return nil, err
	}
	return &event, nil
}
func (pgdb *PostgresqlDB) ViewAllEvents(isOrganizer bool, orgID int) ([]*model.EventDetail, error) {
	var events []*model.EventDetail
	var sql = `SELECT * from events ORDER BY id ASC`
	if isOrganizer {
		sql = fmt.Sprintf("SELECT * from events where owner=%d ORDER BY id ASC", orgID)
	}
	rows, err := pgdb.DB.Query(context.Background(), sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var event model.EventDetail
		err = rows.Scan(&event.EventID,
			&event.EventName,
			&event.OrganizerID,
			&event.Quota,
			&event.RemainingQuota, nil, nil)
		if err != nil {
			return nil, err
		}
		events = append(events, &event)
	}
	return events, nil
}
func (pgdb *PostgresqlDB) EditEvent(eventID int, newName string, newQuota int, applicantID int) (*model.EventDetail, error) {
	event := &model.EventDetail{}
	var trueOwnerID int
	var oldQuota int
	var oldRemainingQuota int
	tx, err := pgdb.DB.BeginTx(context.Background(), pgx.TxOptions{
		IsoLevel: pgx.RepeatableRead,
	})
	if err != nil {
		return nil, errors.Wrap(err, "Unable to make a transaction")
	}
	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback(context.Background())
		} else if err != nil {
			_ = tx.Rollback(context.Background())
		}
	}()
	var sql = `select quota, remaining_quota, owner from events where id=$1`
	err = tx.QueryRow(context.Background(), sql, eventID).Scan(&oldQuota, &oldRemainingQuota, &trueOwnerID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.New("Event not found")
		}
		return nil, err
	}
	if trueOwnerID != applicantID {
		return nil, errors.New("Not Authorized")
	}
	sold := oldQuota - oldRemainingQuota
	if newQuota < sold {
		_ = tx.Rollback(context.Background())
		return nil, errors.New("New quota must be more than sold tickets")
	}
	sql = `UPDATE events SET quota=$1,name=$2,remaining_quota=$3,updated_at=now() WHERE id=$4 RETURNING *`
	err = tx.QueryRow(context.Background(), sql,
		newQuota, newName, newQuota-sold,
		eventID).Scan(&event.EventID,
		&event.EventName,
		&event.OrganizerID,
		&event.Quota,
		&event.RemainingQuota, nil, nil)
	if err != nil {
		return nil, err
	}
	err = tx.Commit(context.Background())
	if err != nil {
		return nil, errors.Wrap(err, "Unable to commit a transaction")
	}
	return event, nil
}
func (pgdb *PostgresqlDB) DeleteEvent(eventId int, applicantID int, admin bool) (string, error) {
	tx, err := pgdb.DB.BeginTx(context.Background(), pgx.TxOptions{
		IsoLevel: pgx.RepeatableRead,
	})
	if err != nil {
		return "", errors.Wrap(err, "Unable to make a transaction")
	}
	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback(context.Background())
		} else if err != nil {
			_ = tx.Rollback(context.Background())
		}
	}()
	var sql = `select owner from events where id=$1`
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
func (pgdb *PostgresqlDB) RefreshEventQuotasFromEntryInReservationsTable() error {
	tx, err := pgdb.DB.BeginTx(context.Background(), pgx.TxOptions{
		IsoLevel: pgx.RepeatableRead,
	})
	if err != nil {
		return errors.Wrap(err, "Unable to make a transaction")
	}
	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback(context.Background())
		} else if err != nil {
			_ = tx.Rollback(context.Background())
		}
	}()
	sql := "select event_id, sum(quota) from reservations group by event_id"
	rows, err := pgdb.DB.Query(context.Background(), sql)
	defer rows.Close()
	if err != nil {
		return err
	}
	for rows.Next() {
		var eventID int
		var quotaBought int
		err = rows.Scan(&eventID, &quotaBought)
		if err != nil {
			return err
		}
		_, _ = pgdb.DB.Exec(context.Background(), "update events set remaining_quota = quota - $1 where id = $2", quotaBought, eventID)
	}
	err = tx.Commit(context.Background())
	if err != nil {
		return errors.Wrap(err, "Unable to commit a transaction")
	}
	return nil
}
func (pgdb *PostgresqlDB) ReclaimEventQuotas(cancelledTickets map[int]int) error {
	tx, err := pgdb.DB.BeginTx(context.Background(), pgx.TxOptions{
		IsoLevel: pgx.RepeatableRead,
	})
	if err != nil {
		return errors.Wrap(err, "Unable to make a transaction")
	}
	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback(context.Background())
		} else if err != nil {
			_ = tx.Rollback(context.Background())
		}
	}()
	var sql = `UPDATE events SET remaining_quota=remaining_quota+$1 WHERE id=$2`
	for id, val := range cancelledTickets {
		_, err := tx.Exec(context.Background(), sql, val, id)
		if err != nil {
			return err
		}
		//if row.RowsAffected() == 0 {
		//	err = errors.New("Cannot reclaim quota")
		//	return err
		//}
	}
	err = tx.Commit(context.Background())
	if err != nil {
		return errors.Wrap(err, "Unable to commit a transaction")
	}

	return nil
}
