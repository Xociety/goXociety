package common

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"reflect"
	"time"

	"github.com/chienfuchen32/goXociety/x/config"
)

func GetNowUnixTimestamp() int {
	return int(time.Now().Unix())
}

func GetNowUnixWeekTimestamp() int {
	return GetNowUnixTimestamp() / config.SevenDaysInSecond
}

func getSHA256Hash(input []byte) string {
	h := sha256.New()
	h.Write(input)
	return hex.EncodeToString(h.Sum(nil))
}

func ConvertIDAndValue(IDs []int, values []string, isReverse bool, forward map[int]string, reverse map[string]int) {
	if len(IDs) == len(values) {
		for i := 0; i < len(IDs); i++ {
			forward[IDs[i]] = values[i]
			if isReverse {
				reverse[values[i]] = IDs[i]
			}
		}
	}
}

func ConvertInterfaceSlice(slice interface{}) ([]interface{}, error) {
	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		return nil, errors.New("ConvertInterfaceSlice() given a non-slice type")
	}
	ret := make([]interface{}, s.Len())
	for i := 0; i < s.Len(); i++ {
		ret[i] = s.Index(i).Interface()
	}
	return ret, nil
}
