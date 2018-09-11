package common

import (
	"strconv"

	"github.com/chienfuchen32/goXociety/x/config"
)

func ParseReactionCountSQL(reactionTable, reactionField, targetTable, targetRelatedField string) (sqlStr string) {
	/* to generate sql like:
	`
		UPDATE post
		SET ` + config.ReactionsMapID2Description[reactionOnPost.ReactionID] + `_count =
		(SELECT COUNT(*) FROM post_reaction
		WHERE post_reaction.post_id = post.post_id AND post_reaction.post_id = $1
			AND post_reaction.reaction_id = $2) WHERE post_id = $1;
	`
	*/
	sqlStr += `UPDATE ` + targetTable + ` SET `
	countSet := 0
	for k, v := range config.ReactionsMapID2Description {
		if countSet != 0 {
			sqlStr += `, `
		}
		sqlStr += v + `_count = (SELECT COUNT(*) FROM ` + reactionTable +
			` WHERE ` + reactionTable + `.` + targetRelatedField + `=` + targetTable + `.` + targetRelatedField +
			` AND ` + reactionTable + `.` + targetRelatedField + `= $1` +
			` AND ` + reactionTable + `.` + reactionField + `= ` + strconv.Itoa(k) + `)`
		countSet++
	}
	sqlStr += ` WHERE ` + targetRelatedField + ` = $1`
	return sqlStr
}
