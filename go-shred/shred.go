package main

import (
	"crypto/rand"
	"io/ioutil"
	"os"
)

func Shred(path string) error {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return err
	}
	fileSize := fileInfo.Size()
	data := make([]byte, fileSize)

	for i := 0; i < 3; i++ {
		_, err = rand.Read(data)
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(path, data, fileInfo.Mode())
		if err != nil {
			return err
		}
	}

	err = os.Remove(path)
	return err
}