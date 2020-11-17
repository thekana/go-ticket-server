package app

import (
	"github.com/jackc/pgerrcode"
	"net/http"
	customError "ticket-reservation/custom_error"
	"ticket-reservation/db/model"
)

type CreateEventParams struct {
	AuthToken string
	Name      string `json:"eventName" validate:"required"`
	Quota     int    `json:"quota" validate:"required"`
}

type CreateEventResult struct {
	EventID int `json:"eventID"`
}

type ViewEventParams struct {
	AuthToken string
	EventID   int `json:"eventID" validate:"required"`
}

type ViewEventResult struct {
	Event interface{} `json:"event"`
}

type ViewAllEventParams struct {
	AuthToken string
}
type ViewAllEventResult struct {
	Events interface{} `json:"events"`
}

type EditEventParams struct {
	AuthToken    string
	EventID      int    `json:"eventID" validate:"required"`
	NewEventName string `json:"newEventName" validate:"required"`
	NewQuota     int    `json:"newQuota" validate:"required"`
	OldQuota     int    `json:"originalQuota" validate:"required"` // Need this for redis logic
}

type EditEventResult struct {
	EditedEvent interface{} `json:"editedEvent"`
}

type DeleteEventParams struct {
	AuthToken string
	EventID   int `json:"eventID" validate:"required"`
}

type DeleteEventResult struct {
	Message string `json:"message"`
}

const RETRY = 20

func (ctx *Context) CreateEvent(params CreateEventParams) (*CreateEventResult, error) {
	logger := ctx.getLogger()

	if err := validateInput(params); err != nil {
		logger.Errorf("validateInput error : %s", err)
		return nil, err
	}

	authRes, err := ctx.authorizeUser(params.AuthToken, []model.Role{model.Organizer})
	if err != nil {
		return nil, err
	}
	eventID, err := ctx.DB.CreateEvent(authRes.User.ID, params.Name, params.Quota)
	if err != nil {
		return nil, &customError.UserError{
			Code:           customError.BadInput,
			Message:        err.Error(),
			HTTPStatusCode: http.StatusBadRequest,
		}
	}
	return &CreateEventResult{EventID: eventID}, nil
}

// GetEventDetail will check always redis
func (ctx *Context) GetEventDetail(params ViewEventParams) (*ViewEventResult, error) {
	logger := ctx.getLogger()

	if err := validateInput(params); err != nil {
		logger.Errorf("validateInput error : %s", err)
		return nil, err
	}
	_, err := ctx.authorizeUser(params.AuthToken, []model.Role{model.Customer, model.Admin, model.Organizer})
	if err != nil {
		return nil, err
	}
	// Since redis only store quota we must query DB every time for other parameters
	eventDetail, err := ctx.DB.ViewEventDetail(params.EventID)
	if err != nil {
		// So many possible error, but high likely going to be
		// event not found
		return nil, &customError.UserError{
			Code:           customError.EventNotFound,
			Message:        err.Error(),
			HTTPStatusCode: http.StatusNotFound,
		}
	}
	// Also get latest data from Redis and return
	// Ignore redis error here
	redisQuota, _ := ctx.RedisCache.GetEventQuota(params.EventID)
	if redisQuota == -1 {
		// not in redis so put it in
		ctx.RedisCache.SetNXEventQuota(params.EventID, eventDetail.RemainingQuota)
	} else {
		eventDetail.RemainingQuota = redisQuota
	}

	return &ViewEventResult{Event: eventDetail}, nil
}

// GetAllEventDetails will always query from DB
func (ctx *Context) GetAllEventDetails(params ViewAllEventParams) (*ViewAllEventResult, error) {
	logger := ctx.getLogger()

	if err := validateInput(params); err != nil {
		logger.Errorf("validateInput error : %s", err)
		return nil, err
	}

	authRes, err := ctx.authorizeUser(params.AuthToken, []model.Role{model.Customer, model.Admin, model.Organizer})
	if err != nil {
		return nil, err
	}
	var event []*model.EventDetail
	event, err = ctx.DB.ViewAllEvents(authRes.IsOrganizer, authRes.User.ID)
	if err != nil {
		return nil, &customError.InternalError{
			Code:    customError.DBError,
			Message: err.Error(),
		}
	}
	return &ViewAllEventResult{Events: event}, nil
}

func (ctx *Context) EditEventDetail(params EditEventParams) (*EditEventResult, error) {
	logger := ctx.getLogger()

	if err := validateInput(params); err != nil {
		logger.Errorf("validateInput error : %s", err)
		return nil, err
	}

	authRes, err := ctx.authorizeUser(params.AuthToken, []model.Role{model.Organizer})
	if err != nil {
		return nil, err
	}

	currentQuota, _ := ctx.RedisCache.GetEventQuota(params.EventID)
	if currentQuota >= 0 {
		// Quota found, perform check here
		sold := params.OldQuota - currentQuota
		if params.NewQuota < sold {
			return nil, &customError.UserError{
				Code:           customError.BadInput,
				Message:        "New quota must be more than sold tickets",
				HTTPStatusCode: http.StatusBadRequest,
			}
		}
		// Everything checks out
		// Update redis so users can continue to send requests
		diff := params.NewQuota - params.OldQuota
		if diff >= 0 {
			err = ctx.RedisCache.IncEventQuota(params.EventID, diff)
		} else {
			err = ctx.RedisCache.DecEventQuota(params.EventID, -diff)
		}
		if err != nil {
			return nil, err
		}
	}

	var record *model.EventDetail
	for i := 0; i < RETRY; i++ {
		record, err = ctx.DB.EditEvent(params.EventID, params.NewEventName, params.NewQuota, authRes.User.ID)
		if err != nil {
			if checkPostgresErrorCode(err, pgerrcode.SerializationFailure) {
				continue
			}
		}
		break
	}

	if err != nil {
		if checkPostgresErrorCode(err, pgerrcode.SerializationFailure) {
			logger.Errorf(err.Error())
			return nil, &customError.InternalError{
				Code:    customError.ConcurrencyIssue,
				Message: "CONCURRENCY ERROR",
			}
		}
		return nil, &customError.InternalError{
			Code:    customError.UnknownError,
			Message: err.Error(),
		}
	}

	return &EditEventResult{EditedEvent: record}, nil
}

func (ctx *Context) DeleteEvent(params DeleteEventParams) (*DeleteEventResult, error) {
	logger := ctx.getLogger()

	if err := validateInput(params); err != nil {
		logger.Errorf("validateInput error : %s", err)
		return nil, err
	}

	authRes, err := ctx.authorizeUser(params.AuthToken, []model.Role{model.Admin, model.Organizer})
	if err != nil {
		return nil, err
	}
	_ = ctx.RedisCache.DelEventQuota(params.EventID)
	var result string
	for i := 0; i < RETRY; i++ {
		result, err = ctx.DB.DeleteEvent(params.EventID, authRes.User.ID, authRes.IsAdmin)
		if err != nil {
			if checkPostgresErrorCode(err, pgerrcode.SerializationFailure) {
				continue
			}
		}
		break
	}
	if err != nil {
		if checkPostgresErrorCode(err, pgerrcode.SerializationFailure) {
			logger.Errorf(err.Error())
			return nil, &customError.InternalError{
				Code:    customError.ConcurrencyIssue,
				Message: "CONCURRENCY ERROR",
			}
		}
		return nil, &customError.InternalError{
			Code:    customError.UnknownError,
			Message: err.Error(),
		}
	}
	return &DeleteEventResult{Message: result}, nil
}
