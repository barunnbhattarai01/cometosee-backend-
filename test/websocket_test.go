package test

import (
	"cometosee/di"
	"cometosee/model"
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
	//di
	wsdi := di.SetupMessage()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wsdi.ServeWS(w, r)
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
	if err := conn.WriteJSON(model.Event{
		Type: model.EventRegister,
		Payload: mustMarshal(model.RegisterEvent{
			Name: "barun",
			Room: "cricket",
		}),
	}); err != nil {
		t.Fatal(err)
	}
	t.Log("user registered")
	t.Log("user is sending msg")
	// semd msg
	if err := conn.WriteJSON(model.Event{
		Type: model.EventSendMessage,
		Payload: mustMarshal(model.SendMessageEvent{
			From:    "barun",
			Message: "welcome to cricket commiunity",
		}),
	}); err != nil {
		t.Fatal(err)
	}
	t.Log("user send message")
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))

	var resp model.Event
	if err := conn.ReadJSON(&resp); err != nil {
		t.Fatal(err)
	}

	if resp.Type != model.EventNewMessage {
		t.Fatalf("expected %s, got %s", model.EventNewMessage, resp.Type)
	}
}
