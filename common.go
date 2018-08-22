package main

import (
	"crypto/sha256"
	"encoding/hex"
	"time"
)

func getNowUnixTimestamp() int {
	return int(time.Now().Unix())
}

func getNowUnixWeekTimestamp() int {
	return getNowUnixTimestamp() / sevenDaysInSecond
}

func getSHA256Hash(input []byte) string {
	h := sha256.New()
	h.Write(input)
	return hex.EncodeToString(h.Sum(nil))
}

func convertIDAndValue(IDs []int, values []string, isReverse bool, forward map[int]string, reverse map[string]int) {
	if len(IDs) == len(values) {
		for i := 0; i < len(IDs); i++ {
			forward[IDs[i]] = values[i]
			if isReverse {
				reverse[values[i]] = IDs[i]
			}
		}
	}
}
