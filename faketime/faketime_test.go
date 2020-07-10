package faketime

import (
	"github.com/gorilla/websocket"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestWebsocket(t *testing.T) {
	t.Run("connect websocket should say hello", func(t *testing.T) {
		server := createServer(chatHandler)
		defer server.Close()

		conn := dialWebsocket(t, server.URL, "/chat")

		message := readText(t, conn, 2*time.Second)
		if message != "hello" {
			t.Errorf("Expected message: %s, got %s", "hello", string(message))
		}
	})

	t.Run("after 'hello', it replied with generated user id",
		func(t *testing.T) {
			server := createServer(chatHandler)
			defer server.Close()

			conn := dialWebsocket(t, server.URL, "/chat")

			_ = readText(t, conn, 2*time.Second)

			userId := readText(t, conn, 2*time.Second)
			if userId != "1234" {
				t.Errorf("Expected user id: %s, got %s", "1234", userId)
			}
		})
}

func createServer(f func(w http.ResponseWriter, r *http.Request)) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(f))
}

func dialWebsocket(t *testing.T, serverURL, path string) *websocket.Conn {
	wsURL := "ws" + strings.TrimPrefix(serverURL, "http") + path
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("could not dial websocket: %v", err)
	}
	return conn
}

func readText(t *testing.T, conn *websocket.Conn, deadline time.Duration) string {
	_ = conn.SetReadDeadline(time.Now().Add(deadline))
	mt, message, err := conn.ReadMessage()
	if err != nil {
		t.Errorf("could not read message: %v", err)
	}
	if mt != websocket.TextMessage {
		t.Errorf("Expected message type: %d, got %d", websocket.TextMessage, mt)
	}
	return string(message)
}
