package main

import (
	"strconv"
)

func parsePlaceSelectAllSQL(name []string, lat, lon []float64) (sqlStr string, args []interface{}) {
	/*
		in order to select multiple place on post in one sql command, this func parse the command and parameters
		basic select:
		sqlStr := `
			SELECT place_id, country_code,
			city_id_1, city_id_2, city_id_3,
			city_id_4, city_id_5,
			lat, lon, name,
			FROM place
			WHERE name=$1 AND lat=$2 AND lon=$3;
		`
	*/
	sqlStr = `
	SELECT place.place_id, 
	place.country_code,
	place.city_id_1, place.city_id_2, place.city_id_3,
	place.city_id_4, place.city_id_5,
	place.lat, place.lon, place.name 
	FROM place WHERE 
	`
	indexArg := 1
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
