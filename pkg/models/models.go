package models

import (
	"time"
)

// User is the user model
type User struct {
	ID          int
	FirstName   string
	LastName    string
	Email       string
	Password    string
	AccessLevel int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Reservation is the reservation model
type Reservation struct {
	ID        int
	FirstName string
	LastName  string
	Email     string
	Phone     string
	StartDate time.Time
	EndDate   time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
	Room      *Rooms
}

type Rooms struct {
	ID         int       `json:"ID"`
	RoomNumber string    `json:"roomnumber"`
	RoomType   *RoomType `json:"roomtype"`
}

type RoomType struct {
	ID   int    `json:"type-id"`
	Name string `json:"roomname"`
}

type MailData struct {
	From    string
	To      string
	Subject string
	Data    string
}
