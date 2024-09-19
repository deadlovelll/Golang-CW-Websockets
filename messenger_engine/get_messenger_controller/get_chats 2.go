package getmessengercontroller

import (
	"encoding/json"
	"fmt"
	basecontroller "messenger_engine/modules/base_controller"
	"net/http"
	"strconv"
)

type GetMessengerController struct {
	*basecontroller.BaseController
}

type MyResponse struct {
	Status       string `json:"STATUS"`
	PresignedURL string `json:"PRESIGNED_URL"`
}

func GetPresignedUrl(url string, responseCh chan<- *MyResponse, errorCh chan<- string) {
	resp, err := http.Get(url)
	if err != nil {
		errorCh <- fmt.Sprintf("Error making request: %v", err)
		return
	}
	defer resp.Body.Close()

	var response MyResponse
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&response)
	if err != nil {
		errorCh <- fmt.Sprintf("Error decoding JSON: %v", err)
		return
	}

	responseCh <- &response
}

func (gmc *GetMessengerController) GetUserChats(UserID int) ([]byte, error) {

	db := gmc.Database.GetConnection()

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

	rows, err := db.Query(query, UserID)
	if err != nil {
		fmt.Printf("Error executing query: %v", err)
		return nil, err
	}
	defer rows.Close()

	var results []map[string]interface{}

	for rows.Next() {
		var (
			content          string
			timestamp        string
			chatID           int
			receiverID       int
			receiverUsername string
		)

		// Scan into variables
		err := rows.Scan(&content, &timestamp, &chatID, &receiverID, &receiverUsername)
		if err != nil {
			fmt.Printf("Error scanning row: %v", err)
			return nil, err
		}

		// Create a map to store the result
		chatsDict := map[string]interface{}{
			"messge_content":            content,
			"message_timestamp":         timestamp,
			"chat_id":                   chatID,
			"message_receiver_id":       receiverID,
			"message_receiver_username": receiverUsername,
		}

		StrUserId := strconv.Itoa(UserID)

		urlString := fmt.Sprintf("http://127.0.0.1:8165?user_id=%s", StrUserId)

		responseCh := make(chan *MyResponse)
		errorCh := make(chan string)

		go GetPresignedUrl(urlString, responseCh, errorCh)

		select {
		case response := <-responseCh:
			fmt.Println(response)
			chatsDict["user_avatar_url"] = response.PresignedURL
		case err := <-errorCh:
			fmt.Printf("Received error: %s\n", err)
		}

		// Append the result to the list
		results = append(results, chatsDict)
	}

	// Check for errors from iterating through rows
	if err := rows.Err(); err != nil {
		fmt.Printf("Error iterating rows: %v", err)
		return nil, err
	}

	// Marshal results into JSON
	jsonData, err := json.Marshal(results)
	if err != nil {
		return nil, err
	}

	return jsonData, nil
}
