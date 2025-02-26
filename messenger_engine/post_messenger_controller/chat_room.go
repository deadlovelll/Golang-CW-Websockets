package postmessengercontroller

import (
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
func (mc *MessageController) SaveMessage(msg Messages.Message) error {
	return mc.Repo.SaveMessage(msg)
}

// SaveMessageReply saves a reply message using the repository.
func (mc *MessageController) SaveMessageReply(msg Messages.MessageReply) error {
	return mc.Repo.SaveMessageReply(msg)
}

// LoadMessages retrieves all messages for a specific chat.
func (mc *MessageController) LoadMessages(chatID int) ([]Messages.Message, error) {
	return mc.Repo.LoadMessages(chatID)
}