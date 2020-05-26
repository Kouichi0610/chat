package chatroom

import (
	"chat/trace"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/stretchr/objx"
)

type Runner interface {
	Run()
	Tracer(t trace.Tracer)
}

type room struct {
	forward chan *message
	join    chan *client
	leave   chan *client
	clients map[*client]bool
	tracer  trace.Tracer
}

func (r *room) Run() {
	go r.run()
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  socketBufferSize,
	WriteBufferSize: messageBufferSize,
}

func (r *room) ServeHTTP(w http.ResponseWriter, rq *http.Request) {
	socket, err := upgrader.Upgrade(w, rq, nil)
	if err != nil {
		log.Fatal("ServeHTP:", err)
	}

	authCookie, err := rq.Cookie("auth")
	if err != nil {
		log.Fatal("failed get auth cookie ", err)
		return
	}
	client := &client{
		socket:   socket,
		send:     make(chan *message, messageBufferSize),
		room:     r,
		userData: objx.MustFromBase64(authCookie.Value),
	}
	r.join <- client
	defer func() {
		r.leave <- client
	}()
	go client.write()
	client.read()
}

func NewRoom() (http.Handler, Runner) {
	r := &room{
		forward: make(chan *message),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
		tracer:  trace.Empty(),
	}
	return r, r
}
func (r *room) Tracer(t trace.Tracer) {
	r.tracer = t
}

func (r *room) run() {
	for {
		select {
		case client := <-r.join:
			r.clients[client] = true
			r.tracer.Trace("新しいクライアントが参加")
		case client := <-r.leave:
			delete(r.clients, client)
			close(client.send)
			r.tracer.Trace("クライアントが退室")
		case msg := <-r.forward:
			for client := range r.clients {
				client.send <- msg
			}
		}
	}
}
