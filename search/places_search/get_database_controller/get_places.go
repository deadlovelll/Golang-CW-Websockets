package getdatabasecontroller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	basecontroller "places_search/modules/base_controller"
)

type GetDatabaseController struct {
	*basecontroller.BaseController
}

type MyResponse struct {
	Status       string   `json:"STATUS"`
	PresignedURL []string `json:"PRESIGNED_URL"`
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

func (gdb *GetDatabaseController) GetPlaceByName(PlaceName string) ([]byte, error) {

	db := gdb.Database.GetConnection()

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

	rows, err := db.Query(query, PlaceName)
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
			createdByUsername string
			placeCommentCount int
			placeLikesCount   int
		)

		// Scan into variables
		err := rows.Scan(&placeID, &placeName, &createdByID, &placeCommentCount, &placeLikesCount, &createdByUsername)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			return nil, err
		}

		// Create a map to store the result
		placesDict := map[string]interface{}{
			"place_id":            placeID,
			"place_name":          placeName,
			"created_by_id":       createdByID,
			"created_by_username": createdByUsername,
			"place_comment_count": placeCommentCount,
			"place_likes_count":   placeLikesCount,
		}

		StrPlaceId := strconv.Itoa(placeID)

		urlString := fmt.Sprintf("http://127.0.0.1:8170?place_id=%s", StrPlaceId)

		responseCh := make(chan *MyResponse)
		errorCh := make(chan string)

		go GetPresignedUrl(urlString, responseCh, errorCh)

		select {
		case response := <-responseCh:
			placesDict["places_media"] = response.PresignedURL
		case err := <-errorCh:
			fmt.Printf("Received error: %s\n", err)
		}

		// Append the result to the list
		results = append(results, placesDict)
	}

	// Check for errors from iterating through rows
	if err := rows.Err(); err != nil {
		log.Printf("Error iterating rows: %v", err)
		return nil, err
	}

	// Marshal results into JSON
	jsonData, err := json.Marshal(results)
	if err != nil {
		return nil, err
	}

	return jsonData, nil
}

func (gdb *GetDatabaseController) GetPlaceWithHashtag(Hashtag string) ([]byte, error) {

	db := gdb.Database.GetConnection()

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

	rows, err := db.Query(query, Hashtag)
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
			createdByUsername string
			placeCommentCount int
			placeLikesCount   int
		)

		// Scan into variables
		err := rows.Scan(&placeID, &placeName, &createdByID, &placeCommentCount, &placeLikesCount, &createdByUsername)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			return nil, err
		}

		// Create a map to store the result
		placesDict := map[string]interface{}{
			"place_id":            placeID,
			"place_name":          placeName,
			"created_by_id":       createdByID,
			"created_by_username": createdByUsername,
			"place_comment_count": placeCommentCount,
			"place_likes_count":   placeLikesCount,
		}

		StrPlaceId := strconv.Itoa(placeID)

		urlString := fmt.Sprintf("http://127.0.0.1:8170?place_id=%s", StrPlaceId)

		responseCh := make(chan *MyResponse)
		errorCh := make(chan string)

		go GetPresignedUrl(urlString, responseCh, errorCh)

		select {
		case response := <-responseCh:
			placesDict["places_media"] = response.PresignedURL
		case err := <-errorCh:
			fmt.Printf("Received error: %s\n", err)
		}

		// Append the result to the list
		results = append(results, placesDict)
	}

	// Check for errors from iterating through rows
	if err := rows.Err(); err != nil {
		log.Printf("Error iterating rows: %v", err)
		return nil, err
	}

	// Marshal results into JSON
	jsonData, err := json.Marshal(results)
	if err != nil {
		return nil, err
	}

	return jsonData, nil
}
