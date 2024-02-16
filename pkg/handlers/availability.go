package handlers

import (
	"log"
	"time"

	"github.com/Man-Crest/GO-Projects/01_bookings/pkg/models"
)

func (m *Repository) Availability(roomType *models.RoomType, desiredStartDate time.Time, desiredEndDate time.Time) ([]models.Rooms, error) {

	rows, err := m.DB.SQL.Query(`
	SELECT DISTINCT rooms.id, rooms.roomnumber, roomtype.name
FROM rooms
JOIN roomtype ON rooms.roomtypeid = roomtype.id
LEFT JOIN reservations ON rooms.id = reservations.room_id AND (
    $1 BETWEEN reservations.start_date AND reservations.end_date
    OR $2 BETWEEN reservations.start_date AND reservations.end_date
    OR ($1 <= reservations.start_date AND $2 >= reservations.end_date)
)
WHERE reservations.room_id IS NULL;`, desiredStartDate, desiredEndDate)

	if err != nil {
		log.Println("Error executing SQL query:", err)
		return nil, err
	}

	defer rows.Close()

	var rooms []models.Rooms

	// Iterate through the rows to fetch available rooms
	for rows.Next() {
		var room models.Rooms
		var RoomTypeName string
		if err := rows.Scan(&room.ID, &room.RoomNumber, &RoomTypeName); err != nil {
			log.Println("Error scanning rows:", err)
			return nil, err
		}
		room.RoomType = &models.RoomType{Name: RoomTypeName}
		if room.RoomType.Name == roomType.Name {
			rooms = append(rooms, room)
		}
	}
	if err := rows.Err(); err != nil {
		log.Println("Error iterating through rows:", err)
		return nil, err
	}
	return rooms, nil
}
