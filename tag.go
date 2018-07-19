package main

import (
	"strconv"
	"strings"
	"unicode"
)

const punctuationHashtag = "#"

// there's no method in package unicode to escape this, even this not cover all test case
const charNotValidHashtag = "`~$^+|><=" + "¦⹋±∓№×°⋯⧸⁄÷−"

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
		if isAddHashtagStr {
			checkValid := false
			if unicode.IsLetter(rContent[i]) || unicode.IsNumber(rContent[i]) { // for speedup check
				checkValid = true
			} else if !strings.ContainsAny(charNotValidHashtag, char) {
				// you can't use unicode.IsSymbol because these [😘, 🉑] are exception
				checkValid = true
			}
			if checkValid {
				hashtag += char
			}
		}
	}
	if isAddHashtagStr && len(hashtag) > 0 { // collection rest info
		hashtags = append(hashtags, strings.ToLower(hashtag))
	}
	return
}

func parsehashtagOnPostSQL(postID int64, hashtagsID []int64) (sqlStrInsert, sqlStrDelete string, args []interface{}) {
	/*
		in order to insert multiple hashtag on post in one sql command, this func parse the command and parameters
		basic insert:
		sqlStr := `
			INSERT INTO post_hashtag (
				hashtag_id, post_id
			)
			VALUES($1,$2);
		`
	*/
	sqlStrInsert = `
		INSERT INTO post_hashtag (
			hashtag_id, post_id
		) 
		VALUES 
	`
	sqlStrDelete = `
		DELETE FROM post_hashtag WHERE post_id = $1 AND hashtag_id NOT IN (
	`
	args = append(args, postID)
	indexArg := 2
	for i := 0; i < len(hashtagsID); i++ {
		if i != 0 {
			sqlStrInsert += `, `
			sqlStrDelete += `, `
		}
		sqlStrInsert += `($` + strconv.Itoa(indexArg) + `,$1)`
		sqlStrDelete += `$` + strconv.Itoa(indexArg)
		args = append(args, hashtagsID[i])
		indexArg++
	}
	sqlStrInsert += `ON CONFLICT ON CONSTRAINT post_hashtag_hashtag_post_unique DO NOTHING;`
	sqlStrDelete += `);`
	return
}
