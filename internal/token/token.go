package token

import (
	"log"
	"strings"
)

// GetToken extracts the token from the Authorization header
func GetToken(authentication string) string {
	log.Printf("authentication = %+v\n", authentication)
	elements := strings.Split(authentication, " ")
	if len(elements) != 2 {
		return ""
	}
	return elements[1]
}
