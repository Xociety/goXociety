package main

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

func loadConfigFromFile(dataPath string, indexID, indexTarget, length int, includeReverse bool, forward map[int]string, reverse map[string]int) {
	if file, err := os.Open(dataPath); err == nil {
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
