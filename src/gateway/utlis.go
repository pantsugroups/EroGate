package gateway

import (
	"log"
	"os"
)

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
func CheckConfFolder() {
	if !PathExist("conf.d/") && !IsDir("conf.d/") {
		err := os.Mkdir("conf.d/", os.ModePerm)
		if err != nil {
			log.Println(err)
		}
	}
}
