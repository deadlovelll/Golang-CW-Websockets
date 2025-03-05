package hashtagcontroller

import (
	"database/sql"
	"encoding/json"
	"fmt"
	basecontroller "hashtags_search/controllers/base_controller"
)

var GetHashtagsFromDBB = defaultGetHashtagsFromDB

// defaultGetHashtagsFromDB is the default implementation that interacts with the database.
func defaultGetHashtagsFromDB(db *sql.DB, hashtag string) ([]map[string]interface{}, error) {
    // Replace with your actual database query logic.
    return nil, fmt.Errorf("database query not implemented")
}


type HashtagProvider interface {
	GetHashtags(hashtag string) ([]byte, error)
}

// HashtagController embeds BaseController, inheriting its behavior
type HashtagController struct {
	basecontroller.BaseController // Embedding BaseController to "inherit" its methods
}

// NewHashtagController is a constructor for HashtagController
func NewHashtagController() *HashtagController {
	return &HashtagController{
		BaseController: basecontroller.BaseController{}, // Initialize the embedded struct
	}
}

// GetHashtags queries the database for hashtags matching the given keyword and returns them as JSON.
//
// Parameters:
//   - hashtag: The hashtag keyword to search for in the database.
//
// Returns:
//   - JSON-encoded byte slice containing the hashtag results.
//   - An error if the query or JSON marshaling fails.
func (hc *HashtagController) GetHashtags(hashtag string) ([]byte, error) {
	fmt.Println("Searching for hashtag:", hashtag)

	db := hc.Database.GetConnection()
	results, err := GetHashtagsFromDB(db, hashtag)
	if err != nil {
		return nil, err
	}

	jsonData, err := json.Marshal(results)
	if err != nil {
		return nil, fmt.Errorf("error marshaling JSON: %v", err)
	}

	return jsonData, nil
}