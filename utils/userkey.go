package utils

import "regexp"

func IsUserKey (memo string) bool {
	if ok, _ := regexp.Match("^d[1-5a-zA-Z_]{33}$", []byte(memo)); ok {
		return true
	}
	return false
}
