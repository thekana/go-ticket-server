package db

// Idea: Objects are not deleted from memory. Use flags

import (
	"github.com/EagleChen/mapmutex"
	"github.com/google/uuid"
	"net/http"
	"sync"
	customError "ticket-reservation/custom_error"
)

var (
	EventNotFoundError = &customError.UserError{
		Code:           customError.UnknownError,
		Message:        "Event not found",
		HTTPStatusCode: http.StatusBadRequest,
	}
	QuotaError = &customError.UserError{
		Code:           0,
		Message:        "Not Enough RemainingQuota",
		HTTPStatusCode: http.StatusBadRequest,
	}
	SoldOutError = &customError.UserError{
		Code:           0,
		Message:        "Sold Out",
		HTTPStatusCode: http.StatusBadRequest,
	}
	UserNotFoundError = &customError.UserError{
		Code:           0,
		Message:        "User not found",
		HTTPStatusCode: http.StatusBadRequest,
	}
	ReservationNotFoundError = &customError.UserError{
		Code:           0,
		Message:        "Reservation not found",
		HTTPStatusCode: http.StatusBadRequest,
	}
)

type RWMap struct {
	sync.RWMutex
	m map[string]interface{}
}

func (r *RWMap) Get(key string) (interface{}, bool) {
	r.RLock()
	defer r.RUnlock()
	item, found := r.m[key]
	return item, found
}

func (r *RWMap) Set(key string, item interface{}) {
	r.Lock()
	defer r.Unlock()
	r.m[key] = item
}

func (r *RWMap) GetMap() map[string]interface{} {
	return r.m
}

type Event struct {
	ID          string
	OrganizerID int
	Name        string
	Quota       int
	Tickets     []*Reservation
	SoldAmount  int  // useful if updated quota less than soldAmount
	Deleted     bool // set this to true after deleted
}

func (e *Event) IsSoldOut() bool {
	if e.Quota-e.SoldAmount <= 0 {
		return true
	} else {
		return false
	}
}

func (e *Event) Delete() {
	e.Deleted = true
	for _, ticket := range e.Tickets {
		// This operation is thread-safe because we the only operation we do on ticket
		// after it is created is voiding it
		ticket.Voided = true
	}
}

func (e *Event) AddReservation(ticket *Reservation) {
	e.Tickets = append(e.Tickets, ticket)
}

type Reservation struct {
	ID         string
	ReservedBy int
	OwnedBy    int
	EventID    string
	Amount     int
	Voided     bool
}

type UserData struct {
	Username     string
	UserID       int
	Reservations []*Reservation // what this user reserves
	Events       []*Event       // what this user owns
	rLock        sync.Mutex
	eLock        sync.Mutex
}

func (d *UserData) AddReservation(res *Reservation) {
	d.rLock.Lock()
	d.Reservations = append(d.Reservations, res)
	d.rLock.Unlock()
}

func (d *UserData) AddEvent(res *Event) {
	d.eLock.Lock()
	d.Events = append(d.Events, res)
	d.eLock.Unlock()
}

type System struct {
	userMap       RWMap
	eventList     []*Event
	eventMap      RWMap
	resourceLock  *mapmutex.Mutex
	eventListLock sync.RWMutex
}

func NewUserData(username string, id int) *UserData {
	return &UserData{
		Username:     username,
		UserID:       id,
		Reservations: nil,
		Events:       nil,
	}
}

func NewEvent(ownerId int, eventName string, quota int) *Event {
	return &Event{
		ID:          uuid.New().String(),
		OrganizerID: ownerId,
		Name:        eventName,
		Quota:       quota,
		Tickets:     nil,
		SoldAmount:  0,
		Deleted:     false,
	}
}

func NewSystem() *System {
	return &System{
		userMap: RWMap{
			m: make(map[string]interface{}),
		},
		eventList: nil,
		eventMap: RWMap{
			m: make(map[string]interface{}),
		},
		resourceLock: mapmutex.NewMapMutex(),
	}
}

func (r *System) AddUserToSystem(user *UserData) {
	r.userMap.Set(string(user.UserID), user)
}

func (r *System) AddEventToSystem(event *Event) {
	r.eventListLock.Lock()
	r.eventList = append(r.eventList, event)
	r.eventListLock.Unlock()
	r.eventMap.Set(event.ID, event)
	user, _ := r.userMap.Get(string(event.OrganizerID))
	user.(*UserData).AddEvent(event)
}

func (r *System) DeleteEvent(eventID string) {
	e, _ := r.eventMap.Get(eventID)
	event := e.(*Event)
	// Acquire event resource lock
	r.resourceLock.TryLock(event.ID)
	event.Delete()
	r.resourceLock.Unlock(event.ID)
}

func (r *System) GetEvent(eventID string) (*Event, bool) {
	e, exist := r.eventMap.Get(eventID)
	return e.(*Event), exist
}

func (r *System) GetEventName(eventID string) string {
	e, exist := r.eventMap.Get(eventID)
	if exist {
		return e.(*Event).Name
	}
	return ""
}

func (r *System) GetAllEventsInSystem() []*Event {
	// Allow dirty read
	return r.eventList
}

func (r *System) GetEventsByUserID(userID int) []*Event {
	// Allow dirty read
	user, _ := r.userMap.Get(string(userID))
	return user.(*UserData).Events
}

func (r *System) GetReservationsByUserID(userID int) []*Reservation {
	// Allow dirty read
	user, _ := r.userMap.Get(string(userID))
	return user.(*UserData).Reservations
}

func (r *System) UserMakeReservation(userID int, eventID string, amount int) (*Reservation, error) {
	event, found := r.GetEvent(eventID)
	if !found || event.Deleted {
		return nil, EventNotFoundError
	}
	r.resourceLock.TryLock(eventID)
	if event.IsSoldOut() {
		return nil, SoldOutError
	}
	if event.Quota-event.SoldAmount-amount < 0 {
		return nil, QuotaError
	}
	ticket := &Reservation{
		ID:         uuid.New().String(),
		ReservedBy: userID,
		OwnedBy:    event.OrganizerID,
		EventID:    event.ID,
		Amount:     amount,
		Voided:     false,
	}
	// Commit ticket to event
	event.SoldAmount += ticket.Amount
	event.AddReservation(ticket)
	r.resourceLock.Unlock(eventID)

	u, _ := r.userMap.Get(string(userID))
	user := u.(*UserData)

	r.resourceLock.TryLock(user.Username)
	user.AddReservation(ticket)
	r.resourceLock.Unlock(user.Username)

	return ticket, nil
}

func (r *System) UserCancelReservation(userID int, reservationID string) error {
	// Lock associated Event here
	u, found := r.userMap.Get(string(userID))
	thisUser := u.(*UserData)
	if !found {
		return UserNotFoundError
	}
	// Now loop through all reservations
	found = false
	var thisReservation *Reservation
	for _, reservation := range thisUser.Reservations {
		if reservation.ID == reservationID && !reservation.Voided {
			found = true
			thisReservation = reservation
			break
		}
	}
	if !found {
		return ReservationNotFoundError
	}
	// Event reclaims quota
	// Acquire event lock
	e, _ := r.eventMap.Get(thisReservation.EventID)
	event := e.(*Event)
	r.resourceLock.TryLock(thisReservation.EventID)
	defer r.resourceLock.Unlock(thisReservation.EventID)
	if !event.Deleted {
		event.SoldAmount -= thisReservation.Amount
	}
	thisReservation.Voided = true
	return nil
}
