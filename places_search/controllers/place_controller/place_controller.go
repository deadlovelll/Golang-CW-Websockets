package placecontroller

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"places_search/utils"
	"strconv"
)

// PlaceController handles database operations.
type PlaceController struct {
	Database *sql.DB // Adjust according to your actual Database connection struct
}

// GetPlaceByName returns places filtered by name.
func (gdb *PlaceController) GetPlaceByName(placeName string) ([]byte, error) {
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

	return gdb.getPlaces(query, placeName)
}

// GetPlaceWithHashtag returns places filtered by hashtag in the description.
func (gdb *PlaceController) GetPlaceWithHashtag(hashtag string) ([]byte, error) {
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
	return gdb.getPlaces(query, hashtag)
}

// getPlaces executes the provided query, processes each row, retrieves the presigned URL,
// and returns the result as JSON.
func (gdb *PlaceController) getPlaces(query string, arg interface{}) ([]byte, error) {
	rows, err := gdb.Database.Query(query, arg)
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

		if err := rows.Scan(&placeID, &placeName, &createdByID, &placeCommentCount, &placeLikesCount, &createdByUsername); err != nil {
			log.Printf("Error scanning row: %v", err)
			return nil, err
		}

		place := map[string]interface{}{
			"place_id":            placeID,
			"place_name":          placeName,
			"created_by_id":       createdByID,
			"created_by_username": createdByUsername,
			"place_comment_count": placeCommentCount,
			"place_likes_count":   placeLikesCount,
		}

		urlString := fmt.Sprintf("http://127.0.0.1:8170?place_id=%s", strconv.Itoa(placeID))
		if presignedResponse, err := utils.GetPresignedURL(urlString); err != nil {
			log.Printf("Error retrieving presigned URL for place %d: %v", placeID, err)
		} else {
			place["places_media"] = presignedResponse.PresignedURL
		}
		results = append(results, place)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating rows: %v", err)
		return nil, err
	}

	jsonData, err := json.Marshal(results)
	if err != nil {
		return nil, err
	}

	return jsonData, nil
}
