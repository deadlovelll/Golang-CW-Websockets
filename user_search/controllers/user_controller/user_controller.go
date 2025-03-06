package usercontroller

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"user_search/controllers/base_controller"
	"user_search/controllers/url_controller"
)

type UserController struct {
	baseCtrl *basecontroller.BaseController
}

type UserControllerInterface interface {
	GetUsers(userID int) ([]byte, error)
	GetUsersByUsername(username string) ([]byte, error)
}

func NewUserController(baseCtrl *basecontroller.BaseController) *UserController {
	return &UserController{baseCtrl: baseCtrl}
}

func (uc *UserController) GetUsers(userId int) ([]byte, error) {
	db := uc.baseCtrl.Database.GetConnection()

	query := getMutualFriendsQuery()

	rows, err := db.Query(query, userId)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %v", err)
	}
	defer rows.Close()

	results, err := scanUsers(rows)
	if err != nil {
		return nil, err
	}

	return json.Marshal(results)
}

func (uc *UserController) GetUsersByUsername(username string) ([]byte, error) {
	db := uc.baseCtrl.Database.GetConnection()

	query := getUsersByUsernameQuery()

	rows, err := db.Query(query, username)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %v", err)
	}
	defer rows.Close()

	results, err := scanUsers(rows)
	if err != nil {
		return nil, err
	}

	return json.Marshal(results)
}

// Helper to scan user rows into a map
func scanUsers(rows *sql.Rows) ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	for rows.Next() {
		var (
			userid           int
			username         string
			verifiedAccount  bool
			fullName         string
			userAvatar       string
			hasActiveStories bool
			commonFriends    string
		)

		if err := rows.Scan(&userid, &username, &verifiedAccount, &fullName, &hasActiveStories, &userAvatar, &commonFriends); err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}

		userDict := map[string]interface{}{
			"userid":           userid,
			"username":         username,
			"verifiedAccount":  verifiedAccount,
			"fullName":         fullName,
			"hasActiveStories": hasActiveStories,
			"commonFriends":    commonFriends,
		}

		userDict["userAvatar"], _ = fetchAvatarURL(userid)
		results = append(results, userDict)
	}
	return results, nil
}

// Helper to fetch avatar URL concurrently
func fetchAvatarURL(userID int) (string, error) {
	url := fmt.Sprintf("http://127.0.0.1:8165?user_id=%d", userID)
	responseCh := make(chan *urlcontroller.UrlResponse)
	errorCh := make(chan string)

	fetcher := &urlcontroller.HttpPresignedUrlFetcher{}
	go fetcher.Fetch(url)

	select {
	case response := <-responseCh:
		return response.PresignedURL, nil
	case err := <-errorCh:
		return "", fmt.Errorf("error fetching avatar: %s", err)
	}
}

// Predefined SQL query strings
func getMutualFriendsQuery() string {
	return mutualFriendsQuery
}

func getUsersByUsernameQuery() string {
	return usersByUsernameQuery
}