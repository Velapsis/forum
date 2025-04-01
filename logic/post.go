package logic

import (
	"forum/web/database"
	"hash/fnv"
	"math/rand/v2"
	"time"
)

func CreatePost(userID int, topic string, title string, post string) {
	database.AddPost(GeneratePostID(), title, post, topic, database.GetUsername(userID), time.Now(), time.Now())
}

func GeneratePostID() int {
	id := fnv.New32a()
	random := rand.IntN(9999999)
	return int(id.Sum32()) + random
}

func AddImage() {
	// TODO: ADD IMAGE TO UPLOADS AND SAVE PATH
}
