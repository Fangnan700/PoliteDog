package render

import (
	"fmt"
	"github.com/fangnan700/PoliteDog/internal/bytesconv"
	"net/http"
)

type StringRender struct {
	Format string
	Data   []any
}

func (s *StringRender) Render(w http.ResponseWriter) error {
	s.WriteContentType(w)

	data := fmt.Sprintf(s.Format, s.Data...)
	_, err := w.Write(bytesconv.StringToBytes(data))

	return err
}

func (s *StringRender) WriteContentType(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
}
