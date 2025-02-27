package messagetypecontroller

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"

	msgModel "messenger_engine/models/message"
	Utils "messenger_engine/utils"
)

// MessageRepository defines the database interaction methods.
type MessageRepository interface {
	LoadMessages(ctx context.Context, chatID int) ([]msgModel.Message, error)
	SaveMessage(ctx context.Context, message msgModel.Message) error
	SaveMessageReply(ctx context.Context, reply msgModel.MessageReply) error
}

// WebSocketManager defines WebSocket client management.
type WebSocketManager interface {
	BroadcastMessage(finalMsg msgModel.FinalMessage)
	BroadcastReply(finalReply msgModel.FinalMessageReply)
}

// MessageService encapsulates message logic.
type MessageService struct {
	Repo MessageRepository
}

// NewMessageService creates a new instance of MessageService.
func NewMessageService(repo MessageRepository) *MessageService {
	return &MessageService{Repo: repo}
}

// FetchChatMessages retrieves messages for a chat.
func (ms *MessageService) FetchChatMessages(chatID int) ([]msgModel.Message, error) {
	return ms.Repo.LoadMessages(context.Background(), chatID)
}

// StoreMessage saves a new message.
func (ms *MessageService) StoreMessage(message msgModel.Message) error {
	return ms.Repo.SaveMessage(context.Background(), message)
}

// StoreMessageReply saves a message reply.
func (ms *MessageService) StoreMessageReply(reply msgModel.MessageReply) error {
	return ms.Repo.SaveMessageReply(context.Background(), reply)
}

// MessageParser extracts structured message data from JSON.
type MessageParser struct{}

// NewMessageParser creates a new instance of MessageParser.
func NewMessageParser() *MessageParser {
	return &MessageParser{}
}

// ExtractChatID extracts chat_id from a message.
func (mp *MessageParser) ExtractChatID(msg map[string]interface{}) (int, error) {
	chatIDFloat, ok := msg["chat_id"].(float64)
	if !ok {
		return 0, fmt.Errorf("invalid chat_id format")
	}
	return int(chatIDFloat), nil
}

// ExtractMessage extracts a new message.
func (mp *MessageParser) ExtractMessage(msg map[string]interface{}) (msgModel.Message, error) {
	messageData, ok := msg["message"].(map[string]interface{})
	if !ok {
		return msgModel.Message{}, fmt.Errorf("invalid message format")
	}
	return parseMessageData(messageData)
}

// ExtractMessageReply extracts a message reply.
func (mp *MessageParser) ExtractMessageReply(msg map[string]interface{}) (msgModel.MessageReply, error) {
	messageReplyData, ok := msg["message"].(map[string]interface{})
	if !ok {
		return msgModel.MessageReply{}, fmt.Errorf("invalid message reply format")
	}
	return parseMessageReplyData(messageReplyData)
}

// MessageTypeHandler processes WebSocket messages.
type MessageTypeHandler struct {
	MessageSvc  *MessageService
	Parser      *MessageParser
	WsManager   WebSocketManager
}

// NewMessageTypeHandler creates a new instance of MessageTypeHandler.
func NewMessageTypeHandler(repo MessageRepository, wsManager WebSocketManager) *MessageTypeHandler {
	return &MessageTypeHandler{
		MessageSvc: NewMessageService(repo),
		Parser:     NewMessageParser(),
		WsManager:  wsManager,
	}
}

// HandleInitialMessage processes messages of type "initial".
func (mh *MessageTypeHandler) HandleInitialMessage(ws *websocket.Conn, msg map[string]interface{}) {
	chatID, err := mh.Parser.ExtractChatID(msg)
	if err != nil {
		sendErrorResponse(ws, err.Error())
		return
	}

	messages, err := mh.MessageSvc.FetchChatMessages(chatID)
	if err != nil {
		log.Printf("Error loading messages: %s", err)
		sendErrorResponse(ws, fmt.Sprintf("Error loading messages: %s", err))
		return
	}

	ws.WriteJSON(map[string]interface{}{
		"type":     "initial",
		"messages": messages,
	})
}

// HandleNewMessage processes new messages.
func (mh *MessageTypeHandler) HandleNewMessage(ws *websocket.Conn, msg map[string]interface{}) {
	messageData, err := mh.Parser.ExtractMessage(msg)
	if err != nil {
		sendErrorResponse(ws, err.Error())
		return
	}

	if err = mh.MessageSvc.StoreMessage(messageData); err != nil {
		log.Printf("Error saving message: %s", err)
		sendErrorResponse(ws, fmt.Sprintf("Error saving message: %s", err))
		return
	}

	finalMsg := msgModel.FinalMessage{Type: "message", Message: messageData}
	mh.WsManager.BroadcastMessage(finalMsg)
}

// HandleMessageReply processes message replies.
func (mh *MessageTypeHandler) HandleMessageReply(ws *websocket.Conn, msg map[string]interface{}) {
	messageReplyData, err := mh.Parser.ExtractMessageReply(msg)
	if err != nil {
		sendErrorResponse(ws, err.Error())
		return
	}

	if err = mh.MessageSvc.StoreMessageReply(messageReplyData); err != nil {
		log.Printf("Error saving reply: %s", err)
		sendErrorResponse(ws, fmt.Sprintf("Error saving reply: %s", err))
		return
	}

	finalReply := msgModel.FinalMessageReply{Type: "message_reply", Message: messageReplyData}
	mh.WsManager.BroadcastReply(finalReply)
}

// sendErrorResponse sends an error message to the WebSocket client.
func sendErrorResponse(ws *websocket.Conn, errorMsg string) {
	ws.WriteJSON(map[string]string{"error": errorMsg})
}

func parseMessageData(data map[string]interface{}) (msgModel.Message, error) {
	messageID, err := Utils.ExtractInt(data, "MessageId")
	if err != nil {
		return msgModel.Message{}, err
	}

	authorID, err := Utils.ExtractInt(data, "AuthorId")
	if err != nil {
		return msgModel.Message{}, err
	}

	timestamp, err := Utils.ExtractInt(data, "Timestamp")
	if err != nil {
		return msgModel.Message{}, err
	}

	receiverID, err := Utils.ExtractInt(data, "ReceiverId")
	if err != nil {
		return msgModel.Message{}, err
	}

	chatID, err := Utils.ExtractInt(data, "ChatId")
	if err != nil {
		return msgModel.Message{}, err
	}

	text, err := Utils.ExtractString(data, "Message")
	if err != nil {
		return msgModel.Message{}, err
	}

	isEdited := Utils.ExtractBool(data, "IsEdited", false)

	return msgModel.Message{
		MessageId:  messageID,
		AuthorId:   authorID,
		Timestamp:  time.Unix(int64(timestamp), 0),
		ReceiverId: receiverID,
		Message:    text,
		ChatId:     chatID,
		IsEdited:   isEdited,
	}, nil
}


func parseMessageReplyData(data map[string]interface{}) (msgModel.MessageReply, error) {
	messageID, err := Utils.ExtractInt(data, "MessageId")
	if err != nil {
		return msgModel.MessageReply{}, err
	}

	authorID, err := Utils.ExtractInt(data, "AuthorId")
	if err != nil {
		return msgModel.MessageReply{}, err
	}

	timestamp, err := Utils.ExtractInt(data, "Timestamp")
	if err != nil {
		return msgModel.MessageReply{}, err
	}

	receiverID, err := Utils.ExtractInt(data, "ReceiverId")
	if err != nil {
		return msgModel.MessageReply{}, err
	}

	chatID, err := Utils.ExtractInt(data, "ChatId")
	if err != nil {
		return msgModel.MessageReply{}, err
	}

	text, err := Utils.ExtractString(data, "Message")
	if err != nil {
		return msgModel.MessageReply{}, err
	}

	isEdited := Utils.ExtractBool(data, "IsEdited", false)

	parentMessageID, err := Utils.ExtractInt(data, "ParentMessageId")
	if err != nil {
		return msgModel.MessageReply{}, err
	}

	return msgModel.MessageReply{
		MessageId:       messageID,
		AuthorId:        authorID,
		Timestamp:       time.Unix(int64(timestamp), 0),
		ReceiverId:      receiverID,
		Message:         text,
		ChatId:          chatID,
		IsEdited:        isEdited,
		ParentMessageId: parentMessageID,
	}, nil
}
