package db

type DBReservationInterface interface {
	MakeReservation()
	ViewReservation()
	ViewAllReservations()
	EditReservation()
	DeleteReservation()
}
