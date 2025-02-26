package message

import (
	"database/sql"
	"time"
)

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

type FinalMessage struct {
	Message Message `json:"message"`
	Type    string  `json:"type"`
}

type FinalMessageReply struct {
	Message MessageReply `json:"message"`
	Type    string       `json:"type"`
}