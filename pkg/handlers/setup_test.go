package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/Man-Crest/GO-Projects/01_bookings/pkg/config"
	"github.com/Man-Crest/GO-Projects/01_bookings/pkg/connection"
	"github.com/Man-Crest/GO-Projects/01_bookings/pkg/render"
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

var app config.AppConfig
var DB connection.DB
var session *scs.SessionManager

func getRoutes() http.Handler {
	app.InProduction = false

	// set up the session
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	tc, err := CreateTestTemplateCache()

	// log.Fatal("tc ------", tc)
	if err != nil {
		log.Fatal("cannot create template cache in test")
	}

	app.TemplateCache = tc
	app.UseCache = true

	repo := NewRepo(&app, DB)
	NewHandlers(repo)

	render.NewTemplates(&app)

	//routes

	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	// mux.Use(NoSurf)
	mux.Use(SessionLoad)

	mux.Get("/", Repo.Home)
	mux.Get("/about", Repo.About)
	mux.Get("/contact", Repo.Contact)
	mux.Get("/make-reservation", Repo.Reservation)
	mux.Post("/make-reservation", Repo.MakeReservation)
	mux.Get("/reservation-summary", Repo.ReservationSummary)
	mux.Get("/generals-quarters", Repo.Generals)
	mux.Get("/majors-suite", Repo.Majors)
	mux.Get("/search-availability", Repo.SearchAvailability)
	mux.Post("/search-availability", Repo.PostAvailability)

	fs := http.FileServer(http.Dir("./static/"))

	mux.Handle("/static/*", http.StripPrefix("/static/", fs))
	return mux

}

func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}

var functions = template.FuncMap{}

var PathToTemplate = "./../../templates"

func CreateTestTemplateCache() (map[string]*template.Template, error) {

	myCache := map[string]*template.Template{}

	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.templ", PathToTemplate))
	if err != nil {
		return myCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.templ", PathToTemplate))
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.templ", PathToTemplate))
			if err != nil {
				return myCache, err
			}
		}
		myCache[name] = ts
	}
	return myCache, nil
}
