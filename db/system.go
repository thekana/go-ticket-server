package db

// Idea: Objects are not deleted from memory. Use flags

// TODO: Complete all methods and start simple testing

type Event struct {
	id         int
	owner      int
	name       string
	quota      int
	tickets    []Reservation
	soldAmount int  // useful if updated quota less than soldAmount
	deleted    bool // set this to true after deleted
}

func (receiver *Event) isSoldOut() bool {
	if receiver.quota-receiver.soldAmount <= 0 {
		return true
	} else {
		return false
	}
}

// Admin or owner may call this function to delete event and void all tickets
func (receiver *Event) delete() {
	receiver.deleted = true
	for _, ticket := range receiver.tickets {
		ticket.voided = true
	}
}

type Reservation struct {
	id         int
	reservedBy int
	ownedBy    int
	eventId    int
	tickets    int
	voided     bool
}

type UserData struct {
	username     string
	userId       int
	reservations []Reservation // what this user reserves
	events       []Event       // what this user owns
}

type System struct {
	userMap   map[int]UserData
	eventList []Event       // allows admin/cust to view/delete all events
	eventMap  map[int]Event // allows for quick access
}
