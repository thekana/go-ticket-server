package main

import (
	"github.com/davecgh/go-spew/spew"
	"ticket-reservation/db"
)

func main() {
	sys := db.NewSystem()
	a := db.NewUserData("king", 1)
	b := db.NewUserData("king2", 2)
	sys.AddUserToSystem(a)
	sys.AddUserToSystem(b)
	spew.Dump(sys)
	a.Username = "hello"
	spew.Dump(a)
	spew.Dump(sys)
	a = nil
	spew.Dump(sys)
}
