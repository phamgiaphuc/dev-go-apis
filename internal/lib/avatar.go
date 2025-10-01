package lib

import (
	"fmt"
	"math/rand"
)

var AVATARS = []string{
	"/static/avatars/ava1.png",
	"/static/avatars/ava2.png",
	"/static/avatars/ava3.png",
	"/static/avatars/ava4.png",
	"/static/avatars/ava5.png",
}

func GetRandomAvatar() string {
	result := fmt.Sprintf("%s%s", SERVER_URL, AVATARS[rand.Intn(len(AVATARS))])
	return result
}

func GetExternalAvatar(username string) string {
	result := fmt.Sprintf("https://avatar.iran.liara.run/username?username=%s", username)
	return result
}
