package postmessengercontroller

import (
	"context"

	BaseController "messenger_engine/controllers/base_controller"
	Messages "messenger_engine/models/message"
	MessageControl "messenger_engine/controllers/message_controller"
)

// MessageController handles message operations using a repository.
type MessageController struct {
	BaseController *BaseController.BaseController
	Repo           MessageControl.MessageRepository
}

// SaveMessage saves a new message using the repository.
func (mc *MessageController) SaveMessage(ctx context.Context, msg Messages.Message) error {
	return mc.Repo.SaveMessage(ctx, msg)
}

// SaveMessageReply saves a reply message using the repository.
func (mc *MessageController) SaveMessageReply(ctx context.Context, msg Messages.MessageReply) error {
	return mc.Repo.SaveMessageReply(ctx, msg)
}

// LoadMessages retrieves all messages for a specific chat.
func (mc *MessageController) LoadMessages(ctx context.Context, chatID int) ([]Messages.Message, error) {
	return mc.Repo.LoadMessages(ctx, chatID)
}