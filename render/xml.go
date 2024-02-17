package render

import (
	"encoding/xml"
	"net/http"
)

type XMLRender struct {
	Data any
}

func (x *XMLRender) Render(w http.ResponseWriter) error {
	x.WriteContentType(w)

	err := xml.NewEncoder(w).Encode(x.Data)
	return err
}

func (x *XMLRender) WriteContentType(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
}
