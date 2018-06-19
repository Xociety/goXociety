package main

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

func loadPostTypeFromFConfig(dataPath string) {
	if file, err := os.Open(dataPath); err == nil {
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			strSplitter := strings.Split(scanner.Text(), ";")
			if len(strSplitter) > 1 {
				if i, err := strconv.Atoi(strSplitter[0]); err == nil {
					postTypeMapID2Type[i] = strSplitter[1]
					postTypeMapType2ID[strSplitter[1]] = i
				}
			}
		}
	}
}

func loadActionsTypeFromFConfig(dataPath string) {
	if file, err := os.Open(dataPath); err == nil {
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			strSplitter := strings.Split(scanner.Text(), ";")
			if len(strSplitter) > 1 {
				if i, err := strconv.Atoi(strSplitter[0]); err == nil {
					actionsTypeMapID2Type[i] = strSplitter[1]
					actionsTypeMapType2ID[strSplitter[1]] = i
				}
			}
		}
	}
}
