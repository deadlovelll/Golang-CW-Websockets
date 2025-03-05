package hashtagcontroller

import (
	"database/sql"
	"fmt"
	"strings"
)

// GetHashtagsFromDB queries the database to find hashtags matching the given keyword.
//
// Parameters:
//   - db: The database connection.
//   - hashtag: The hashtag keyword to search for.
//
// Returns:
//   - A slice of maps containing hashtag names and their match counts.
//   - An error if the query execution or row scanning fails.
func GetHashtagsFromDB(db *sql.DB, hashtag string) ([]map[string]interface{}, error) {
	cleanedHashtag := strings.TrimSpace(hashtag)

	query := fmt.Sprintf(`
		WITH unique_hashtags_per_place AS (
			SELECT
				bp.id AS place_id,
				unnest(
					regexp_matches(bp.description, '#%s\\w*', 'g')
				) AS hashtag
			FROM base_place bp
			GROUP BY bp.id, hashtag
		)
		SELECT hashtag, COUNT(*) AS match_count
		FROM unique_hashtags_per_place
		GROUP BY hashtag
		ORDER BY match_count DESC
	`, cleanedHashtag)

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %v", err)
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		var hashtag string
		var matchCount int
		if err := rows.Scan(&hashtag, &matchCount); err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}

		results = append(results, map[string]interface{}{
			"name":        hashtag,
			"match_count": matchCount,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}

	return results, nil
}
