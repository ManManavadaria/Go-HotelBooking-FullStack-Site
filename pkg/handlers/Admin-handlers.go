package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Man-Crest/GO-Projects/01_bookings/pkg/forms"
	"github.com/Man-Crest/GO-Projects/01_bookings/pkg/models"
	"github.com/Man-Crest/GO-Projects/01_bookings/pkg/render"
)

func (m *Repository) AdminDashboard(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(r, w, "admin-dashboard.page.templ", &models.TemplateData{})
}

// type idate struct {
// 	startdate string
// 	enddate   string
// }

// var formatedDate []idate

func (m *Repository) AdminAllReservations(w http.ResponseWriter, r *http.Request) {

	ctx, cancle := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancle()

	var reservations []models.Reservation

	query := `SELECT
    r.id,
    r.first_name,
    r.last_name,
    r.email,
    r.phone,
    r.start_date,
    r.end_date,
    r.created_at,
    r.updated_at,
    ro.roomnumber,
    rt.name AS roomtypename
FROM
    reservations r
JOIN
    rooms ro ON r.room_id = ro.id
JOIN
    roomtype rt ON ro.roomtypeid = rt.id;
`
	rows, err := m.DB.SQL.QueryContext(ctx, query)
	defer rows.Close()

	if err != nil {
		log.Fatal(err, "error retriving all reservations from DB")
	}

	for rows.Next() {
		i := models.Reservation{
			Room: &models.Rooms{
				RoomType: &models.RoomType{},
			},
		}
		err := rows.Scan(
			&i.ID,
			&i.FirstName,
			&i.LastName,
			&i.Email,
			&i.Phone,
			&i.StartDate,
			&i.EndDate,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Room.RoomNumber,
			&i.Room.RoomType.Name,
		)

		if err != nil {
			log.Fatal(err, "error scanning rows of all rservations")
		}
		reservations = append(reservations, i)
	}
	if err := rows.Err(); err != nil {
		log.Println("Error iterating through rows:", err)
	}

	data := make(map[string]interface{})
	data["reservations"] = reservations

	render.RenderTemplate(r, w, "admin-all-reservations.page.templ", &models.TemplateData{
		Data: data,
	})
}

func (m *Repository) AdminNewReservations(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(r, w, "admin-new-reservations.page.templ", &models.TemplateData{})
}

func (m *Repository) AdminReservationsDetail(w http.ResponseWriter, r *http.Request) {

	ctx, cancle := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancle()
	query := `SELECT
    r.id,
    r.first_name,
    r.last_name,
    r.email,
    r.phone,
    r.start_date,
    r.end_date,
    r.created_at,
    r.updated_at,
    ro.roomnumber,
    rt.name AS roomtypename
FROM
    reservations r
JOIN
    rooms ro ON r.room_id = ro.id
JOIN
    roomtype rt ON ro.roomtypeid = rt.id 
WHERE r.id = $1	
;
`
	uri := strings.Split(r.RequestURI, "/")

	id, err := strconv.Atoi(uri[4])
	if err != nil {
		log.Fatal(err)
	}

	rows := m.DB.SQL.QueryRowContext(ctx, query, id)

	// if err != nil {
	// 	log.Fatal(err, "error retriving all reservations from DB")
	// }

	var i = &models.Reservation{
		Room: &models.Rooms{
			RoomType: &models.RoomType{},
		},
	}

	err = rows.Scan(
		&i.ID,
		&i.FirstName,
		&i.LastName,
		&i.Email,
		&i.Phone,
		&i.StartDate,
		&i.EndDate,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Room.RoomNumber,
		&i.Room.RoomType.Name,
	)
	if err != nil {
		log.Fatal(err, "error scanning rows of all rservations")
	}

	if err := rows.Err(); err != nil {
		log.Println("Error iterating through rows:", err)
	}

	data := make(map[string]interface{})
	data["reservation-detail"] = i

	render.RenderTemplate(r, w, "reservations-detail.page.templ", &models.TemplateData{
		Data: data,
	})
}

func (m *Repository) AdminUpdateReservation(w http.ResponseWriter, r *http.Request) {

	ctx, cancle := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancle()
	uri := strings.Split(r.RequestURI, "/")

	id, err := strconv.Atoi(uri[4])
	if err != nil {
		log.Fatal(err)
	}

	updateQuery := `update reservations set first_name = $1, last_name = $2, email = $3, phone = $4, updated_at = $5
	where id = $6`

	err = r.ParseForm()
	if err != nil {
		log.Fatal("err parsing form", err)
	}

	updateRes := m.DB.SQL.QueryRowContext(ctx, updateQuery,
		r.Form.Get("firstname"),
		r.Form.Get("lastname"),
		r.Form.Get("email"),
		r.Form.Get("phone"),
		time.Now(),
		id)
	if err != nil {
		log.Fatal("error executing query of update", err)
	}
	fmt.Println(updateRes)

	m.App.Session.Put(r.Context(), "flash", "Updated successfully")
	http.Redirect(w, r, "/admin/reservations-all", http.StatusSeeOther)
}

func (m *Repository) AdminCancleReservation(w http.ResponseWriter, r *http.Request) {
	uri := strings.Split(r.RequestURI, "/")

	id, err := strconv.Atoi(uri[4])
	if err != nil {
		log.Fatal(err)
	}

	_ = m.DB.SQL.QueryRow(`DELETE FROM reservations WHERE id = $1;`, id)

	http.Redirect(w, r, "/admin/reservations-all", http.StatusSeeOther)
}

// func (m *Repository) AdminCalendar(w http.ResponseWriter, r *http.Request) {
// 	Year, _, Day := time.Now().Date()

// 	currentTime := time.Now()
// 	currentMonth := int(currentTime.Month())

// 	data := make(map[string]interface{})

// 	data["year"] = Year
// 	data["month"] = currentMonth
// 	data["day"] = Day

// 	render.RenderTemplate(r, w, "admin-calendar.page.templ", &models.TemplateData{
// 		Data: data,
// 	})

// 	rows, err := m.DB.SQL.Query(`
// 	SELECT DISTINCT rooms.id, rooms.roomnumber, roomtype.name
// FROM rooms
// JOIN roomtype ON rooms.roomtypeid = roomtype.id
// LEFT JOIN reservations ON rooms.id = reservations.room_id AND (
//     $1 BETWEEN reservations.start_date AND reservations.end_date
//     OR $2 BETWEEN reservations.start_date AND reservations.end_date
//     OR ($1 <= reservations.start_date AND $2 >= reservations.end_date)
// )
// WHERE reservations.room_id IS NULL;`, StartDate, EndDate)

// 	if err != nil {
// 		log.Println("Error executing SQL query:", err)
// 	}

// 	defer rows.Close()

// 	var rooms []models.Rooms

// 	// Iterate through the rows to fetch available rooms
// 	for rows.Next() {
// 		var room models.Rooms
// 		var RoomTypeName string
// 		if err := rows.Scan(&room.ID, &room.RoomNumber, &RoomTypeName); err != nil {
// 			log.Println("Error scanning rows:", err)
// 		}
// 		room.RoomType = &models.RoomType{Name: RoomTypeName}
// 		if room.RoomType.Name == roomType.Name {
// 			rooms = append(rooms, room)
// 		}
// 	}
// 	if err := rows.Err(); err != nil {

// 	}
// }

func (m *Repository) PostAdminReservationsCalendar(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()

	if err != nil {
		log.Println("Error executing SQL query:", err)
	}
	desiredStartDate := r.Form.Get("start")
	desiredEndDate := r.Form.Get("end")

	rows, err := m.DB.SQL.Query(`
	SELECT DISTINCT rooms.id, rooms.roomnumber, roomtype.name
FROM rooms
JOIN roomtype ON rooms.roomtypeid = roomtype.id
LEFT JOIN reservations ON rooms.id = reservations.room_id 
WHERE (($1 NOT BETWEEN reservations.start_date AND reservations.end_date)
  AND ($2 NOT BETWEEN reservations.start_date AND reservations.end_date)
  AND NOT($1 < reservations.start_date AND $2 > reservations.end_date)
  )`, desiredStartDate, desiredEndDate)

	if err != nil {
		log.Println("Error executing SQL query:", err)
	}

	defer rows.Close()

	// var roomsAvailable []models.Rooms
	var allrooms []models.Rooms

	// Iterate through the rows to fetch available rooms
	for rows.Next() {
		var room models.Rooms
		var RoomTypeName string
		if err := rows.Scan(&room.ID, &room.RoomNumber, &RoomTypeName); err != nil {
			log.Println("Error scanning rows:", err)
		}
		room.RoomType = &models.RoomType{Name: RoomTypeName}
		allrooms = append(allrooms, room)
	}
	if err := rows.Err(); err != nil {
		log.Println("Error iterating through rows:", err)
	}

	data := make(map[string]interface{})
	data["allrooms"] = allrooms

	render.RenderTemplate(r, w, "admin-calendar-rooms.page.templ", &models.TemplateData{
		Data: data,
	})
}

func (m *Repository) AdminReservationsCalendar(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(r, w, "admin-calendar.page.templ", &models.TemplateData{
		Form: forms.New(nil),
	})
}
