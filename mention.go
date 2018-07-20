package main

import (
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

const punctuationHashtag = "#"
const punctuationTag = "@"

// there's no method in package unicode to escape this, even this not cover all test case
const charNotValidHashtag = "`~$^+|><=" + "¦⹋±∓№×°⋯⧸⁄÷−"
const replNotValidTag = "[^a-z0-9_]+"

func checkMention(content string) (hashtags []string, tags []string) {
	var regTag, _ = regexp.Compile(replNotValidTag)
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
	tag := ""
	isAddTagStr := false
	for i := 0; i < len(rContent); i++ {
		char := string(rContent[i])
		if char != "_" {
			if unicode.IsPunct(rContent[i]) ||
				unicode.IsSpace(rContent[i]) {
				if isAddTagStr && len(regTag.ReplaceAllString(tag, "")) > 0 {
					tags = append(tags, strings.ToLower(regTag.ReplaceAllString(tag, "")))
				}
				tag = ""
				switch char {
				case punctuationTag:
					isAddTagStr = true
					continue // for loop
				default:
					isAddTagStr = false
					tag = ""
				}
			}
		}
		if isAddTagStr {
			tag += char
		}
	}
	if isAddTagStr && len(regTag.ReplaceAllString(tag, "")) > 0 { // collection rest info
		tags = append(tags, strings.ToLower(regTag.ReplaceAllString(tag, "")))
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
func parseTagOnPostInserSQL(postID int64, tags []tagOnPostSetAPI) (sqlStr string, args []interface{}) {
	/*
		in order to insert multiple hashtag on post in one sql command, this func parse the command and parameters
		basic insert:
		sqlStr := `
			INSERT INTO post_tag_xuser (
				post_id, user_id,
				x, y,
				valid,
				createtime, updatetime
			)
			VALUES($1,$2,$3,$4,$5,$6,$7);
		`
	*/
	sqlStr = `
		INSERT INTO post_tag_xuser (
			post_id, user_id,
			x, y,
			valid,
			createtime, updatetime
		)
		VALUES 
	`
	args = append(args, postID)
	indexArg := 2
	timestamp := getNowUnixTimestamp()
	for i := 0; i < len(tags); i++ {
		if i != 0 {
			sqlStr += `, `
		}
		sqlStr += `($1,$` + strconv.Itoa(indexArg) + `,` +
			`$` + strconv.Itoa(indexArg+1) + `,` + `$` + strconv.Itoa(indexArg+2) + `,` +
			`$` + strconv.Itoa(indexArg+3) + `,` +
			`$` + strconv.Itoa(indexArg+4) + `,` + `$` + strconv.Itoa(indexArg+5) + `)`
		indexArg += 6
		args = append(args, tags[i].UserID, tags[i].X, tags[i].Y, false, timestamp, timestamp)
	}
	sqlStr += `ON CONFLICT ON CONSTRAINT post_tag_post_id_user_id_unique DO NOTHING;`
	return
}
