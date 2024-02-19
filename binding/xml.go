package binding

import (
	"bytes"
	"encoding/xml"
	"io"
	"net/http"
)

type xmlBinding struct {
}

func (j *xmlBinding) Name() string {
	return "xml"
}

func (j *xmlBinding) Bind(r *http.Request, obj any) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	r.Body = io.NopCloser(bytes.NewReader(body))
	decoder := xml.NewDecoder(bytes.NewReader(body))

	return decoder.Decode(obj)
}
