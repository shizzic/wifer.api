package chat

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/lxzan/gws"
)

type Handler struct {
	ID int // user_id
}

const (
	PingInterval = 5 * time.Second
	PingWait     = 10 * time.Second
)

var clients sync.Map

// Подключаю юзера к сокету
func Connect(w http.ResponseWriter, r *http.Request, id int) {
	var upgrader = gws.NewUpgrader(&Handler{id}, &gws.ServerOption{
		ParallelEnabled:   true,                                 // Parallel message processing
		Recovery:          gws.Recovery,                         // Exception recovery
		PermessageDeflate: gws.PermessageDeflate{Enabled: true}, // Enable compression
	})
	socket, _ := upgrader.Upgrade(w, r)
	go func() {
		socket.ReadLoop() // Blocking prevents the context from being GC.
	}()
}

func (c *Handler) OnMessage(socket *gws.Conn, message *gws.Message) {
	_ = socket.SetDeadline(time.Now().Add(PingInterval + PingWait))

	// fmt.Print("\n", message.Data)
	defer message.Close()
	socket.WriteMessage(message.Opcode, message.Bytes())
}

func (c *Handler) OnOpen(socket *gws.Conn) {
	_ = socket.SetDeadline(time.Now().Add(PingInterval + PingWait))

	// Удаляю уже существующее соединение для того же юзера
	if v, exist := clients.LoadOrStore(c.ID, socket); exist {
		client := v.(*gws.Conn)
		clients.Delete(c.ID)
		client.NetConn().Close() // Закрываю соеденение
		clients.Store(c.ID, socket)
		fmt.Print("\n", "Closed existed")
	}
	fmt.Print("\n", "Opened")
}

func (c *Handler) OnClose(socket *gws.Conn, err error) {
	quit(c.ID)
	fmt.Print("\n", "Closed")
}

func (c *Handler) OnPing(socket *gws.Conn, payload []byte) {
	_ = socket.SetDeadline(time.Now().Add(PingInterval + PingWait))
	_ = socket.WritePong(nil)
}
func (c *Handler) OnPong(socket *gws.Conn, payload []byte) {}
