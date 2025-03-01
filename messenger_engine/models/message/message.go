package message

import (
	"database/sql"
	"time"
)

// Mesaage represents a message structure.
type Message struct {
	MessageId       int           `json:"message_id"`
	AuthorId        int           `json:"author_id"`
	Timestamp       time.Time     `json:"timestamp"`
	ReceiverId      int           `json:"receiver_id"`
	Message         string        `json:"message"`
	ChatId          int           `json:"chat_id"`
	IsEdited        bool          `json:"is_edited"`
	ParentMessageId sql.NullInt64 `json:"parent_message_id"`
}

// MessageReply represents a reply to a message.
type MessageReply struct {
	MessageId       int       `json:"message_id"`
	AuthorId        int       `json:"author_id"`
	Timestamp       time.Time `json:"timestamp"`
	ReceiverId      int       `json:"receiver_id"`
	Message         string    `json:"message"`
	ChatId          int       `json:"chat_id"`
	IsEdited        bool      `json:"is_edited"`
	ParentMessageId int       `json:"parent_message_id"`
}

// FinalMessage represents the final message format to be sent to the client.
type FinalMessage struct {
	Message Message `json:"message"`
	Type    string  `json:"type"`
}

// FinalMessageReply represents the final reply format to be sent to the client.
type FinalMessageReply struct {
	Message MessageReply `json:"message"`
	Type    string       `json:"type"`
}