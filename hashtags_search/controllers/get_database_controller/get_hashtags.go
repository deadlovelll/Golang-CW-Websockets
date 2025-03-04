package getdatabasecontroller

import (
	"encoding/json"
	"fmt"
	basecontroller "hashtags_search/modules/base_controller"
	"net/http"
	"strings"
)

type GetDatabaseController struct {
	basecontroller.BaseController
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

func (gdb *GetDatabaseController) GetHashtags(Hashtag string) ([]byte, error) {

	fmt.Println(Hashtag)

	db := gdb.Database.GetConnection()

	cleanedHashtag := strings.TrimSpace(Hashtag)

	query := fmt.Sprintf(`
		WITH unique_hashtags_per_place AS (
			SELECT
				bp.id AS place_id,
				unnest(
					regexp_matches(
						bp.description, 
						'#%s\w*',
						'g'
					)
				) AS hashtag
			FROM base_place bp
			GROUP BY bp.id, hashtag
		)
		SELECT 
			hashtag, 
			COUNT(*) AS match_count
		FROM unique_hashtags_per_place
		GROUP BY hashtag
		ORDER BY match_count DESC
		`, cleanedHashtag)

	rows, err := db.Query(query)

	if err != nil {
		fmt.Printf("Error occured during exucting query: %s", err)
	}

	defer rows.Close()

	var results []map[string]interface{}

	for rows.Next() {

		var hashtag string
		var match_count int

		if err := rows.Scan(&hashtag, &match_count); err != nil {
			fmt.Printf("An error ocurred during scanning")
		}

		hashtag_dict := map[string]interface{}{
			"name":        hashtag,
			"match_count": match_count,
		}

		results = append(results, hashtag_dict)

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
