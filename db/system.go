package db

// Idea: Objects are not deleted from memory. Use flags

import (
	"fmt"
	"github.com/EagleChen/mapmutex"
	"github.com/google/uuid"
)

// TODO: Complete all methods and start simple testing

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
	UserMap      map[int]*UserData
	EventList    []*Event          // allows admin/cust to view/delete all events
	EventMap     map[string]*Event // allows for quick access
	ResourceLock *mapmutex.Mutex
}

func (receiver *System) AddUserToSystem(user *UserData) {
	receiver.UserMap[user.UserId] = user
}

func (receiver *System) AddEventToSystem(event *Event) {
	receiver.EventList = append(receiver.EventList, event)
	receiver.EventMap[event.Id] = event
}

// TODO: Add helper functions

func NewSystem() *System {
	return &System{
		UserMap:      make(map[int]*UserData),
		EventList:    nil,
		EventMap:     make(map[string]*Event),
		ResourceLock: mapmutex.NewMapMutex(),
	}
}

//var system System = NewSystem()
