package forms

import (
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestHas(t *testing.T) {

	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	has := form.Has("whatever", r)

	if has {
		t.Errorf("form shows field when it does not")
	}

	postedData := url.Values{}
	postedData.Add("a", "bazs")
	postedData.Add("c", "dasda")
	postedData.Add("e", "fasd")
	form = New(postedData)

	has = form.Has("a", r)
	if !has {
		t.Errorf("Form have to show the value but it did not have one %s", form)
	}
	has = form.Has("c", r)
	if !has {
		t.Errorf("Form have to show the value but it did not have one %s", form)
	}
	has = form.Has("e", r)
	if !has {
		t.Errorf("Form have to show the value but it did not have one %s", form)
	}
}

func TestValid(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)

	form := New(r.PostForm)
	if !form.Valid() {
		t.Errorf("form does not have any error so that it does'nt show")
	}

	postData := url.Values{}
	postData.Add("a", "bqwe")

	form = New(postData)
	form.Has("b", r)

	if form.Valid() {
		t.Errorf("it supposed to show an error")
	}
}

// func TestNew(t *testing.T) {
// 	r := httptest.NewRequest("POST", "/whatever", nil)

// 	form := New(r.PostForm)

// 	if form.Errors != nil {
// 		t.Errorf("form does not have to show an error field because it is nil: %s", form.Errors)
// 	} else if form.Values != nil {
// 		t.Errorf("form does not have to show a url.Values field because nil : %s", form.Values)
// 	}
// }
