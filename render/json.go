package render

import (
	"encoding/json"
	"net/http"
)

type JSONRender struct {
	Data any
}

func (j *JSONRender) Render(w http.ResponseWriter) error {
	j.WriteContentType(w)

	jd, err := json.Marshal(j.Data)
	if err != nil {
		return err
	}

	_, err = w.Write(jd)
	return err
}

func (j *JSONRender) WriteContentType(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
}
