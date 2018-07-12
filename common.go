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
