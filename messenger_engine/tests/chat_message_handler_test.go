// file: chatmessagehandler_test.go
package tests

import (
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"messenger_engine/controllers/broadcast_controller" // assumed package path for broadcast type
	"messenger_engine/controllers/websocket_controller/handlers/chat_message_handler"
	"messenger_engine/controllers/message_controller"
	"messenger_engine/models/message"

	"github.com/gorilla/websocket"
)

// DummyMessageController implements required methods for testing.
type DummyMessageController struct{}

func (d *DummyMessageController) SaveMessage(msg message.Message) error {
	return nil
}

func (d *DummyMessageController) SaveMessageReply(msg message.MessageReply) error {
	return nil
}

func (d *DummyMessageController) LoadMessages(chatId int) ([]message.Message, error) {
	// Return one dummy message for testing.
	return []message.Message{
		{
			MessageId:  1,
			Message:    "Hello",
			AuthorId:   1,
			ChatId:     chatId,
			ReceiverId: 2,
			Timestamp:  time.Now(),
			IsEdited:   false,
		},
	}, nil
}

func TestChatMessageHandler_InitialMessage(t *testing.T) {
	// Create dummy dependencies.
	dummyMsgCtrl := &messagecontroller.MessageController{}
	dummyBroadcast := &broadcastcontroller.Broadcast{
		Clients:          make(map[*websocket.Conn]bool),
		Broadcast:        make(chan message.FinalMessage, 1),
		RepliesBroadcast: make(chan message.FinalMessageReply, 1),
	}
	upgrader := websocket.Upgrader{}

	// Create the ChatMessageHandler.
	handler := chatmessagehandler.NewChatMessageHandler(upgrader, dummyMsgCtrl, dummyBroadcast)

	// Set up a test server.
	server := httptest.NewServer(handler)
	defer server.Close()

	wsURL, err := url.Parse(server.URL)
	if err != nil {
		t.Fatalf("failed to parse server URL: %v", err)
	}
	wsURL.Scheme = "ws"

	// Dial the WebSocket.
	ws, _, err := websocket.DefaultDialer.Dial(wsURL.String(), nil)
	if err != nil {
		t.Fatalf("failed to connect to WebSocket: %v", err)
	}
	defer ws.Close()

	// Send an "initial" message with a chat_id.
	request := map[string]interface{}{
		"type":    "initial",
		"chat_id": 10,
	}
	if err := ws.WriteJSON(request); err != nil {
		t.Fatalf("failed to write JSON: %v", err)
	}

	// Read the JSON response.
	var response map[string]interface{}
	if err := ws.ReadJSON(&response); err != nil {
		t.Fatalf("failed to read JSON response: %v", err)
	}

	// Validate that the response has type "initial" and includes messages.
	if response["type"] != "initial" {
		t.Errorf("expected response type 'initial', got: %v", response["type"])
	}
	if _, ok := response["messages"]; !ok {
		t.Errorf("expected response to contain key 'messages'")
	}
}

func TestChatMessageHandler_MessageAndReply(t *testing.T) {
	// Create dummy dependencies.
	dummyMsgCtrl := &messagecontroller.MessageController{}
	dummyBroadcast := &broadcastcontroller.Broadcast{
		Clients:          make(map[*websocket.Conn]bool),
		Broadcast:        make(chan message.FinalMessage, 1),
		RepliesBroadcast: make(chan message.FinalMessageReply, 1),
	}
	upgrader := websocket.Upgrader{}
	handler := chatmessagehandler.NewChatMessageHandler(upgrader, dummyMsgCtrl, dummyBroadcast)

	// Set up a test server.
	server := httptest.NewServer(handler)
	defer server.Close()

	wsURL, err := url.Parse(server.URL)
	if err != nil {
		t.Fatalf("failed to parse server URL: %v", err)
	}
	wsURL.Scheme = "ws"

	// Connect the client.
	ws, _, err := websocket.DefaultDialer.Dial(wsURL.String(), nil)
	if err != nil {
		t.Fatalf("failed to connect to WebSocket: %v", err)
	}
	defer ws.Close()

	// --- Test "message" type ---
	testMessage := map[string]interface{}{
		"type": "message",
		"message": map[string]interface{}{
			"MessageId":  float64(1),
			"AuthorId":   float64(1),
			"Timestamp":  float64(time.Now().Unix()),
			"ReceiverId": float64(2),
			"Message":    "Test message",
			"ChatId":     float64(10),
			"IsEdited":   false,
		},
	}
	if err := ws.WriteJSON(testMessage); err != nil {
		t.Fatalf("failed to write JSON message: %v", err)
	}
	// The "message" type does not send a response directly. Instead, check the broadcast channel.
	select {
	case finalMsg := <-dummyBroadcast.Broadcast:
		if finalMsg.Type != "message" {
			t.Errorf("expected broadcast message type 'message', got: %s", finalMsg.Type)
		}
	case <-time.After(time.Second):
		t.Errorf("timeout waiting for broadcast message")
	}

	// --- Test "message_reply" type ---
	testReply := map[string]interface{}{
		"type": "message_reply",
		"message": map[string]interface{}{
			"MessageId":       float64(2),
			"AuthorId":        float64(1),
			"Timestamp":       float64(time.Now().Unix()),
			"ReceiverId":      float64(2),
			"Message":         "Test reply",
			"ChatId":          float64(10),
			"IsEdited":        false,
			"ParentMessageId": float64(1),
		},
	}
	if err := ws.WriteJSON(testReply); err != nil {
		t.Fatalf("failed to write JSON reply: %v", err)
	}
	select {
	case finalReply := <-dummyBroadcast.RepliesBroadcast:
		if finalReply.Type != "message_reply" {
			t.Errorf("expected broadcast reply type 'message_reply', got: %s", finalReply.Type)
		}
	case <-time.After(time.Second):
		t.Errorf("timeout waiting for broadcast reply")
	}
}
