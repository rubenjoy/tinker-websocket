package main

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
		server := httptest.NewServer(http.HandlerFunc(chatHandler))
		defer server.Close()

		wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/chat"
		conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
		if err != nil {
			t.Fatal(err)
		}

		_ = conn.SetReadDeadline(time.Now().Add(2 * time.Second))
		_, message, err := conn.ReadMessage()
		if err != nil {
			t.Fatal(err)
		}
		if string(message) != "hello" {
			t.Errorf("Expected message: %s, got %s", "hello", string(message))
		}
	})
}
