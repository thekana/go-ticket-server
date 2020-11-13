package db

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"ticket-reservation/db/model"
)

type DBReservationInterface interface {
	MakeReservation(userID int, eventID int, amount int) (*model.ReservationTicket, error)
	ViewAllReservations(userID int) ([]*model.ReservationDetail, error)
	CancelReservationBatch(userID int, reservationIDs []int) ([]*model.DeletedTicket, map[int]int, error)
	MakeReservationBatch(jobs []*model.ReservationRequest, remainingQuotaMap map[int]int) ([]*model.ReservationTicket, error)
}

func (pgdb *PostgresqlDB) MakeReservationBatch(jobs []*model.ReservationRequest, remainingQuotaMap map[int]int) ([]*model.ReservationTicket, error) {
	var results []*model.ReservationTicket
	var data model.ReservationTicket

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
	// Update batch into events FIXME:
	//for k, v := range remainingQuotaMap {
	//	_, err = tx.Exec(context.Background(), `UPDATE events SET remaining_quota=remaining_quota-$1 WHERE id=$2`, v, k)
	//}
	err = tx.Commit(context.Background())
	if err != nil {
		return nil, errors.Wrap(err, "Unable to commit a transaction")
	}
	return results, nil
}

func (pgdb *PostgresqlDB) MakeReservation(userID int, eventID int, amount int) (*model.ReservationTicket, error) {
	var res model.ReservationTicket
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
	_, err = tx.Exec(context.Background(), sql, amount, eventID)
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

func (pgdb *PostgresqlDB) CancelReservationBatch(userID int, reservationIDs []int) ([]*model.DeletedTicket, map[int]int, error) {
	tx, err := pgdb.DB.BeginTx(context.Background(), pgx.TxOptions{
		IsoLevel: pgx.RepeatableRead,
	})
	if err != nil {
		return nil, nil, errors.Wrap(err, "Unable to make a transaction")
	}
	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback(context.Background())
		} else if err != nil {
			_ = tx.Rollback(context.Background())
		}
	}()
	var sql = `DELETE from reservations where id=$1 and user_id=$2 RETURNING event_id, quota`
	var deletedQuota int
	var eventID int
	var deletedTickets []*model.DeletedTicket
	var quotaToReclaim = make(map[int]int)
	for _, id := range reservationIDs {
		err = tx.QueryRow(context.Background(), sql, id, userID).Scan(&eventID, &deletedQuota)
		if err != nil {
			if err == pgx.ErrNoRows {
				continue
			}
			return nil, nil, err
		}
		deletedTickets = append(deletedTickets, &model.DeletedTicket{
			ReservationID: id,
			EventID:       eventID,
			Amount:        deletedQuota,
		})
		quotaToReclaim[eventID] += deletedQuota
	}
	err = tx.Commit(context.Background())
	if err != nil {
		return nil, nil, errors.Wrap(err, "Unable to commit a transaction")
	}

	return deletedTickets, quotaToReclaim, nil
}
