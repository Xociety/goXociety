package main

import (
	"strings"
	"unicode"
)

const punctuationHashtag = "#"

// there's no method in package unicode to escape this, even this not cover all test case
const charNotValidHashtag = "`~$^+|><=¦⹋±∓№×°⋯ฯ⧸⁄÷ºª−"

func checkMention(content string) (hashtags []string, tags []string) {
	rContent := []rune(content)
	hashtag := ""
	isAddHashtagStr := false
	for i := 0; i < len(rContent); i++ {
		char := string(rContent[i])
		if char != "_" {
			if unicode.IsPunct(rContent[i]) ||
				unicode.IsSpace(rContent[i]) {
				if isAddHashtagStr && len(hashtag) > 0 {
					hashtags = append(hashtags, strings.ToLower(hashtag))
				}
				hashtag = ""
				switch char {
				case punctuationHashtag:
					isAddHashtagStr = true
					continue // for loop
				default:
					isAddHashtagStr = false
					hashtag = ""
				}
			}
		}
		if !strings.ContainsAny(charNotValidHashtag, char) {
			if isAddHashtagStr {
				hashtag += char
			}
		}
	}
	if isAddHashtagStr && len(hashtag) > 0 { // collection rest info
		hashtags = append(hashtags, strings.ToLower(hashtag))
	}
	return
}

func checkMentionBK(content string) (hashtags []string, tags []string) {
	rContent := []rune(content)
	hashtag := ""
	isAddHashtagStr := false
	for i := 0; i < len(rContent); i++ {
		char := string(rContent[i])
		if unicode.IsPunct(rContent[i]) || unicode.IsSpace(rContent[i]) {
			switch char {
			case "_":
			case punctuationHashtag:
				if isAddHashtagStr && len(hashtag) > 0 {
					hashtags = append(hashtags, hashtag)
				}
				hashtag = ""
				isAddHashtagStr = true
				continue // for loop
			default:
				if isAddHashtagStr && len(hashtag) > 0 {
					hashtags = append(hashtags, hashtag)
				}
				isAddHashtagStr = false
				hashtag = ""
			}
		}
		if isAddHashtagStr {
			hashtag += char
		}
	}
	if isAddHashtagStr && len(hashtag) > 0 {
		hashtags = append(hashtags, hashtag)
		isAddHashtagStr = false
	}
	return
}
