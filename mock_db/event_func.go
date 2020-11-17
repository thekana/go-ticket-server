package mock_db

import (
	"net/http"
	customError "ticket-reservation/custom_error"
	"ticket-reservation/db/model"
)

// MockGetEventDetailForEvent1Only only has eventID 1 in the database. It will return event struct for event 1 only
func MockGetEventDetailForEvent1Only(eventId int) (*model.EventDetail, error) {
	if eventId != 1 {
		return nil, &customError.UserError{
			Code:           customError.EventNotFound,
			Message:        "Some error",
			HTTPStatusCode: http.StatusNotFound,
		}
	}
	return &model.EventDetail{
		EventID:        1,
		OrganizerID:    1,
		EventName:      "mock data",
		Quota:          1,
		RemainingQuota: 1,
	}, nil
}
