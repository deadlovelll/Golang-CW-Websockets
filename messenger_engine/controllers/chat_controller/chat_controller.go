package chatcontroller

import (
	"encoding/json"
	"fmt"
	"strconv"

	BaseController "messenger_engine/controllers/base_controller"
	UlrResponse "messenger_engine/models/presigned_url"
	UrlController "messenger_engine/controllers/url_controller"
)

// ChatController manages chat-related operations, including fetching user chats.
type ChatController struct {
	*BaseController.BaseController      // Embeds the base controller for shared functionality
	UrlFetcher UrlController.HttpPresignedUrlFetcher // Fetcher for retrieving presigned URLs for avatars
}

// GetUserChats retrieves a user's chats from the database by their user ID.
// It fetches the most recent chat message for each chat, including the content, timestamp, receiver info, and avatar URL.
func (gmc *ChatController) GetUserChats(UserID int) ([]byte, error) {
	// Get the database connection
	db := gmc.Database.GetConnection()

	// SQL query to fetch the most recent chat message for each chat
	query := `
		SELECT 
			bcm.content, 
			bcm.timestamp, 
			bcm.chat_id, 
			bcm.receiver_id, 
			base_user.username AS receiver_username 

		FROM base_chatmessage AS bcm
		
		LEFT JOIN (
			SELECT id, username
			FROM base_user
		) base_user ON bcm.receiver_id = base_user.id
		
		WHERE bcm.author_id = $1

		GROUP BY bcm.content, bcm.timestamp, bcm.chat_id, bcm.receiver_id, username
		ORDER BY timestamp DESC

		LIMIT 1
	`

	// Execute the query
	rows, err := db.Query(query, UserID)
	if err != nil {
		// Log the error and return it
		fmt.Printf("Error executing query: %v", err)
		return nil, err
	}
	defer rows.Close() // Ensure the rows are closed after processing

	var results []map[string]interface{} // Holds the results to be returned

	// Iterate through the query results
	for rows.Next() {
		var (
			content          string
			timestamp        string
			chatID           int
			receiverID       int
			receiverUsername string
		)

		// Scan the row into the variables
		err := rows.Scan(&content, &timestamp, &chatID, &receiverID, &receiverUsername)
		if err != nil {
			// Log scanning error and return it
			fmt.Printf("Error scanning row: %v", err)
			return nil, err
		}

		// Create a map to store the chat data
		chatsDict := map[string]interface{}{
			"messge_content":            content,
			"message_timestamp":         timestamp,
			"chat_id":                   chatID,
			"message_receiver_id":       receiverID,
			"message_receiver_username": receiverUsername,
		}

		// Convert the user ID to string for the URL
		StrUserId := strconv.Itoa(UserID)

		// Prepare the URL for fetching the avatar
		urlString := fmt.Sprintf("http://127.0.0.1:8165?user_id=%s", StrUserId)

		// Channels for receiving the presigned URL or error response
		responseCh := make(chan *UlrResponse.PresignedUrlResponse)
		errorCh := make(chan string)

		// Asynchronously fetch the presigned URL for the user's avatar
		go gmc.UrlFetcher.Fetch(urlString)

		// Handle the response or error from the URL fetch
		select {
		case response := <-responseCh:
			// Add the fetched avatar URL to the chat data
			chatsDict["user_avatar_url"] = response.PresignedURL
		case err := <-errorCh:
			// Log any error that occurs while fetching the avatar URL
			fmt.Printf("Received error: %s\n", err)
		}

		// Append the chat data to the results list
		results = append(results, chatsDict)
	}

	// Check for errors during row iteration
	if err := rows.Err(); err != nil {
		// Log any error and return it
		fmt.Printf("Error iterating rows: %v", err)
		return nil, err
	}

	// Marshal the results into JSON format for the response
	jsonData, err := json.Marshal(results)
	if err != nil {
		// Handle any error during JSON marshaling
		return nil, err
	}

	return jsonData, nil
}
