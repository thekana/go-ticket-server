package app

import (
	"fmt"
	"strings"
	customError "ticket-reservation/custom_error"
	"ticket-reservation/db/model"
)

type CreateEventParams struct {
	AuthToken string `json:"authToken" validate:"required"`
	Name      string `json:"eventName" validate:"required"`
	Quota     int    `json:"quota" validate:"required"`
}

type CreateEventResult struct {
	EventID string `json:"eventID"`
}

type ViewEventParams struct {
	AuthToken string `json:"authToken" validate:"required"`
	EventID   string `json:"eventID" validate:"required"`
}

type ViewEventResult struct {
	Event interface{} `json:"event"`
}

type ViewAllEventParams struct {
	AuthToken string `json:"authToken" validate:"required"`
}
type ViewAllEventResult struct {
	Events interface{} `json:"events"`
}

type EditEventParams struct {
	AuthToken    string `json:"authToken" validate:"required"`
	EventID      string `json:"eventID" validate:"required"`
	NewEventName string `json:"newEventName" validate:"required"`
	NewQuota     int    `json:"newQuota" validate:"required"`
}

type EditEventResult struct {
	EditedEvent interface{} `json:"editedEvent"`
}

func (ctx *Context) CreateEvent(params CreateEventParams) (*CreateEventResult, error) {
	logger := ctx.getLogger()

	if err := validateInput(params); err != nil {
		logger.Errorf("validateInput error : %s", err)
		return nil, err
	}

	_, claims, err := ctx.verifyToken(params.AuthToken)
	if err != nil {
		return nil, err
	}
	roles := fmt.Sprint((*claims)["role"])
	if !strings.Contains(roles, "organizer") {
		return nil, &customError.AuthorizationError{
			Code:    customError.Unauthorized,
			Message: "for organizers only",
		}
	}
	eventID, err := ctx.DB.CreateEvent(int((*claims)["uid"].(float64)), params.Name, params.Quota)
	if err != nil {
		return nil, err
	}
	return &CreateEventResult{EventID: eventID}, nil
}

func (ctx *Context) GetEventDetail(params ViewEventParams) (*ViewEventResult, error) {
	logger := ctx.getLogger()

	if err := validateInput(params); err != nil {
		logger.Errorf("validateInput error : %s", err)
		return nil, err
	}

	tokenValid, _, err := ctx.verifyToken(params.AuthToken)
	if err != nil {
		return nil, err
	}
	if !tokenValid {
		return nil, &customError.AuthorizationError{
			Code:    customError.InvalidAuthToken,
			Message: "invalid token",
		}
	}
	eventDetail, err := ctx.DB.ViewEventDetail(params.EventID)
	if err != nil {
		return nil, err
	}
	return &ViewEventResult{Event: eventDetail}, nil
}

func (ctx *Context) GetAllEventDetails(params ViewAllEventParams) (*ViewAllEventResult, error) {
	logger := ctx.getLogger()

	if err := validateInput(params); err != nil {
		logger.Errorf("validateInput error : %s", err)
		return nil, err
	}

	tokenValid, claims, err := ctx.verifyToken(params.AuthToken)
	if err != nil {
		return nil, err
	}
	if !tokenValid {
		return nil, &customError.AuthorizationError{
			Code:    customError.InvalidAuthToken,
			Message: "invalid token",
		}
	}
	roles := fmt.Sprint((*claims)["role"])
	userID := int((*claims)["uid"].(float64))
	var event []*model.EventDetail
	if strings.Contains(roles, "organizer") {
		event, err = ctx.DB.OrganizerViewAllEvents(userID)
	} else {
		event, err = ctx.DB.ViewAllEvents()
	}
	if err != nil {
		return nil, err
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
	// Now check if the event is owned by user
	record, err := ctx.DB.EditEvent(params.EventID, params.NewEventName, params.NewQuota, int(authRes.User.ID))
	if err != nil {
		return nil, err
	}
	return &EditEventResult{EditedEvent: record}, nil
}
