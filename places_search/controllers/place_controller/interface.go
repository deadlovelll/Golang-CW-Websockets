package placecontroller

// PlaceControllerInterface defines the methods required by the WebSocketHandler.
type PlaceControllerInterface interface {
	// GetPlaceByName searches places by name and returns a JSON response.
	GetPlaceByName(query string) ([]byte, error)
	// GetPlaceWithHashtag searches places by hashtag and returns a JSON response.
	GetPlaceWithHashtag(query string) ([]byte, error)
}
