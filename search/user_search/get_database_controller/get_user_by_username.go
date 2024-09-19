package getdatabasecontroller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func GetPresignedUrlUsername(url string, responseCh chan<- *MyResponse, errorCh chan<- string) {
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

func (gdb *GetDatabaseController) GetUsersByUsername(Username string) ([]byte, error) {
	db := gdb.Database.GetConnection()

	query := `
	WITH UserAFriends AS (
		SELECT bff.user_id
		FROM base_friendslist bfl
		INNER JOIN base_friendslist_friends bff
			ON bff.friendslist_id = bfl.id
		WHERE bfl.user_id IN (
			SELECT id FROM base_user WHERE username LIKE $1 || '%'
		)
	),
	UserBFriends AS (
		SELECT bff.user_id
		FROM base_friendslist bfl
		INNER JOIN base_friendslist_friends bff
			ON bff.friendslist_id = bfl.id
		WHERE bfl.user_id = 2
	),
	CommonFriends AS (
		SELECT u.username
		FROM base_user u
		INNER JOIN UserAFriends af ON u.id = af.user_id
		INNER JOIN UserBFriends bf ON u.id = bf.user_id
	),
	FriendCount AS (
		SELECT COUNT(*) AS count
		FROM CommonFriends
	)
	SELECT 
	    base_user.id,
	    base_user.username,
	    base_user.verified_account,
	    COALESCE(CONCAT(base_user.first_name, ' ', base_user.last_name), base_user.first_name) AS full_name,
	    EXISTS (
	        SELECT 1
	        FROM base_userstory
	        WHERE base_userstory.posted_by_id = base_user.id
	        AND base_userstory.is_visible = true
	    ) AS has_active_stories,
	    COALESCE(base_useravatar.image_path, 'default_avatar_image_path') AS avatar_image,
	    CASE
	        WHEN fc.count = 0 THEN ''
	        WHEN fc.count = 1 THEN
	            (SELECT username FROM CommonFriends ORDER BY username LIMIT 1)
	        ELSE
	            CONCAT(
	                (SELECT username FROM CommonFriends ORDER BY username LIMIT 1),
	                ' + ',
	                fc.count,
	                ' Mutuals'
	            )
	    END AS common_friends
	FROM base_user
	LEFT JOIN base_useravatar
	    ON base_user.id = base_useravatar.avatar_user_id
	    AND base_useravatar.is_current = true
	LEFT JOIN CommonFriends
	    ON base_user.username = CommonFriends.username
	LEFT JOIN FriendCount fc ON true
	WHERE base_user.username LIKE $1 || '%'
	GROUP BY base_user.id, base_user.username, base_user.verified_account, base_user.first_name, base_user.last_name, base_useravatar.avatar_user_id, base_useravatar.is_current, base_useravatar.image_path, fc.count;
	`

	rows, err := db.Query(query, Username)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return nil, err
	}
	defer rows.Close()

	var results []map[string]interface{}

	for rows.Next() {
		var (
			username         string
			userid           int
			verifiedAccount  bool
			fullName         string
			userAvatar       string
			hasActiveStories bool
			commonFriends    string
		)

		err := rows.Scan(&userid, &username, &verifiedAccount, &fullName, &hasActiveStories, &userAvatar, &commonFriends)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			return nil, err
		}

		userDict := map[string]interface{}{
			"userid":           userid,
			"username":         username,
			"verifiedAccount":  verifiedAccount,
			"fullName":         fullName,
			"hasActiveStories": hasActiveStories,
			"commonFriends":    commonFriends,
		}

		StrUserId := strconv.Itoa(userid)

		urlString := fmt.Sprintf("http://127.0.0.1:8165?user_id=%s", StrUserId)

		responseCh := make(chan *MyResponse)
		errorCh := make(chan string)

		go GetPresignedUrl(urlString, responseCh, errorCh)

		select {
		case response := <-responseCh:
			userDict["userAvatar"] = response.PresignedURL
		case err := <-errorCh:
			fmt.Printf("Received error: %s\n", err)
		}

		results = append(results, userDict)
	}

	if len(results) == 0 {
		// Return nil if no results are found
		return nil, nil
	}

	jsonData, err := json.Marshal(results)
	if err != nil {
		return nil, err
	}

	return jsonData, nil
}
