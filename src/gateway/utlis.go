package main

import (
	"os"
	"strings"
)

func Secret2Route(secret string) string {
	str := strings.Replace(secret, " ", "", -1)
	return str
}
func ParseUrl(url string) string {
	if url[len(url)-1:] == "/" {
		return url[:len(url)-1]
	} else {
		return url
	}
}
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}
func PathExist(_path string) bool {
	_, err := os.Stat(_path)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}
