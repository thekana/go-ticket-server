package model

type Role int

const (
	Admin Role = iota + 1
	Organizer
	Customer
)
