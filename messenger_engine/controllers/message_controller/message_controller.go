package messagecontroller

import (
	"context"
	"database/sql"
	"fmt"

	Messages "messenger_engine/models/message"
)

// MessageRepository defines operations for persisting messages.
type MessageRepository interface {
	SaveMessage(ctx context.Context, msg Messages.Message) error
	SaveMessageReply(ctx context.Context, msg Messages.MessageReply) error
	LoadMessages(ctx context.Context, chatID int) ([]Messages.Message, error)
}

// PostgresMessageRepository implements MessageRepository using PostgreSQL.
type PostgresMessageRepository struct {
	DB *sql.DB
}

// NewPostgresMessageRepository returns a new PostgresMessageRepository.
func NewPostgresMessageRepository(db *sql.DB) MessageRepository {
	return &PostgresMessageRepository{DB: db}
}

// SaveMessage inserts a new message into the database.
func (r *PostgresMessageRepository) SaveMessage(ctx context.Context, msg Messages.Message) error {
	query := `
		INSERT INTO base_chatmessage 
			(content, timestamp, author_id, chat_id, receiver_id, is_edited, parent)
		VALUES ($1, $2, $3, $4, $5, false, null)
	`
	_, err := r.DB.ExecContext(ctx, query, msg.Message, msg.Timestamp, msg.AuthorId, msg.ChatId, msg.ReceiverId)
	return err
}

// SaveMessageReply inserts a new reply message into the database.
func (r *PostgresMessageRepository) SaveMessageReply(ctx context.Context, msg Messages.MessageReply) error {
	query := `
		INSERT INTO base_chatmessage 
			(content, timestamp, author_id, chat_id, receiver_id, parent_id, is_edited)
		VALUES ($1, $2, $3, $4, $5, $6, false)
	`
	_, err := r.DB.ExecContext(ctx, query, msg.Message, msg.Timestamp, msg.AuthorId, msg.ChatId, msg.ReceiverId, msg.ParentMessageId)
	return err
}

// LoadMessages retrieves messages for a specific chat.
func (r *PostgresMessageRepository) LoadMessages(ctx context.Context, chatID int) ([]Messages.Message, error) {
	query := `
		SELECT message_id, content, is_edited, timestamp, author_id, chat_id, receiver_id, parent
		FROM base_chatmessage
		WHERE chat_id = $1
	`
	rows, err := r.DB.QueryContext(ctx, query, chatID)
	if err != nil {
		return nil, fmt.Errorf("querying messages: %w", err)
	}
	defer rows.Close()

	var messages []Messages.Message
	for rows.Next() {
		var msg Messages.Message
		if err := rows.Scan(
			&msg.MessageId,
			&msg.Message,
			&msg.IsEdited,
			&msg.Timestamp,
			&msg.AuthorId,
			&msg.ChatId,
			&msg.ReceiverId,
			&msg.ParentMessageId,
		); err != nil {
			return nil, fmt.Errorf("scanning message: %w", err)
		}
		messages = append(messages, msg)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}
	return messages, nil
}
