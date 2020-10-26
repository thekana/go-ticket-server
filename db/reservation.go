package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"ticket-reservation/db/model"
)

type DBReservationInterface interface {
	MakeReservation(userID int, eventID int, amount int) (*model.ReservationDetail, error)
	ViewAllReservations(userID int) ([]*model.ReservationDetail, error)
	CancelReservation(userID int, reservationID int) (string, error)
}

func (pgdb *PostgresqlDB) MakeReservation(userID int, eventID int, amount int) (*model.ReservationDetail, error) {
	var res model.ReservationDetail
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
	var sql = `INSERT INTO reservations (user_id,event_id,quota) VALUES ($1,$2,$3) RETURNING id,user_id,event_id,quota;`
	err = tx.QueryRow(context.Background(), sql, userID, eventID, amount).Scan(&res.ReservationID, &res.UserID, &res.EventID, &res.Tickets)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.New("event not found")
		}
		return nil, err
	}
	err = tx.QueryRow(context.Background(), `SELECT name, owner FROM events where id=$1`, eventID).Scan(&res.EventName, &res.OrganizerID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.New("event not found")
		}
		return nil, err
	}
	err = tx.Commit(context.Background())
	return &res, nil
}

func (pgdb *PostgresqlDB) ViewAllReservations(userID int) ([]*model.ReservationDetail, error) {
	var reservations []*model.ReservationDetail
	rows, err := pgdb.DB.Query(context.Background(),
		`SELECT r.*, e.name, e.owner from reservations r
			JOIN events e on (e.id = r.event_id)
			WHERE user_id=$1
			ORDER BY id ASC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var reservation model.ReservationDetail
		err = rows.Scan(&reservation.ReservationID, &reservation.UserID, &reservation.EventID, &reservation.Tickets, nil, &reservation.EventName, &reservation.OrganizerID)
		if err != nil {
			return nil, err
		}
		reservations = append(reservations, &reservation)
	}
	return reservations, nil
}

func (pgdb *PostgresqlDB) CancelReservation(userID int, reservationID int) (string, error) {
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
	var sql = `DELETE from reservations where id=$1 and user_id=$2`
	row, err := tx.Exec(context.Background(), sql, reservationID, userID)
	if err != nil {
		return "", err
	}
	if row.RowsAffected() == 0 {
		return "", errors.New("delete failed")
	}
	err = tx.Commit(context.Background())
	if err != nil {
		return "", errors.Wrap(err, "Unable to commit a transaction")
	}
	return fmt.Sprintf("Reservation %d Cancelled", reservationID), nil
}
