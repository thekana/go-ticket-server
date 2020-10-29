package app

import (
	"github.com/jackc/pgerrcode"
	"net/http"
	customError "ticket-reservation/custom_error"
	"ticket-reservation/db/model"
	"time"
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

type ReservationQueueResult struct {
	ticket *model.ReservationDetail
	err    error
}

type ReservationQueueElem struct {
	UserID  int
	EventID int
	Amount  int
	c       chan ReservationQueueResult
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
	elem := &ReservationQueueElem{
		UserID:  authRes.User.ID,
		EventID: params.EventID,
		Amount:  params.Amount,
		c:       make(chan ReservationQueueResult, 1),
	}
	ctx.Queue.PushBack(elem)
	var result ReservationQueueResult
	defer close(elem.c)
	select {
	case b := <-elem.c:
		result = b
	case <-time.After(1 * time.Second):
		return nil, &customError.InternalError{
			Code:    100,
			Message: "Time out",
		}
	}
	if result.err != nil {
		if checkPostgresErrorCode(result.err, pgerrcode.SerializationFailure) {
			return nil, &customError.InternalError{
				Code:    69,
				Message: "CONCURRENT ERROR",
			}
		}
		if checkPostgresErrorCode(result.err, pgerrcode.CheckViolation) {
			return nil, &customError.UserError{
				Code:           9,
				Message:        "Not enough quota",
				HTTPStatusCode: http.StatusBadRequest,
			}
		}
		return nil, &customError.UserError{
			Code:           0,
			Message:        result.err.Error(),
			HTTPStatusCode: http.StatusBadRequest,
		}
	}

	return &MakeReservationResult{Ticket: result.ticket}, nil
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
