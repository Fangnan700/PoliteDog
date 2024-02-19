package binding

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type jsonBinding struct {
	DisallowUnknownFields bool
}

func (j *jsonBinding) Name() string {
	return "json"
}

func (j *jsonBinding) Bind(r *http.Request, obj any) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	r.Body = io.NopCloser(bytes.NewReader(body))
	decoder := json.NewDecoder(bytes.NewReader(body))
	// 是否校验json对应结构体字段
	if j.DisallowUnknownFields {
		decoder.DisallowUnknownFields()
	}

	return decoder.Decode(obj)
}
