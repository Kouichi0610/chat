package chatroom

import (
	"time"

	"github.com/gorilla/websocket"
)

// チャットを行う一人のユーザ
type client struct {
	socket   *websocket.Conn //
	send     chan *message   // メッセージが送られるチャネル
	room     *room           // 参加しているチャットルーム
	userData map[string]interface{}
}

func (c *client) read() {
	defer c.socket.Close()

	for {
		var msg *message
		err := c.socket.ReadJSON(&msg)
		if err != nil {
			return
		}
		msg.When = time.Now()
		msg.Name = c.userData["name"].(string)
		c.room.forward <- msg
	}
}

func (c *client) write() {
	defer c.socket.Close()
	for msg := range c.send {
		err := c.socket.WriteJSON(msg)
		if err != nil {
			return
		}
	}
}
