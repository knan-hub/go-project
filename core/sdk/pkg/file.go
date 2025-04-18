package pkg

import (
	"log"
	"os"
)

func PathExist(s string) bool {
	fileInfo, err := os.Stat(s)
	if err != nil {
		log.Println(err)
		return false
	}
	return fileInfo.IsDir()
}

func PathCreate(s string) error {
	return os.MkdirAll(s, os.ModePerm)
}
