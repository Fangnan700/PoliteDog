package render

import (
	"github.com/fangnan700/PoliteDog/internal/bytesconv"
	"html/template"
	"net/http"
)

type HTMLRender struct {
	Name     string
	Data     any
	Template *template.Template
	IsTmpl   bool
}

func (h *HTMLRender) Render(w http.ResponseWriter) error {
	h.WriteContentType(w)

	if h.IsTmpl {
		// 使用模板
		err := h.Template.ExecuteTemplate(w, h.Name, h.Data)
		if err != nil {
			return err
		}
	} else {
		// 不使用模板
		_, err := w.Write(bytesconv.StringToBytes(h.Data.(string)))
		if err != nil {
			return err
		}
	}

	return nil
}

func (h *HTMLRender) WriteContentType(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
}
