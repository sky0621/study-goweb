package main

import (
	"github.com/gorilla/websocket"
)

// クライアントは、チャットを行っている１人のユーザを表す。
type client struct {
	// このクライアントのためのWebSocket
	socket *websocket.Conn
	// メッセージが送られるチャネル
	send chan []byte
	// このクライアントが参加しているチャットルーム
	room *room
}

func (c *client) read() {
	for {
		if _, msg, err := c.socket.ReadMessage(); err == nil {
			c.room.forward <- msg
		} else {
			break
		}
	}
	c.socket.Close()
}

func (c *client) write() {
	for msg := range c.send {
		if err := c.socket.WriteMessage(websocket.TextMessage, msg);
			err != nil {
			break
		}
	}
	c.socket.Close()
}

