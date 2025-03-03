package placecontroller

import (
	"places_search/controllers/base_controller"
	"places_search/modules/database/database"
	"encoding/json"
	"fmt"
	"log"
	"places_search/utils"
	"strconv"
)

// PlaceController handles operations related to retrieving place information from the database.
//
// Fields:
//   - Database: A pointer to an sql.DB instance representing the database connection.
// PlaceController handles database operations and embeds BaseController for database access.
type PlaceController struct {
	*basecontroller.BaseController
}

// NewPlaceController creates a new instance of PlaceController.
func NewPlaceController(db *database.Database) *PlaceController {
	return &PlaceController{
		BaseController: &basecontroller.BaseController{Database: db},
	}
}

// GetPlaceByName retrieves places based on a partial name match.
//
// It searches for places whose names contain the given query string (case-insensitive) and excludes drafts.
//
// Parameters:
//   - placeName: A string representing part of the place name to search for.
//
// Returns:
//   - []byte: A JSON-encoded list of place details matching the search criteria.
//   - error: An error if the query execution or processing fails.
func (pc *PlaceController) GetPlaceByName(placeName string) ([]byte, error) {
	query := `
	SELECT 
		base_place.id AS place_id, 
		base_place.place_name, 
		base_place.created_by_id,
		COUNT(DISTINCT base_placecomment.id) AS place_comment_count,
		COUNT(DISTINCT base_place_place_likes.id) AS place_likes_count,
		COALESCE(base_user.username, 'None') AS user_username
	FROM base_place
		LEFT JOIN base_placecomment 
			ON base_placecomment.place_room_id = base_place.id 
		LEFT JOIN base_place_place_likes
			ON base_place_place_likes.place_id = base_place.id
		LEFT JOIN base_user
			ON base_place.created_by_id = base_user.id
		LEFT JOIN base_placephoto
			ON base_placephoto.parent_place_id = base_place.id
	WHERE base_place.place_name ILIKE '%' || $1 || '%' AND base_place.is_draft = false
	GROUP BY base_user.username, base_place.id, base_place.place_name, base_place.created_by_id;
	`

	return pc.getPlaces(query, placeName)
}

// GetPlaceWithHashtag retrieves places containing a specified hashtag in the description.
//
// The function searches for hashtags within the `description` column and excludes drafts.
//
// Parameters:
//   - hashtag: A string representing the hashtag to filter places.
//
// Returns:
//   - []byte: A JSON-encoded list of places matching the search criteria.
//   - error: An error if the query execution or processing fails.
func (pc *PlaceController) GetPlaceWithHashtag(hashtag string) ([]byte, error) {
	query := `
	SELECT 
		base_place.id AS place_id, 
		base_place.place_name, 
		base_place.created_by_id,
		COUNT(DISTINCT base_placecomment.id) AS place_comment_count,
		COUNT(DISTINCT base_place_place_likes.id) AS place_likes_count,
		COALESCE(base_user.username, 'None') AS user_username
	FROM base_place
		LEFT JOIN base_placecomment 
			ON base_placecomment.place_room_id = base_place.id 
		LEFT JOIN base_place_place_likes
			ON base_place_place_likes.place_id = base_place.id
		LEFT JOIN base_user
			ON base_place.created_by_id = base_user.id
		LEFT JOIN base_placephoto
			ON base_placephoto.parent_place_id = base_place.id
	WHERE base_place.description LIKE '%' || $1 || '%' AND base_place.is_draft = false
	GROUP BY base_user.username, base_place.id, base_place.place_name, base_place.created_by_id;
	`
	return pc.getPlaces(query, hashtag)
}

// getPlaces executes a query to fetch place details and appends presigned media URLs.
//
// The function runs a database query, iterates over results, retrieves media URLs, and
// returns a JSON-encoded list of places.
//
// Parameters:
//   - query: A string containing the SQL query to execute.
//   - arg: The argument value to use in the SQL query (e.g., place name or hashtag).
//
// Returns:
//   - []byte: A JSON-encoded list of places.
//   - error: An error if the query execution or JSON encoding fails.
func (pc *PlaceController) getPlaces(query string, arg interface{}) ([]byte, error) {
	rows, err := pc.Database.Query(query, arg)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return nil, err
	}
	defer rows.Close()

	var results []map[string]interface{}

	for rows.Next() {
		var (
			placeID           int
			placeName         string
			createdByID       int
			placeCommentCount int
			placeLikesCount   int
			createdByUsername string
		)

		// Scan the row into variables.
		if err := rows.Scan(&placeID, &placeName, &createdByID, &placeCommentCount, &placeLikesCount, &createdByUsername); err != nil {
			log.Printf("Error scanning row: %v", err)
			return nil, err
		}

		// Construct the place data.
		place := map[string]interface{}{
			"place_id":            placeID,
			"place_name":          placeName,
			"created_by_id":       createdByID,
			"created_by_username": createdByUsername,
			"place_comment_count": placeCommentCount,
			"place_likes_count":   placeLikesCount,
		}

		// Fetch the presigned URL for place media.
		urlString := fmt.Sprintf("http://127.0.0.1:8170?place_id=%s", strconv.Itoa(placeID))
		if presignedResponse, err := utils.GetPresignedURL(urlString); err != nil {
			log.Printf("Error retrieving presigned URL for place %d: %v", placeID, err)
		} else {
			place["places_media"] = presignedResponse.PresignedURL
		}

		results = append(results, place)
	}

	// Check for iteration errors.
	if err := rows.Err(); err != nil {
		log.Printf("Error iterating rows: %v", err)
		return nil, err
	}

	// Convert results to JSON.
	jsonData, err := json.Marshal(results)
	if err != nil {
		return nil, err
	}

	return jsonData, nil
}
