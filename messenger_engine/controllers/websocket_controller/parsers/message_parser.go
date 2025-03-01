package parsers

import (
	"fmt"
	"time"

	Messages "messenger_engine/models/message"
)

// Parser encapsulates methods for parsing message data.
type Parser struct{}

// New returns a new instance of Parser.
//
// Returns:
//   - A pointer to a newly initialized Parser.
func New() *Parser {
	return &Parser{}
}

// ParseChatID extracts the chat ID from the incoming message.
//
// Parameters:
//   - msg: A map[string]interface{} representing the incoming JSON message.
//
// Returns:
//   - An integer representing the chat ID.
//   - An error if the chat_id is missing or invalid.
func (p *Parser) ParseChatID(msg map[string]interface{}) (int, error) {
	chatIdFloat, ok := msg["chat_id"].(float64)
	if !ok {
		return 0, fmt.Errorf("invalid chat_id")
	}
	return int(chatIdFloat), nil
}

// ParseMessageData extracts and converts a message from the incoming JSON payload.
//
// Parameters:
//   - msg: A map[string]interface{} representing the incoming JSON message.
//
// Returns:
//   - A Messages.Message struct containing the parsed message data.
//   - An error if the message format is invalid.
func (p *Parser) ParseMessageData(msg map[string]interface{}) (Messages.Message, error) {
	messageData, ok := msg["message"].(map[string]interface{})
	if !ok {
		return Messages.Message{}, fmt.Errorf("invalid message format")
	}

	return Messages.Message{
		MessageId:  int(messageData["MessageId"].(float64)),
		AuthorId:   int(messageData["AuthorId"].(float64)),
		Timestamp:  time.Unix(int64(messageData["Timestamp"].(float64)), 0),
		ReceiverId: int(messageData["ReceiverId"].(float64)),
		Message:    messageData["Message"].(string),
		ChatId:     int(messageData["ChatId"].(float64)),
		IsEdited:   messageData["IsEdited"].(bool),
	}, nil
}

// ParseMessageReplyData extracts and converts a message reply from the incoming JSON payload.
//
// Parameters:
//   - msg: A map[string]interface{} representing the incoming JSON message.
//
// Returns:
//   - A Messages.MessageReply struct containing the parsed reply data.
//   - An error if the message reply format is invalid.
func (p *Parser) ParseMessageReplyData(msg map[string]interface{}) (Messages.MessageReply, error) {
	messageReplyData, ok := msg["message"].(map[string]interface{})
	if !ok {
		return Messages.MessageReply{}, fmt.Errorf("invalid message format")
	}

	return Messages.MessageReply{
		MessageId:       int(messageReplyData["MessageId"].(float64)),
		AuthorId:        int(messageReplyData["AuthorId"].(float64)),
		Timestamp:       time.Unix(int64(messageReplyData["Timestamp"].(float64)), 0),
		ReceiverId:      int(messageReplyData["ReceiverId"].(float64)),
		Message:         messageReplyData["Message"].(string),
		ChatId:          int(messageReplyData["ChatId"].(float64)),
		IsEdited:        messageReplyData["IsEdited"].(bool),
		ParentMessageId: int(messageReplyData["ParentMessageId"].(float64)),
	}, nil
}
