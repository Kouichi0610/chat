// テンプレートハンドラ
package template

import (
	"fmt"
	"net/http"
	"path/filepath"
	"sync"
	"text/template"

	"github.com/stretchr/objx"
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

	data := map[string]interface{}{
		"Host": r.Host,
	}
	authCookie, err := r.Cookie("auth")
	if err == nil {
		data["UserData"] = objx.MustFromBase64(authCookie.Value)
		fmt.Printf("User:%s\n", data["UserData"])
	}
	t.tmpl.Execute(w, data)
}
