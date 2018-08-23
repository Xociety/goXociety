package main

import (
	"strconv"
)

func parsePlaceSelectSQL(name []string, lat, lon []float64) (sqlStr string, args []interface{}) {
	/*
		in order to select multiple place on post in one sql command, this func parse the command and parameters
		basic select:
		sqlStr := `
			SELECT place_id
			FROM place
			WHERE name=$1 AND lat=$2 AND lon=$3;
		`
	*/
	sqlStr = `SELECT place.place_id FROM place WHERE `
	indexArg := 0
	for i := 0; i < len(name); i++ {
		if i != 0 {
			sqlStr += `OR `
		}
		sqlStr += `(name=$` + strconv.Itoa(indexArg) + ` AND lat=$` + strconv.Itoa(indexArg+1) + ` AND lon=$` + strconv.Itoa(indexArg+2) + `) `
		indexArg += 3
		args = append(args, name[i], lat[i], lon[i])
	}
	sqlStr += `;`
	return
}
