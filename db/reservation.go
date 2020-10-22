package db

import (
	"fmt"
	"ticket-reservation/db/model"
)

type DBReservationInterface interface {
	MakeReservation(userID int, eventID string, amount int) (*model.ReservationDetail, error)
	ViewAllReservations(userID int) ([]*model.ReservationDetail, error)
	CancelReservation(userID int, reservationID string) (string, error)
}

func (pgdb *PostgresqlDB) MakeReservation(userID int, eventID string, amount int) (*model.ReservationDetail, error) {
	ticket, err := pgdb.MemoryDB.UserMakeReservation(userID, eventID, amount)
	if err != nil {
		return nil, err
	}
	return &model.ReservationDetail{
		ReservationID: ticket.ID,
		EventID:       ticket.EventID,
		OrganizerID:   ticket.OwnedBy,
		UserID:        ticket.ReservedBy,
		EventName:     pgdb.MemoryDB.GetEventName(ticket.EventID),
		Tickets:       ticket.Amount,
	}, nil
}

func (pgdb *PostgresqlDB) ViewAllReservations(userID int) ([]*model.ReservationDetail, error) {
	tickets, err := pgdb.MemoryDB.UserViewReservations(userID)
	if err != nil {
		return nil, err
	}
	var res []*model.ReservationDetail
	for _, ticket := range tickets {
		if ticket.Voided {
			continue
		}
		res = append(res, &model.ReservationDetail{
			ReservationID: ticket.ID,
			EventID:       ticket.EventID,
			OrganizerID:   ticket.OwnedBy,
			UserID:        ticket.ReservedBy,
			EventName:     pgdb.MemoryDB.GetEventName(ticket.EventID),
			Tickets:       ticket.Amount,
		})
	}
	return res, nil
}

func (pgdb *PostgresqlDB) CancelReservation(userID int, reservationID string) (string, error) {
	err := pgdb.MemoryDB.UserCancelReservation(userID, reservationID)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Reservation %s Cancelled", reservationID), nil

}
