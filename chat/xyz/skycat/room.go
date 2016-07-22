package main

type room struct {
	// 他のクライアントに転送するためのメッセージを保持するチャネル
	forward chan []byte
	// チャットルームに参加しようとしているクライアントのためのチャネル
	join chan *client
	// チャットルームから退室しようとしているクライアントのためのチャネル
	leave chan *client
	// 在室している全てのクライアントを保持
	clients map[*client]bool
}

func newRoom() *room {
	return &room {
		forward: make(chan []byte),
		join: make(chan *client),
		leave: make(chan *client),
		clients: make(map[*client]bool),
	}
}

func (r *room) run() {
	// 無限ループ
	for {
		// 共有されているメモリに対して同期化や変更が必要な任意の箇所で使える。
		select {
			case client := <-r.join:
				// 参加
				r.clients[client] = true
			case client := <-r.leave:
				// 退室
				delete(r.clients, client)
				close(client.send)
			case msg := <-r.forward:
				// 全てのクライアントにメッセージを転送
				for client := range r.clients {
					select {
						case client.send <- msg:
						default:
							// 送信に失敗
							delete(r.clients, client)
							close(client.send)
					}
				}
		}
	}
}

const (
	socketBufferSize = 1024
	messageBufferSize = 256
)
var upgrader = &websocket.Upgrader{ReadBufferSize:
	socketBufferSize, WriteBufferSize: socketBufferSize}
func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP:", err)
		return
	}
	client := &client{
		socket: socket,
		send: make(chan []byte, messageBufferSize),
		room: r,
	}
	r.join <- client
	defer func() { r.leave <- client }()
	go client.write()
	client.read()
}
