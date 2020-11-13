package app

import (
	"github.com/jackc/pgerrcode"
	"net/http"
	customError "ticket-reservation/custom_error"
	"ticket-reservation/db/model"
)

type MakeReservationParams struct {
	AuthToken string
	EventID   int `json:"eventID" validate:"required"`
	Amount    int `json:"amount" validate:"required"`
}

type MakeReservationResult struct {
	Ticket *model.ReservationTicket `json:"ticket"`
}

type ViewReservationsParams struct {
	AuthToken string
}

type ViewReservationsResult struct {
	Tickets []*model.ReservationDetail `json:"tickets"`
}

type CancelReservationParams struct {
	AuthToken      string
	ReservationIDs []int `json:"reservationId" validate:"required"`
}

type CancelReservationResults struct {
	DeletedTickets []*model.DeletedTicket `json:"deleted" validate:"required"`
}

type ReservationQueueResult struct {
	ticket *model.ReservationTicket
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
			Code:           customError.Unauthorized,
			Message:        err.Error(),
			HTTPStatusCode: http.StatusUnauthorized,
		}
	}
	elem := &ReservationQueueElem{
		UserID:  authRes.User.ID,
		EventID: params.EventID,
		Amount:  params.Amount,
		c:       make(chan ReservationQueueResult, 1),
	}
	ctx.My.QueueChan <- elem
	var result ReservationQueueResult
	defer close(elem.c)
	select {
	case b := <-elem.c:
		result = b
	}
	if result.err != nil {
		if checkPostgresErrorCode(result.err, pgerrcode.SerializationFailure) {
			return nil, &customError.InternalError{
				Code:    customError.ConcurrencyIssue,
				Message: "DB Concurrency ERROR",
			}
		}
		if checkPostgresErrorCode(result.err, pgerrcode.CheckViolation) {
			return nil, &customError.UserError{
				Code:           customError.InsufficientQuota,
				Message:        "Not enough quota",
				HTTPStatusCode: http.StatusBadRequest,
			}
		}
		return nil, &customError.UserError{
			Code:           customError.UnknownError,
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
	authRes, err := ctx.authorizeUser(params.AuthToken, []model.Role{model.Customer, model.Admin})
	if err != nil {
		return nil, err
	}
	tickets, err := ctx.DB.ViewAllReservations(authRes.User.ID)
	if err != nil {
		return nil, &customError.InternalError{
			Code:    customError.DBError,
			Message: err.Error(),
		}
	}
	return &ViewReservationsResult{Tickets: tickets}, nil
}

func (ctx *Context) CancelReservation(params CancelReservationParams) (*CancelReservationResults, error) {
	logger := ctx.getLogger()

	if err := validateInput(params); err != nil {
		logger.Errorf("validateInput error : %s", err)
		return nil, err
	}
	authRes, err := ctx.authorizeUser(params.AuthToken, []model.Role{model.Customer})

	if err != nil {
		return nil, err
	}
	deletedTickets, quotaToReclaims, err := ctx.DB.CancelReservationBatch(authRes.User.ID, params.ReservationIDs)

	if err != nil {
		return nil, &customError.InternalError{
			Code:    customError.DBError,
			Message: err.Error(),
		}
	}
	// At this point, all legit reservations are already deleted from DB
	// There should not be any error when reclaiming quotas
	// Keep retrying if there is error

	for k, v := range quotaToReclaims {
		for {
			err = ctx.RedisCache.IncEventQuota(k, v)
			if err == nil {
				break
			}
			logger.Errorf(err.Error())
		}

	}
	for {
		err = ctx.DB.ReclaimEventQuotas(quotaToReclaims)
		if err == nil {
			break
		}
		logger.Errorf(err.Error())
	}

	return &CancelReservationResults{DeletedTickets: deletedTickets}, nil
}
