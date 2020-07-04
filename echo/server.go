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

	_ = conn.WriteMessage(websocket.TextMessage, []byte("hello"))
	_ = conn.WriteMessage(websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseNormalClosure, "bye"))
}

func attachLogHandler(conn *websocket.Conn) {
	conn.SetPongHandler(logPongMessage)
	conn.SetPingHandler(logPingMessage)
	conn.SetCloseHandler(logCloseMessage)
}

func logCloseMessage(code int, text string) error {
	log.Printf("close, code: %d, text: %s\n", code, text)
	return nil
}

type closeHandler func(int, string) error
type pingHandler func(string) error

func sendBackClose(conn *websocket.Conn) closeHandler {
	return func(i int, s string) error {
		err := conn.WriteControl(websocket.CloseMessage,
			websocket.FormatCloseMessage(i, s), time.Now().Add(2*time.Second))
		if err != nil {
			log.Println("write control: ", err)
		}
		return logCloseMessage(i, s)
	}
}

func logPongMessage(data string) error {
	return dumpMessage("pong message: ", data)
}

func logPingMessage(data string) error {
	return dumpMessage("ping message: ", data)
}

func dumpMessage(msg ...string) error {
	log.Println(msg)
	return nil
}

func responseWithPong(conn *websocket.Conn) pingHandler {
	return func(text string) error {
		err := conn.WriteControl(websocket.PongMessage, []byte(text), time.Now().Add(2*time.Second))
		if err != nil {
			log.Println("write control: ", err)
		}
		return dumpMessage("ping message: ", text)
	}
}

func dropControlMessage(r *http.Request) bool {
	query := r.URL.Query()
	_, present := query["drop"]
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

	_ = conn.WriteMessage(websocket.PingMessage, []byte("pmci"))

	go func() {
		<-time.After(3 * time.Second)
		_ = conn.WriteMessage(websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseGoingAway, "after 3 seconds"))
	}()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			return
		}
		_ = dumpMessage("wait ping: ", string(message))
	}
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
	} else {
		conn.SetPingHandler(responseWithPong(conn))
		conn.SetCloseHandler(sendBackClose(conn))
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
	timer := time.NewTimer(10 * time.Second)
	defer timer.Stop()

	for {
		select {
		case t := <-ticker.C:
			err := conn.WriteMessage(websocket.TextMessage, []byte(t.String()))
			if err != nil {
				log.Println("write message: ", err)
				return
			}
		case <-timer.C:
			_ = dumpMessage("10 seconds timed out")
			_ = conn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseGoingAway, "after 10 seconds"), time.Now().Add(time.Second))
			return
		}
	}
}

func main() {
	http.HandleFunc("/echo", echoMessage)
	http.HandleFunc("/hello", echoHello)
	http.HandleFunc("/ping", ping)
	http.HandleFunc("/flood", flood)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
