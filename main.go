package main

import (
	"chat/auth"
	"chat/chatroom"
	"chat/template"
	"flag"
	"log"
	"net/http"
	"os"
)

func main() {
	var addr = flag.String("addr", ":8080", "アプリケーションのアドレス")
	flag.Parse()

	http.Handle("/chat", auth.MustAuth(template.New("chat.html")))
	// chat.htmlでWebSocketを生成している
	r, runner := chatroom.NewRoom()
	//r.tracer = trace.New(os.Stdout)
	http.Handle("/room", r)
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
