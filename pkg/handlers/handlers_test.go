package handlers

import (
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Man-Crest/GO-Projects/01_bookings/pkg/models"
)

type postData struct {
	key   string
	value string
}

var theTests = []struct {
	name               string
	url                string
	method             string
	params             []postData
	expectedStatusCode int
}{
	{"home", "/", "GET", []postData{}, http.StatusOK},
	{"about", "/about", "GET", []postData{}, http.StatusOK},
	{"generals-quarters", "/generals-quarters", "GET", []postData{}, http.StatusOK},
	{"majors-suite", "/majors-suite", "GET", []postData{}, http.StatusOK},
	{"search-availability", "/search-availability", "GET", []postData{}, http.StatusOK},
	// {"search-availability", "/search-availability", "POST", []postData{
	// 	{key: "start", value: "2020-01-01"},
	// 	{key: "end", value: "2020-01-05"},
	// 	{key: "roomType", value: "general"},
	// }, http.StatusOK},
	{"contact", "/contact", "GET", []postData{}, http.StatusOK},
	// 	{"make-reservation", "/make-reservation", "POST", []postData{
	// 		{key: "first_name", value: "man"},
	// 		{key: "last_name", value: "patel"},
	// 		{key: "email", value: "man@gmail.com"},
	// 		{key: "phone", value: "1010101010"},
	// 	}, http.StatusOK},
}

func TestGetHandlers(t *testing.T) {
	routes := getRoutes()

	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	for _, e := range theTests {
		if e.method == "GET" {
			resp, err := ts.Client().Get(ts.URL + e.url)
			if err != nil {
				t.Log(err, "Get methods condition LOG")
				t.Fatalf("route: %s, method: %s , handlers: %s ,err: %s", e.url, e.method, e.name, err)
			}
			if resp.StatusCode != e.expectedStatusCode {
				t.Errorf("for %s expected %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
			}
		}
	}
}

func TestReservation(t *testing.T) {

	reservation := models.Reservation{
		Room: &models.Rooms{
			ID:         1,
			RoomNumber: "401",
			RoomType: &models.RoomType{
				Name: "general",
			},
		},
	}

	req, _ := http.NewRequest("Get", "/make-reservation", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)

	handler := http.HandlerFunc(Repo.Reservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("reservation handler returning response code %d", rr.Code)
	}
}

func getCtx(req *http.Request) context.Context {
	ctx, err := session.Load(req.Context(), req.Header.Get("X-Session"))

	if err != nil {
		log.Println(err)
	}
	return ctx
}

// func TestMakeReservationHandler(t *testing.T) {
// 	// Create a new instance of your Repository with mocked dependencies
// 	repo := &Repository{
// 		App:
// 		DB:
// 	}

// 	// Create a new HTTP request for the handler
// 	req := httptest.NewRequest("POST", "/make-reservation", strings.NewReader("first_name=John&last_name=Doe&email=johndoe@example.com&phone=123456789"))

// 	// Create a new ResponseRecorder (which satisfies http.ResponseWriter) to record the response
// 	rr := httptest.NewRecorder()

// 	// Call the handler function with the ResponseRecorder and Request
// 	repo.MakeReservation(rr, req)

// 	// Check the status code returned by the handler
// 	if status := rr.Code; status != http.StatusSeeOther {
// 		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusSeeOther)
// 	}

// 	// Check if the reservation is properly stored in the session
// 	reservation, ok := repo.App.Session.Get(req.Context(), "reservation").(models.Reservation)
// 	if !ok {
// 		t.Error("expected reservation to be stored in session")
// 	}

// 	// Check the correctness of the reservation details
// 	if reservation.FirstName != "John" || reservation.LastName != "Doe" || reservation.Email != "johndoe@example.com" || reservation.Phone != "123456789" {
// 		t.Errorf("handler did not store correct reservation details in session")
// 	}
// }
