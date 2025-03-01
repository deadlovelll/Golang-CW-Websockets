package message

import (
	"database/sql"
	"time"
)

// Message represents a chat message structure.
//
// Fields:
//   - MessageId: Unique identifier for the message.
//   - AuthorId: ID of the user who sent the message.
//   - Timestamp: Time when the message was created.
//   - ReceiverId: ID of the user or group receiving the message.
//   - Message: The actual message content.
//   - ChatId: ID of the chat where the message belongs.
//   - IsEdited: Indicates if the message has been edited.
//   - ParentMessageId: Optional ID of the parent message (for replies).
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

// MessageReply represents a reply to an existing message.
//
// Fields:
//   - MessageId: Unique identifier for the reply message.
//   - AuthorId: ID of the user who sent the reply.
//   - Timestamp: Time when the reply was created.
//   - ReceiverId: ID of the user or group receiving the reply.
//   - Message: The content of the reply message.
//   - ChatId: ID of the chat where the reply belongs.
//   - IsEdited: Indicates if the reply has been edited.
//   - ParentMessageId: ID of the original message being replied to.
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
//
// Fields:
//   - Message: The actual message data.
//   - Type: Specifies the type of message (e.g., "text", "image").
type FinalMessage struct {
	Message Message `json:"message"`
	Type    string  `json:"type"`
}

// FinalMessageReply represents the final reply message format to be sent to the client.
//
// Fields:
//   - Message: The actual reply message data.
//   - Type: Specifies the type of reply (e.g., "reply").
type FinalMessageReply struct {
	Message MessageReply `json:"message"`
	Type    string       `json:"type"`
}
