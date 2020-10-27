package app

import (
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"net/http"
	customError "ticket-reservation/custom_error"
	"ticket-reservation/db/model"
)

type MakeReservationParams struct {
	AuthToken string `json:"authToken" validate:"required"`
	EventID   int    `json:"eventID" validate:"required"`
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
	ReservationID int    `json:"reservationId" validate:"required"`
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
		return nil, &customError.AuthorizationError{
			Code:           0,
			Message:        err.Error(),
			HTTPStatusCode: http.StatusForbidden,
		}
	}
	var ticket *model.ReservationDetail
	// Retry 10 times
	for i := 0; i < 10; i++ {
		ticket, err = ctx.DB.MakeReservation(authRes.User.ID, params.EventID, params.Amount)
		if err != nil {
			if pgErr, ok := err.(*pgconn.PgError); ok {
				if pgErr.Code == pgerrcode.SerializationFailure {
					fmt.Printf("\n[%d -> retrying user %d reserves event %d for %d]\n", i, authRes.User.ID, params.EventID, params.Amount)
					continue
				}
			}
		}
		break
	}
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == pgerrcode.SerializationFailure {
				return nil, &customError.InternalError{
					Code:    69,
					Message: "CONCURRENT ERROR",
				}
			}
		}
		return nil, &customError.UserError{
			Code:           0,
			Message:        err.Error(),
			HTTPStatusCode: http.StatusBadRequest,
		}
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
		return nil, &customError.AuthorizationError{
			Code:           0,
			Message:        err.Error(),
			HTTPStatusCode: http.StatusForbidden,
		}
	}
	tickets, err := ctx.DB.ViewAllReservations(int(authRes.User.ID))
	if err != nil {
		return nil, &customError.UserError{
			Code:           0,
			Message:        err.Error(),
			HTTPStatusCode: http.StatusBadRequest,
		}
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
		return nil, &customError.AuthorizationError{
			Code:           0,
			Message:        err.Error(),
			HTTPStatusCode: http.StatusForbidden,
		}
	}
	message, err := ctx.DB.CancelReservation(int(authRes.User.ID), params.ReservationID)

	if err != nil {
		return nil, &customError.UserError{
			Code:           0,
			Message:        err.Error(),
			HTTPStatusCode: http.StatusBadRequest,
		}
	}
	return &CancelReservationResult{Message: message}, nil
}
