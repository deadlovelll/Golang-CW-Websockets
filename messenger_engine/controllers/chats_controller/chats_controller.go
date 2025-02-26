package chatscontroller

import (
	"fmt"

	BaseController "messenger_engine/controllers/base_controller"
)

// GetMessengerController handles fetching user messages
type GetMessengerController struct {
	*BaseController.BaseController
	Repo               MessengerRepository
	PresignedURLClient *PresignedURLService
}

// GetUserChats retrieves the latest user chats
func (gmc *GetMessengerController) GetUserChats(ctx context.Context, userID int) ([]byte, error) {
	messages, err := gmc.Repo.GetUserChats(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}

	var results []map[string]interface{}

	for _, msg := range messages {
		chatDict := map[string]interface{}{
			"message_content":            msg.Content,
			"message_timestamp":         msg.Timestamp,
			"chat_id":                   msg.ChatID,
			"message_receiver_id":       msg.ReceiverID,
			"message_receiver_username": msg.ReceiverUsername,
		}

		// Fetch user avatar URL with timeout
		ctxTimeout, cancel := context.WithTimeout(ctx, 3*time.Second)
		defer cancel()

		if avatarURL, err := gmc.PresignedURLClient.GetPresignedURL(ctxTimeout, userID); err == nil {
			chatDict["user_avatar_url"] = avatarURL
		} else {
			fmt.Printf("Failed to fetch avatar URL: %v\n", err)
		}

		results = append(results, chatDict)
	}

	return json.Marshal(results)
}