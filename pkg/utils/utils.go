package utils

import "strings"

func GetIdFromParams(str string) string {
	// get id param. 0th param is empty space, 1st is joke, 2nd is id
	id := strings.Split(str, "/")[2]
	return id
}
