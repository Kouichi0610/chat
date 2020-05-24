package auth

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/facebook"
	"github.com/stretchr/gomniauth/providers/github"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/objx"
)

func SetProviders() {
	// セキュリティキーは任意
	gomniauth.SetSecurityKey("k098Biub7nZ")
	gomniauth.WithProviders(
		facebook.New("", "", "http://localhost:8080/auth/callback/facebook"),
		github.New("", "", "http://localhost:8080/auth/callback/github"),
		google.New("419582938575-c208rt5d9vvbn269q8o4j7iijb01b05t.apps.googleusercontent.com", "K4LnW-9k6ZT79RhA9Y35MzQ6", "http://localhost:8080/auth/callback/google"),
	)
}

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
	prv := segs[3]
	switch action {
	case "login":
		provider, err := gomniauth.Provider(prv)
		if err != nil {
			log.Fatalln("認証プロバイダの取得に失敗 ", prv, "-", err)
		}
		loginUrl, err := provider.GetBeginAuthURL(nil, nil)
		if err != nil {
			log.Fatalln("GetBeginAuthURLに失敗 ", prv, "-", err)
		}
		w.Header().Set("Location", loginUrl)
		w.WriteHeader(http.StatusTemporaryRedirect)
	case "callback":
		provider, err := gomniauth.Provider(prv)
		if err != nil {
			log.Fatalln("認証プロバイダの取得に失敗 ", prv, "-", err)
		}
		creds, err := provider.CompleteAuth(objx.MustFromURLQuery(r.URL.RawQuery))
		if err != nil {
			log.Fatalln("認証を完了できませんでした ", prv, "-", err)
		}
		user, err := provider.GetUser(creds)
		if err != nil {
			log.Fatalln("ユーザの取得に失敗しました。 ", prv, "-", err)
		}
		authCookieValue := objx.New(map[string]interface{}{
			"name": user.Name(),
		}).MustBase64()
		http.SetCookie(w, &http.Cookie{
			Name:  "auth",
			Value: authCookieValue,
			Path:  "/",
		})

		w.Header()["Location"] = []string{"/chat"}
		w.WriteHeader(http.StatusTemporaryRedirect)
	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "%s には非対応", action)
	}
}
