// テンプレートハンドラ
package template

import (
	"net/http"
	"path/filepath"
	"sync"
	"text/template"
)

func New(path string) http.Handler {
	return &templateHandler{
		path: path,
	}
}

type templateHandler struct {
	once sync.Once
	path string
	tmpl *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		path := filepath.Join("templates", t.path)
		t.tmpl = template.Must(template.ParseFiles(path))
	})
	t.tmpl.Execute(w, r)
}
