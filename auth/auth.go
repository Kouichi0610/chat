package auth

import (
	"net/http"
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
