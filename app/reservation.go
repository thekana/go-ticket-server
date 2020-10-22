package app

import "ticket-reservation/db/model"

type MakeReservationParams struct {
	AuthToken string `json:"authToken" validate:"required"`
	EventID   string `json:"eventID" validate:"required"`
	Amount    int    `json:"amount" validate:"required"`
}

type MakeReservationResult struct {
	Ticket *model.ReservationDetail `json:"ticket"`
}

type ViewReservationsParams struct {
	AuthToken string `json:"authToken" validate:"required"`
}

type ViewReservationsResult struct {
	Tickets []*model.ReservationDetail `json:"tickets"`
}

type CancelReservationParams struct {
	AuthToken     string `json:"authToken" validate:"required"`
	ReservationID string `json:"reservationId" validate:"required"`
}

type CancelReservationResult struct {
	Message string `json:"message"`
}

func (ctx *Context) MakeReservation(params MakeReservationParams) (*MakeReservationResult, error) {
	logger := ctx.getLogger()

	if err := validateInput(params); err != nil {
		logger.Errorf("validateInput error : %s", err)
		return nil, err
	}
	authRes, err := ctx.authorizeUser(params.AuthToken, []model.Role{model.Customer})
	if err != nil {
		return nil, err
	}
	ticket, err := ctx.DB.MakeReservation(int(authRes.User.ID), params.EventID, params.Amount)
	if err != nil {
		return nil, err
	}
	return &MakeReservationResult{Ticket: ticket}, nil
}

func (ctx *Context) ViewReservations(params ViewReservationsParams) (*ViewReservationsResult, error) {
	logger := ctx.getLogger()

	if err := validateInput(params); err != nil {
		logger.Errorf("validateInput error : %s", err)
		return nil, err
	}
	authRes, err := ctx.authorizeUser(params.AuthToken, []model.Role{model.Customer, model.Organizer, model.Admin})
	if err != nil {
		return nil, err
	}
	tickets, err := ctx.DB.ViewAllReservations(int(authRes.User.ID))
	if err != nil {
		return nil, err
	}
	return &ViewReservationsResult{Tickets: tickets}, nil
}

func (ctx *Context) CancelReservation(params CancelReservationParams) (*CancelReservationResult, error) {
	logger := ctx.getLogger()

	if err := validateInput(params); err != nil {
		logger.Errorf("validateInput error : %s", err)
		return nil, err
	}
	authRes, err := ctx.authorizeUser(params.AuthToken, []model.Role{model.Customer})

	if err != nil {
		return nil, err
	}
	message, err := ctx.DB.CancelReservation(int(authRes.User.ID), params.ReservationID)

	if err != nil {
		return nil, err
	}
	return &CancelReservationResult{Message: message}, nil
}
