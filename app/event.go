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

// GetEventDetail will check Redis first then DB
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
	// TODO: Get data from cache first
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
			return nil, &customError.InternalError{
				Code:    customError.ConcurrencyIssue,
				Message: "CONCURRENCY ERROR",
			}
			logger.Errorf(err.Error())
		}
		return nil, &customError.InternalError{
			Code:    customError.UnknownError,
			Message: err.Error(),
		}
	}
	// TODO: Change key in redis
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
			return nil, &customError.InternalError{
				Code:    customError.ConcurrencyIssue,
				Message: "CONCURRENCY ERROR",
			}
			logger.Errorf(err.Error())
		}
		return nil, &customError.InternalError{
			Code:    customError.UnknownError,
			Message: err.Error(),
		}
	}
	//TODO: Delete key from redis
	return &DeleteEventResult{Message: result}, nil
}
