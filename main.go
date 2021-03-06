package main

import (
	"chat/auth"
	"chat/chatroom"
	"chat/template"
	"chat/trace"
	"flag"
	"log"
	"net/http"
	"os"
)

func main() {
	var addr = flag.String("addr", ":8080", "アプリケーションのアドレス")
	flag.Parse()

	auth.SetProviders()

	r, runner := chatroom.NewRoom()
	runner.Tracer(trace.New(os.Stdout))

	// chat.htmlでWebSocketを生成している
	http.Handle("/chat", auth.MustAuth(template.New("chat.html")))
	http.Handle("/login", template.New("login.html"))
	http.HandleFunc("/auth/", auth.LoginHandler)
	http.Handle("/room", r)

	// ログアウト処理　クッキーを削除している
	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:   "auth",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		})
		w.Header()["Location"] = []string{"/chat"}
		w.WriteHeader(http.StatusTemporaryRedirect)
	})

	// チャットルームを開始する
	runner.Run()

	// Webサーバを起動
	log.Println("Webサーバを開始 Port:", *addr)
	err := http.ListenAndServe(*addr, nil)
	if err == nil {
		log.Fatalln("failed listenandserve ", err)
		os.Exit(0)
	}
}
