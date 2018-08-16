package main

import (
	"bufio"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

func loadConfigFromFile(filePath string, indexID, indexTarget, length int, includeReverse bool, forward map[int]string, reverse map[string]int) {
	file, err := os.Open(filePath)
	if err == nil {
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			strSplitter := strings.Split(scanner.Text(), ";")
			if len(strSplitter) == length {
				if i, err := strconv.Atoi(strSplitter[indexID]); err == nil {
					forward[i] = strSplitter[indexTarget]
					if includeReverse {
						reverse[strSplitter[indexTarget]] = i
					}
				}
			}
		}
	}
}

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
func ioReaderFromFile(filePath string) (file io.Reader, err error) {
	file, err = os.Open(filePath)
	if err != nil {
		return nil, err
	}
	return bufio.NewReader(file), err
}
