package test

import (
	"cometosee/controller"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func mustMarshal(v interface{}) json.RawMessage {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return b
}

func TestWebSocket(t *testing.T) {
	manager := controller.NewManger(context.TODO())

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		manager.ServeWS(w, r)
	}))
	defer server.Close()

	u := url.URL{
		Scheme: "ws",
		Host:   server.Listener.Addr().String(),
		Path:   "/",
	}

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		t.Fatalf("dial error: %v", err)
	}
	defer conn.Close()

	t.Log("user is registering..")
	// register user
	if err := conn.WriteJSON(controller.Event{
		Type: controller.EventRegister,
		Payload: mustMarshal(controller.RegisterEvent{
			Name: "barun",
			Room: "cricket",
		}),
	}); err != nil {
		t.Fatal(err)
	}
	t.Log("user registered")
	t.Log("user is sending msg")
	// semd msg
	if err := conn.WriteJSON(controller.Event{
		Type: controller.EventSendMessage,
		Payload: mustMarshal(controller.SendMessageEvent{
			From:    "barun",
			Message: "welcome to cricket commiunity",
		}),
	}); err != nil {
		t.Fatal(err)
	}
	t.Log("user send message")
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))

	var resp controller.Event
	if err := conn.ReadJSON(&resp); err != nil {
		t.Fatal(err)
	}

	if resp.Type != controller.EventNewMessage {
		t.Fatalf("expected %s, got %s", controller.EventNewMessage, resp.Type)
	}
}
