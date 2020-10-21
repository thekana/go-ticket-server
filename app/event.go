package app

import (
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

type DeleteEventParams struct {
	AuthToken string `json:"authToken" validate:"required"`
	EventID   string `json:"eventID" validate:"required"`
}

type DeleteEventResult struct {
	Message string `json:"message"`
}

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
	eventID, err := ctx.DB.CreateEvent(int(authRes.User.ID), params.Name, params.Quota)
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
	_, err := ctx.authorizeUser(params.AuthToken, []model.Role{model.Customer, model.Admin, model.Organizer})
	if err != nil {
		return nil, err
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

	authRes, err := ctx.authorizeUser(params.AuthToken, []model.Role{model.Customer, model.Admin, model.Organizer})
	if err != nil {
		return nil, err
	}
	var event []*model.EventDetail
	if authRes.IsOrganizer {
		event, err = ctx.DB.OrganizerViewAllEvents(int(authRes.User.ID))
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

	record, err := ctx.DB.EditEvent(params.EventID, params.NewEventName, params.NewQuota, int(authRes.User.ID))
	if err != nil {
		return nil, err
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

	result, err := ctx.DB.DeleteEvent(params.EventID, int(authRes.User.ID), authRes.IsAdmin)

	if err != nil {
		return nil, err
	}

	return &DeleteEventResult{Message: result}, nil
}
