package chat

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
	"wifer/server/crud/create"
	"wifer/server/crud/delete"
	"wifer/server/structs"

	"github.com/lxzan/gws"
)

type Handler struct {
	ID    int // user_id
	Props *structs.Props
}

const (
	PingInterval = 5 * time.Second
	PingWait     = 10 * time.Second
)

var clients sync.Map

// Подключаю юзера к сокету
func Connect(props *structs.Props, w http.ResponseWriter, r *http.Request, id int) {
	var upgrader = gws.NewUpgrader(&Handler{id, props}, &gws.ServerOption{
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
	_ = socket.SetDeadline(time.Now().Add(PingInterval + PingWait)) // Обновляю остаток времени для следующей проверки соединения

	var msg structs.Message
	if err := json.Unmarshal(message.Bytes(), &msg); err != nil {
		socket.NetConn().Close()
	}

	// Отправляю любое действие связанное с сообщением таргету, если он в сети (письмо/чтение/смена доступа в чате)
	if v, exist := clients.Load(msg.Target); exist {
		client := v.(*gws.Conn)
		client.WriteMessage(message.Opcode, message.Bytes())
	}

	switch msg.Api {
	case "message":
		write(c.Props, &msg)
	case "view":
		view(c.Props, &msg)
	case "access": // сменить доступ к бесплатному чату для собеседника
		var target structs.Target
		target.Target = msg.Target

		if msg.Access {
			create.TargetAccess(c.Props, &target, c.ID)
		} else {
			delete.TargetAccess(c.Props, &target, c.ID)
		}
	}

	defer message.Close()
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
}

func (c *Handler) OnClose(socket *gws.Conn, err error) {
	quit(c.ID)
}

func (c *Handler) OnPing(socket *gws.Conn, payload []byte) {
	_ = socket.SetDeadline(time.Now().Add(PingInterval + PingWait))
	_ = socket.WritePong(nil)
}
func (c *Handler) OnPong(socket *gws.Conn, payload []byte) {}
