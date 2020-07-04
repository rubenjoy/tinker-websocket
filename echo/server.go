package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  2048,
	WriteBufferSize: 2048,
}

// echoHello upgrades the HTTP request to websocket connection
// then send 'hello' message and send close control message
func echoHello(w http.ResponseWriter, r *http.Request) {
	log.Println("accept hello request...")

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("ugprade: ", err)
	}
	defer conn.Close()
	if dropControlMessage(r) {
		attachLogHandler(conn)
	}

	conn.WriteMessage(websocket.TextMessage, []byte("hello"))
	conn.WriteMessage(websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseNormalClosure, "bye"))
}

func loopRead(conn *websocket.Conn, tag string) {
	for {
		mt, msg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("tag: '%s', read message: '%s'", tag, err)
			return
		}
		if websocket.TextMessage == mt {
			log.Printf("tag: '%s', message: %s", tag, msg)
		}
	}
}

func attachLogHandler(conn *websocket.Conn) {
	conn.SetPongHandler(logPongMessage)
	conn.SetPingHandler(logPingMessage)
	conn.SetCloseHandler(dumpCloseMessage)
}

func dumpCloseMessage(code int, text string) error {
	log.Printf("close, code: %d, text: %s\n", code, text)
	return nil
}

type closeHandler func(int, string) error
type pongHandler func(string) error
type pingHandler func(string) error

func sendBackClose(conn *websocket.Conn) closeHandler {
	return func(i int, s string) error {
		conn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(i, s), time.Now().Add(2*time.Second))
		return dumpCloseMessage(i, s)
	}
}

func logPongMessage(data string) error {
	log.Println("pong message: ", data)
	return nil
}

func sendBackPong(conn *websocket.Conn) pongHandler {
	return func(text string) error {
		conn.WriteMessage(websocket.PongMessage, []byte(text))
		return logPongMessage(text)
	}
}

func logPingMessage(data string) error {
	log.Println("ping message: ", data)
	return nil
}

func dumpMessage(msg string) error {
	log.Println(msg)
	return nil
}

func responseWithPong(conn *websocket.Conn) pingHandler {
	return func(text string) error {
		return logPingMessage(text)
	}
}

func dropControlMessage(r *http.Request) bool {
	query := r.URL.Query()
	_, present := query["dropcontrol"]
	return present
}

// ping upgrades the HTTP request to websocket connection
// then send ping control message and wait the pong response in 10 seconds
// finally it sends close control message
func ping(w http.ResponseWriter, r *http.Request) {
	log.Println("accept ping request...")

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade: ", err)
	}
	defer conn.Close()
	if dropControlMessage(r) {
		attachLogHandler(conn)
	}

	conn.WriteMessage(websocket.PingMessage, []byte("pmci"))
	conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "bye"))
}

// echoMessage upgrades the HTTP request into a websocket
// and will echo back the received message.
// When error, it closes the connection.
func echoMessage(w http.ResponseWriter, r *http.Request) {
	log.Println("accept echo request")

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade: ", err)
	}
	defer conn.Close()
	if dropControlMessage(r) {
		attachLogHandler(conn)
	}

	for {
		mt, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("read message: ", err)
			break
		}
		if websocket.TextMessage == mt {
			log.Println("echo message: ", string(msg))
		}
		err = conn.WriteMessage(mt, msg)
		if err != nil {
			log.Println("write message: ", err)
			break
		}
	}
}

func flood(w http.ResponseWriter, r *http.Request) {
	log.Println("accepting flood request...")

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade: ", err)
	}
	defer conn.Close()
	if dropControlMessage(r) {
		attachLogHandler(conn)
	}

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

}

func main() {
	http.HandleFunc("/echo", echoMessage)
	http.HandleFunc("/hello", echoHello)
	http.HandleFunc("/ping", ping)
	http.HandleFunc("/flood", flood)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
