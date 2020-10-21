package db

import "ticket-reservation/db/model"

type DBReservationInterface interface {
	MakeReservation(userID int, eventID string, amount int) (*model.ReservationDetail, error)
	//ViewReservation()
	//ViewAllReservations()
	//EditReservation()
	//DeleteReservation()
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
		EventName:     ticket.EventName,
		Tickets:       ticket.Amount,
	}, nil
}

//func (pgdb *PostgresqlDB)
//func (pgdb *PostgresqlDB)
//func (pgdb *PostgresqlDB)
//func (pgdb *PostgresqlDB)
