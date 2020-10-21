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
