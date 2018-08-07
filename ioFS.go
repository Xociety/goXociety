package main

import (
	"bufio"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

func loadConfigFromFile(dataPath string, indexID, indexTarget, length int, includeReverse bool, forward map[int]string, reverse map[string]int) {
	file, err := os.Open(dataPath)
	defer file.Close()
	if err == nil {
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

func writeFile(filePath string, content string) {
	err := ioutil.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		log.Fatal(err)
	}
	return
}
