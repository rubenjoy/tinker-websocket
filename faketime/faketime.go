package faketime

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

var upgrade = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

const deadline = time.Second

func chatHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrade.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	defer conn.Close()

	if err := writeText(conn, "hello", deadline); err != nil {
		log.Println(err)
	}
	if err := writeText(conn, "1234", deadline); err != nil {
		log.Println(err)
	}

	if err := conn.WriteControl(websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseNormalClosure, "bye"),
		time.Now().Add(deadline)); err != nil {
		log.Println(err)
	}
}

func writeText(conn *websocket.Conn, message string, d time.Duration) error {
	err := conn.SetWriteDeadline(time.Now().Add(d))
	if err != nil {
		return err
	}
	return conn.WriteMessage(websocket.TextMessage, []byte(message))
}
