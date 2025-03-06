package usercontroller

const mutualFriendsQuery = `
	WITH mutual_friends AS (
		...
	)
	SELECT
		...
	FROM base_user bu
	...
	WHERE bu.id::text LIKE $1 || '%';`

const usersByUsernameQuery = `
	WITH UserAFriends AS (
		...
	)
	SELECT
		...
	FROM base_user
	...
	WHERE base_user.username LIKE $1 || '%'
	GROUP BY
		...`
