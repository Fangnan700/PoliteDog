package render

import (
	"errors"
	"fmt"
	"net/http"
)

type RedirectRender struct {
	Code     int
	Request  *http.Request
	Location string
}

func (r *RedirectRender) Render(w http.ResponseWriter) error {
	r.WriteContentType(w)

	if r.Code < http.StatusMultipleChoices || r.Code > http.StatusPermanentRedirect && r.Code != http.StatusCreated {
		return errors.New(fmt.Sprintf("Can net redirect with status code: %d", r.Code))
	}

	http.Redirect(w, r.Request, r.Location, r.Code)
	return nil
}

func (r *RedirectRender) WriteContentType(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
}
