package forms

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type Form struct {
	url.Values
	Errors errors
}

func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

func (f *Form) Has(field string, r *http.Request) bool {
	s := f.Get(field)
	if s == "" {
		f.Errors.Add(field, fmt.Sprintf("%s is invalid", field))
		return false
	} else {
		if !f.MinLength(field, r) {
			f.Errors.Add(field, fmt.Sprintf("%s must have minimum 3 chars", field))
		}
	}
	return true
}

func (f *Form) MinLength(field string, r *http.Request) bool {
	trrimed := strings.TrimSpace(f.Get(field))
	if len(trrimed) <= 3 {
		return false
	}
	return true
}
func (f *Form) Required(fields ...string) bool {

	for _, field := range fields {
		if len(field) > 0 {
			return true
		}
	}
	return false
}
