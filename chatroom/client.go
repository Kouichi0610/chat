package chatroom

import (
	"github.com/gorilla/websocket"
)

// チャットを行う一人のユーザ
type client struct {
	socket *websocket.Conn //
	send   chan []byte     // メッセージが送られるチャネル
	room   *room           // 参加しているチャットルーム
}

func (c *client) read() {
	defer c.socket.Close()

	for {
		_, msg, err := c.socket.ReadMessage()
		if err != nil {
			return
		}
		c.room.forward <- msg
	}
}

func (c *client) write() {
	defer c.socket.Close()
	for msg := range c.send {
		err := c.socket.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			return
		}
	}
}
