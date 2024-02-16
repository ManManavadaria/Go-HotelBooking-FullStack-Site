package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Man-Crest/GO-Projects/01_bookings/pkg/config"
	"github.com/Man-Crest/GO-Projects/01_bookings/pkg/connection"
	"github.com/Man-Crest/GO-Projects/01_bookings/pkg/forms"
	"github.com/Man-Crest/GO-Projects/01_bookings/pkg/helpers"
	"github.com/Man-Crest/GO-Projects/01_bookings/pkg/models"
	"github.com/Man-Crest/GO-Projects/01_bookings/pkg/render"
)

// Repo the repository used by the handlers
var Repo *Repository

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
	DB  connection.DB
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig, db connection.DB) *Repository {
	return &Repository{
		App: a,
		DB:  db,
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(r, w, "home.page.templ", &models.TemplateData{})
}

func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(r, w, "about.page.templ", &models.TemplateData{})
}

func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(r, w, "contact.page.templ", &models.TemplateData{})
}

func (m *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(r, w, "generals.page.templ", &models.TemplateData{})
}

func (m *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(r, w, "majors.page.templ", &models.TemplateData{})
}

func (m *Repository) MakeReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	if err != nil {
		http.Error(w, "Error parsing form data", http.StatusInternalServerError)
		return
	}

	availability, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)

	if !ok {
		log.Println("can't get item from session")
		m.App.Session.Put(r.Context(), "error", "Can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	reservations := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Email:     r.Form.Get("email"),
		Phone:     r.Form.Get("phone"),
		StartDate: availability.StartDate,
		EndDate:   availability.EndDate,
	}

	form := forms.New(r.PostForm)
	form.Has("first_name", r)
	form.Has("last_name", r)
	form.Has("email", r)
	form.Has("phone", r)

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservations
		render.RenderTemplate(r, w, "make-reservation.page.templ", &models.TemplateData{
			Form: form,
			Data: data,
		})
	}
	m.App.Session.Put(r.Context(), "reservation", reservations)
	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
}

func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {

	var emptyReservation models.Reservation
	data := make(map[string]interface{})

	data["reservation"] = emptyReservation
	render.RenderTemplate(r, w, "make-reservation.page.templ", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	room, ok1 := m.App.Session.Get(r.Context(), "room").(models.Rooms)
	if !ok && !ok1 {
		log.Println("can't get reservation from session in ReservationSummery")
		m.App.Session.Put(r.Context(), "error", "Can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	// m.App.Session.Remove(r.Context(), "reservation")
	data := make(map[string]interface{})
	data["reservation"] = reservation
	data["room"] = room

	render.RenderTemplate(r, w, "reservation-summary.page.templ", &models.TemplateData{
		Data: data,
	})
}

func (m *Repository) PostReservationSummary(w http.ResponseWriter, r *http.Request) {
	reservations, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	room, ok1 := m.App.Session.Get(r.Context(), "room").(models.Rooms)
	if !ok && !ok1 {
		log.Println("can't get reservation from session in PostReservationSummery")
		m.App.Session.Put(r.Context(), "error", "Can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	query := `
	INSERT INTO reservations (first_name, last_name, email, phone, start_date, end_date, room_id)
	SELECT $1, $2, $3, $4, $5, $6, r.id
	FROM rooms r
	INNER JOIN roomtype rt ON r.roomtypeid = rt.id
	WHERE r.roomnumber = $7 AND rt.name = $8
	`
	_, err := m.DB.SQL.Exec(query, reservations.FirstName, reservations.LastName, reservations.Email, reservations.Phone, reservations.StartDate, reservations.EndDate, room.RoomNumber, room.RoomType.Name)

	if err != nil {
		// Handle the error here
		fmt.Println("error parsing query sql", err)
		return
	}
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	m.App.Session.Put(r.Context(), "flash", "Booking Confirmed")

	htmlMsg := fmt.Sprintf(`<h1>Congratulations , %s </h1><h2>your boking is confirmed from <b>%s</b> to <b>%s</b></h2><br><h2>Room no: %s</h2><br><h2>Room Type: %s</h2>`, reservations.FirstName, reservations.StartDate.Format("2024-01-01"), reservations.EndDate.Format("2024-01-01"), room.RoomNumber, room.RoomType.Name)

	mailData := models.MailData{
		From:    "man.m@crestinfosystems.com",
		To:      reservations.Email,
		Subject: "Booking Confermed",
		Data:    htmlMsg,
	}
	m.App.MailChan <- mailData

	m.App.Session.Remove(r.Context(), "reservation")
	m.App.Session.Remove(r.Context(), "room")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (m *Repository) SearchAvailability(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(r, w, "search-availability.page.templ", &models.TemplateData{})
}

func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form data", http.StatusInternalServerError)
		return
	}
	s := r.Form.Get("start")
	e := r.Form.Get("end")
	room_type := r.Form.Get("roomType")

	// roomType, err := strconv.Atoi(room_type)

	if err != nil {
		// Handle parsing error
		http.Error(w, "Error parsing room type", http.StatusInternalServerError)
		return
	}

	startDate, err := time.Parse("2006-01-02", s)

	if err != nil {
		// Handle parsing error
		http.Error(w, "Error parsing start date", http.StatusInternalServerError)
		return
	}
	endDate, err := time.Parse("2006-01-02", e)
	if err != nil {
		// Handle parsing error
		http.Error(w, "Error parsing end date", http.StatusInternalServerError)
		return
	}

	roomType := &models.RoomType{
		Name: room_type,
	}
	rooms, err := m.Availability(roomType, startDate, endDate)

	if err != nil {
		log.Fatal(err, ": error auccerd in availability handler")
	}

	reservations := &models.Reservation{
		StartDate: startDate,
		EndDate:   endDate,
	}

	m.App.Session.Put(r.Context(), "reservation", reservations)
	m.App.Session.Put(r.Context(), "rooms", rooms)
	http.Redirect(w, r, "/room-availability", http.StatusSeeOther)
}

func (m *Repository) RoomAvailability(w http.ResponseWriter, r *http.Request) {
	rooms, ok := m.App.Session.Get(r.Context(), "rooms").([]models.Rooms)
	if !ok {
		log.Println("can't get item from session")
		m.App.Session.Put(r.Context(), "error", "Can't get rooms from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// m.App.Session.Remove(r.Context(), "rooms")

	data := make(map[string]interface{})
	data["rooms"] = rooms

	render.RenderTemplate(r, w, "room-availability.page.templ", &models.TemplateData{
		Data: data,
	})
}

func (m *Repository) PostRoomAvailability(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()

	if err != nil {
		log.Fatal(err, "not getting data from room availability form")
	}

	ID := r.Form.Get("roomid")
	id, _ := strconv.Atoi(ID)
	room := models.Rooms{
		ID:         id,
		RoomNumber: r.Form.Get("roomnumber"),
		RoomType: &models.RoomType{
			Name: r.Form.Get("roomname"),
		},
	}
	m.App.Session.Put(r.Context(), "room", room)
	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}

func (m *Repository) LoginHandler(w http.ResponseWriter, r *http.Request) {

	render.RenderTemplate(r, w, "authentication.page.templ", &models.TemplateData{
		Form: forms.New(nil),
	})
}

func (m *Repository) PostLoginHandler(w http.ResponseWriter, r *http.Request) {

	err := m.App.Session.RenewToken(r.Context())
	if err != nil {
		log.Fatal(err)
	}

	err = r.ParseForm()
	if err != nil {
		log.Fatal("error parsing form in sign in handler", err)
	}
	email := r.Form.Get("email")
	password := r.Form.Get("password")

	form := forms.New(r.PostForm)
	form.Has("email", r)
	form.Has("password", r)

	if !form.Valid() {
		render.RenderTemplate(r, w, "authentication.page.templ", &models.TemplateData{
			Form: form,
		})
	} else {
		var storedpassword string
		var user_id int
		err = m.DB.SQL.QueryRow("SELECT password,ID FROM users WHERE email = $1", email).Scan(&storedpassword, &user_id)
		if err != nil {
			log.Fatal("err in sign in query :", err)
			if err == sql.ErrNoRows {
				http.Redirect(w, r, "/login-user", http.StatusSeeOther)
			}
		}
		if password == storedpassword {
			m.App.Session.Put(r.Context(), "flash", "Successfully logged in")
			m.App.Session.Put(r.Context(), "IsLoggedIn", true)
			m.App.Session.Put(r.Context(), "user_id", user_id)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			fmt.Println(storedpassword)
			fmt.Println(user_id)
		} else {
			m.App.Session.Put(r.Context(), "error", "Please enter valid credentials")
			// m.App.Session.Put(r.Context(), "IsLoggedIn", "false")
			http.Redirect(w, r, "/login-user", http.StatusSeeOther)
		}
	}
}

func (m *Repository) SignupHandler(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(r, w, "signup.page.templ", &models.TemplateData{
		Form: forms.New(nil),
	})
}

func (m *Repository) PostSignupHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}

	user := &models.User{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Password:  r.Form.Get("password"),
		Email:     r.Form.Get("email"),
	}

	form := forms.New(r.PostForm)
	form.Has("first_name", r)
	form.Has("last_name", r)
	form.Has("email", r)
	form.Has("password", r)

	if !form.Valid() {
		render.RenderTemplate(r, w, "signup.page.templ", &models.TemplateData{
			Form: form,
		})
	}

	query := `INSERT INTO users (firstname,lastname,email,password)
	VALUES ($1,$2,$3,$4)`

	_, err = m.DB.SQL.Exec(query, user.FirstName, user.LastName, user.Email, user.Password)

	if err != nil {
		log.Println(err)
	}

	eq := struct {
		IsLoggedIn bool
	}{
		IsLoggedIn: true,
	}
	m.App.Session.Put(r.Context(), "eq", eq)
	m.App.Session.Put(r.Context(), "flash", "Successfully Sign Up")
	http.Redirect(w, r, "/login-user", http.StatusSeeOther)
}

func (m *Repository) Logout(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.Destroy(r.Context())
	_ = m.App.Session.RenewToken(r.Context())

	http.Redirect(w, r, "/login-user", http.StatusSeeOther)
}

// func (m *Repository) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
// 	_, ok := m.App.Session.Get(r.Context(), "user_id").(models.User)

// 	if ok != true {
// 		log.Println("user is not logged in")
// 		http.Redirect(w, r, "/login-user", http.StatusSeeOther)
// 	}

// 	query := `UPDATE users
// 	SET
// 		FirstName = $1,
// 		LastName = $2,
// 		Password = $3,
// 	WHERE
// 		email = $4;
// 	`
// 	err := r.ParseForm()

// 	if err != nil {
// 		log.Println(err)
// 	}

// 	_ = m.DB.SQL.QueryRowContext(r.Context(), query, r.Form.Get("firstname"), r.Form.Get("lasttname"), r.Form.Get("email"), r.Form.Get("password"))

// }
