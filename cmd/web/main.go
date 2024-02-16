package main

import (
	"encoding/gob"
	"fmt"
	"os"

	"log"
	"net/http"
	"time"

	"github.com/Man-Crest/GO-Projects/01_bookings/pkg/config"
	"github.com/Man-Crest/GO-Projects/01_bookings/pkg/connection"
	"github.com/Man-Crest/GO-Projects/01_bookings/pkg/handlers"
	"github.com/Man-Crest/GO-Projects/01_bookings/pkg/helpers"
	"github.com/Man-Crest/GO-Projects/01_bookings/pkg/models"
	"github.com/Man-Crest/GO-Projects/01_bookings/pkg/render"
	"github.com/alexedwards/scs/v2"
)

const portNumber = ":8000"

var app config.AppConfig
var session *scs.SessionManager

var db *connection.DB
var err error

// main is the main function
func main() {
	db, err = connection.ConnectSQL()
	defer db.SQL.Close()

	defer close(app.MailChan)

	// Create a new repository with the AppConfig
	repo := handlers.NewRepo(&app, *db)
	handlers.NewHandlers(repo)
	helpers.NewHelpers(&app)

	err := run()
	if err != nil {
		fmt.Println(err, "at run function")
	}

	ReadMail()
	fmt.Println(fmt.Sprintf("Staring application on port %s", portNumber))
	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err, "error in serving at port ")
	}

}

func run() error {
	app.InProduction = false

	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register([]models.Rooms{})
	gob.Register(models.Rooms{})

	infolog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infolog
	errorlog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime)
	app.ErorLog = errorlog
	// set up the session
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache", err)
	}

	mailChan := make(chan models.MailData)
	app.MailChan = mailChan

	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app, *db)

	handlers.NewHandlers(repo)

	render.NewTemplates(&app)

	return err
}
