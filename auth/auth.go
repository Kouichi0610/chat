package auth

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

func MustAuth(handler http.Handler) http.Handler {
	return &authHandler{
		next: handler,
	}
}

type authHandler struct {
	next http.Handler
}

func (h *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("auth")
	if err == http.ErrNoCookie {
		// 認証されていなければ"/login"ページにリダイレクトされる
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}
	if err != nil {
		panic(err.Error())
	}
	h.next.ServeHTTP(w, r)
}

/*
	サードパーティへのログイン処理を受け持つ
	パスの形式: /auth/{action}/{provider}
*/
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	segs := strings.Split(r.URL.Path, "/")
	action := segs[2]
	provider := segs[3]
	switch action {
	case "login":
		log.Println("TODO:Login ", provider)
	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "%s には非対応", action)
	}
}
