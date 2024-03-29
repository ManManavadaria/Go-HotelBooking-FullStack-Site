package main

import (
	"net/http"

	"github.com/Man-Crest/GO-Projects/01_bookings/pkg/config"
	"github.com/Man-Crest/GO-Projects/01_bookings/pkg/handlers"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func routes(app *config.AppConfig) http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	// mux.Use(NoSurf)
	mux.Use(SessionLoad)

	mux.Get("/", handlers.Repo.Home)
	mux.Get("/about", handlers.Repo.About)
	mux.Get("/contact", handlers.Repo.Contact)
	mux.Get("/make-reservation", handlers.Repo.Reservation)
	mux.Post("/make-reservation", handlers.Repo.MakeReservation)
	mux.Get("/reservation-summary", handlers.Repo.ReservationSummary)
	mux.Get("/confirm-booking", handlers.Repo.PostReservationSummary)
	mux.Get("/rooms", handlers.Repo.Rooms)
	mux.Get("/search-availability", handlers.Repo.SelectRooms)
	mux.Post("/search-availability", handlers.Repo.PostAvailability)
	mux.Get("/room-availability", handlers.Repo.RoomAvailability)
	mux.Post("/room-availability", handlers.Repo.PostRoomAvailability)
	mux.Get("/login-user", handlers.Repo.LoginHandler)
	mux.Post("/login-user", handlers.Repo.PostLoginHandler)
	mux.Post("/signup-user", handlers.Repo.PostSignupHandler)
	mux.Get("/user/logout", handlers.Repo.Logout)

	fs := http.FileServer(http.Dir("./static/"))

	mux.Handle("/static/*", http.StripPrefix("/static/", fs))

	mux.Route("/admin", func(mux chi.Router) {
		mux.Use(Auth)
		mux.Use(AdminAuth)
		mux.Get("/dashboard", handlers.Repo.AdminDashboard)
		mux.Get("/reservations-all", handlers.Repo.AdminAllReservations)
		mux.Get("/reservations-new", handlers.Repo.AdminDashboard)
		mux.Get("/reservations-calendar", handlers.Repo.AdminNewReservations)
		mux.Get("/reservations/{src}/{id}", handlers.Repo.AdminReservationsDetail)
		mux.Post("/reservations/cancle/{id}", handlers.Repo.AdminCancleReservation)
		mux.Post("/reservations/update/{id}", handlers.Repo.AdminUpdateReservation)
		mux.Get("/reservation-calendar", handlers.Repo.AdminReservationsCalendar)
		mux.Post("/admin-avaibility-check", handlers.Repo.PostAdminReservationsCalendar)
		mux.Get("/all-users", handlers.Repo.AdminGetAllUsers)
		mux.Get("/user/{id}", handlers.Repo.AdminGetUserById)
		mux.Post("/user/update/{id}", handlers.Repo.AdminUpdateUserById)
		mux.Post("/user/delete/{id}", handlers.Repo.AdminDeleteUserById)
	})

	return mux
}
