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
	MakeReservationBatch(jobs []*model.ReservationRequest, remainingQuotaMap map[int]int) ([]*model.ReservationDetail, error)
}

func (pgdb *PostgresqlDB) MakeReservationBatch(jobs []*model.ReservationRequest, remainingQuotaMap map[int]int) ([]*model.ReservationDetail, error) {
	var results []*model.ReservationDetail
	var data model.ReservationDetail

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
	// Insert batch into reservations
	for _, job := range jobs {
		err = tx.QueryRow(context.Background(),
			`INSERT INTO reservations (user_id,event_id,quota) VALUES ($1,$2,$3) RETURNING id,user_id,event_id,quota;`,
			job.UserID, job.EventID, job.Amount).Scan(&data.ReservationID, &data.UserID, &data.EventID, &data.Tickets)
		if err != nil {
			return nil, err
		}
		results = append(results, &data)
	}
	// Update batch into events
	for k, v := range remainingQuotaMap {
		_, err = tx.Exec(context.Background(), `UPDATE events SET remaining_quota=remaining_quota-$1 WHERE id=$2`, v, k)
	}
	err = tx.Commit(context.Background())
	if err != nil {
		return nil, errors.Wrap(err, "Unable to commit a transaction")
	}
	return results, nil
}

func (pgdb *PostgresqlDB) MakeReservation(userID int, eventID int, amount int) (*model.ReservationDetail, error) {
	var res model.ReservationDetail
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
	// Insert new reservation
	var sql = `INSERT INTO reservations (user_id,event_id,quota) VALUES ($1,$2,$3) RETURNING id,user_id,event_id,quota;`
	err = tx.QueryRow(context.Background(), sql, userID, eventID, amount).Scan(&res.ReservationID, &res.UserID, &res.EventID, &res.Tickets)
	if err != nil {
		return nil, err
	}
	// Deduct quota
	sql = `UPDATE events SET remaining_quota=remaining_quota-$1 WHERE id=$2 RETURNING name, owner`
	err = tx.QueryRow(context.Background(), sql, amount, eventID).Scan(&res.EventName, &res.OrganizerID)
	if err != nil {
		return nil, err
	}
	err = tx.Commit(context.Background())
	if err != nil {
		return nil, errors.Wrap(err, "Unable to commit a transaction")
	}
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
	var deletedQuota int
	var eventID int
	var sql = `DELETE from reservations where id=$1 and user_id=$2 RETURNING event_id, quota`
	err = tx.QueryRow(context.Background(), sql, reservationID, userID).Scan(&eventID, &deletedQuota)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", errors.Wrap(err, "Invalid ID")
		}
		return "", err
	}
	// Add back quota
	sql = `UPDATE events SET remaining_quota=remaining_quota+$1 WHERE id=$2`
	row, err := tx.Exec(context.Background(), sql, deletedQuota, eventID)
	if err != nil {
		return "", err
	}
	if row.RowsAffected() == 0 {
		err = errors.New("Cannot reclaim quota")
		return "", err
	}
	err = tx.Commit(context.Background())
	if err != nil {
		return "", errors.Wrap(err, "Unable to commit a transaction")
	}
	return fmt.Sprintf("Reservation %d Cancelled", reservationID), nil
}
