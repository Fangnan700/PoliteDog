package render

import "net/http"

type Render interface {
	Render(w http.ResponseWriter) error
	WriteContentType(w http.ResponseWriter)
}
