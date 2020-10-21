package db

// Idea: Objects are not deleted from memory. Use flags

import (
	"fmt"
	"github.com/EagleChen/mapmutex"
	"github.com/google/uuid"
)

type Event struct {
	Id          string
	OrganizerID int
	Name        string
	Quota       int
	Tickets     []*Reservation
	SoldAmount  int  // useful if updated quota less than soldAmount
	Deleted     bool // set this to true after deleted
}

func NewEvent(ownerId int, eventName string, quota int) *Event {
	return &Event{
		Id:          uuid.New().String(),
		OrganizerID: ownerId,
		Name:        eventName,
		Quota:       quota,
		Tickets:     nil,
		SoldAmount:  0,
		Deleted:     false,
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

type Reservation struct {
	Id         string
	ReservedBy int
	OwnedBy    int
	EventId    string
	Tickets    int
	Voided     bool
}

func NewReservation(eventId string, userId int, ownerId int, amountReserved int) *Reservation {
	return &Reservation{
		Id:         uuid.New().String(),
		ReservedBy: userId,
		OwnedBy:    ownerId,
		EventId:    eventId,
		Tickets:    amountReserved,
		Voided:     false,
	}
}

type UserData struct {
	Username     string
	UserId       int
	Reservations []*Reservation // what this user reserves
	Events       []*Event       // what this user owns
}

func NewUserData(username string, id int) *UserData {
	return &UserData{
		Username:     username,
		UserId:       id,
		Reservations: nil,
		Events:       nil,
	}
}

func (receiver *UserData) AddReservation(res *Reservation) {
	receiver.Reservations = append(receiver.Reservations, res)
}

func (receiver *UserData) AddEvent(res *Event) {
	receiver.Events = append(receiver.Events, res)
	fmt.Print(receiver.Events)
}

type System struct {
	userMap      map[int]*UserData
	eventList    []*Event
	eventMap     map[string]*Event
	resourceLock *mapmutex.Mutex
}

// Every mutation & read should be called from System functions so it can handle locks

// TODO: Add locks
func (receiver *System) AddUserToSystem(user *UserData) {
	receiver.userMap[user.UserId] = user
}

// TODO: Add locks
func (receiver *System) AddEventToSystem(event *Event) {
	receiver.eventList = append(receiver.eventList, event)
	receiver.eventMap[event.Id] = event
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

func (receiver *System) GetEventsOwnedByUser(uid int) []*Event {
	return receiver.userMap[uid].Events
}

// TODO: Add helper functions

func NewSystem() *System {
	return &System{
		userMap:      make(map[int]*UserData),
		eventList:    nil,
		eventMap:     make(map[string]*Event),
		resourceLock: mapmutex.NewMapMutex(),
	}
}
