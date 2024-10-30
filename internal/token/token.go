package token

import "strings"

// GetToken extracts the token from the Authorization header
func GetToken(authentication string) string {
	elements := strings.Split(authentication, " ")
	if len(elements) != 2 {
		return ""
	}
	return elements[1]
}
