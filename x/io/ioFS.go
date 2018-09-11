package io

import (
	"bufio"
	"io"
	"io/ioutil"
	"log"
	"os"
)

func loadListFromFile(filePath string) (list []string) {
	file, err := os.Open(filePath)
	if err == nil {
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			list = append(list, scanner.Text())
		}
	}
	return list
}

func writeFile(filePath, content string) {
	err := ioutil.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		log.Fatal(err)
	}
	return
}
func readFile(filePath string) (data []byte, err error) {
	return ioutil.ReadFile(filePath)
}
func IoReaderFromFile(filePath string) (file io.Reader, err error) {
	file, err = os.Open(filePath)
	if err != nil {
		return nil, err
	}
	return bufio.NewReader(file), err
}
