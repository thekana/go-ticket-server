package db

// Idea: Objects are not deleted from memory. Use flags

import (
	"fmt"
	"github.com/EagleChen/mapmutex"
	"github.com/google/uuid"
)

type Event struct {
	ID          string
	OrganizerID int
	Name        string
	Quota       int
	Tickets     []*Reservation
	SoldAmount  int  // useful if updated quota less than soldAmount
	Deleted     bool // set this to true after deleted
}

type Reservation struct {
	ID         string
	ReservedBy int
	OwnedBy    int
	EventID    string
	Tickets    int
	Voided     bool
}

type UserData struct {
	Username     string
	UserID       int
	Reservations []*Reservation // what this user reserves
	Events       []*Event       // what this user owns
}

type System struct {
	userMap      map[int]*UserData
	eventList    []*Event
	eventMap     map[string]*Event
	resourceLock *mapmutex.Mutex
}

func NewUserData(username string, id int) *UserData {
	return &UserData{
		Username:     username,
		UserID:       id,
		Reservations: nil,
		Events:       nil,
	}
}

func NewReservation(eventId string, userId int, ownerId int, amountReserved int) *Reservation {
	return &Reservation{
		ID:         uuid.New().String(),
		ReservedBy: userId,
		OwnedBy:    ownerId,
		EventID:    eventId,
		Tickets:    amountReserved,
		Voided:     false,
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
		userMap:      make(map[int]*UserData),
		eventList:    nil,
		eventMap:     make(map[string]*Event),
		resourceLock: mapmutex.NewMapMutex(),
	}
}

func (receiver *Event) IsSoldOut() bool {
	if receiver.Quota-receiver.SoldAmount <= 0 {
		return true
	} else {
		return false
	}
}

// Admin or owner may call this function to delete event and void all Tickets
func (receiver *Event) Delete() {
	receiver.Deleted = true
	for _, ticket := range receiver.Tickets {
		ticket.Voided = true
	}
}

func (receiver *UserData) AddReservation(res *Reservation) {
	receiver.Reservations = append(receiver.Reservations, res)
}

// Every mutation & read should be called from System functions so it can handle locks

func (receiver *UserData) AddEvent(res *Event) {
	receiver.Events = append(receiver.Events, res)
	fmt.Print(receiver.Events)
}

// TODO: Add locks
func (receiver *System) AddUserToSystem(user *UserData) {
	receiver.userMap[user.UserID] = user
}

// TODO: Add locks
func (receiver *System) AddEventToSystem(event *Event) {
	receiver.eventList = append(receiver.eventList, event)
	receiver.eventMap[event.ID] = event
	receiver.userMap[event.OrganizerID].AddEvent(event)
}

// TODO: Add locks
func (receiver *System) DeleteEvent(eventID string) {
	e, _ := receiver.eventMap[eventID]
	e.Delete()
}

// TODO: Add locks
func (receiver *System) GetEvent(eventID string) (*Event, bool) {
	e, exist := receiver.eventMap[eventID]
	return e, exist
}

func (receiver *System) GetAllEvents() []*Event {
	return receiver.eventList
}

// TODO: Add helper functions

func (receiver *System) GetEventsOwnedByUser(uid int) []*Event {
	return receiver.userMap[uid].Events
}
