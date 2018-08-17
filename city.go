package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type administrativeArea struct {
	data map[string]map[string]map[string]map[string]string // map[countryCode][admin1][admin2][name]latLon
}

func loadAdministrativeArea(dataPath string) {
	/*
		lat, lon, name, admin1, admin2, cc
	*/
	aa := administrativeArea{data: make(map[string]map[string]map[string]map[string]string)}
	count := 0
	if file, err := os.Open("./config/postgres/rg_cities1000.csv"); err == nil {
		/*
			6.2,124.83333,"T""boli",Davao,,PH => 6.2,124.83333,T"boli,Davao,,PH
			38.89511,-77.03637,"Washington, D.C.","Washington, D.C.",,US => 38.89511,-77.03637,'Washington, D.C.','Washington, D.C.',,US
			48.73748,36.01987,"Yur""yivka",Dnipropetrovsk,,UA => 48.73748,36.01987,Yur"yivka,Dnipropetrovsk,,UA
		*/
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			rStr := []rune(scanner.Text())
			column := []string{}
			isQuotationMark := false
			// indexLastQuotationMark := -1
			str := ""
			for i := 0; i < len(rStr); i++ {
				char := string(rStr[i])
				if char == `"` {
					if i+1 < len(rStr) {
						if string(rStr[i+1]) == `"` {
							i = i + 1
							str += char
						} else {
							if isQuotationMark {
								isQuotationMark = false
								continue
							}
						}
						isQuotationMark = true
					}
				} else if char == `,` {
					if isQuotationMark {
						str += char
						continue
					}
					column = append(column, str)
					str = ""
					continue
				} else {
					str += char
				}
			}
			if len(str) > 0 {
				column = append(column, str)
				str = ""
			}
			if len(column) != 6 {
				log.Println("column len", count, len(column))
				for i := 0; i < len(column); i++ {
					log.Println(column[i])
				}
				continue
			}
			latLon := column[0] + "," + column[1]
			name := column[2]
			admin1 := column[3]
			admin2 := column[4]
			cc := column[5]
			if aa.data[cc] == nil {
				aa.data[cc] = make(map[string]map[string]map[string]string)
			}
			if aa.data[cc][admin1] == nil {
				aa.data[cc][admin1] = make(map[string]map[string]string)
			}
			if aa.data[cc][admin1][admin2] == nil {
				aa.data[cc][admin1][admin2] = make(map[string]string)
			}
			aa.data[cc][admin1][admin2][name] = latLon
			count++
		}
	}
	b, err := json.Marshal(aa.data)
	if err != nil {
		fmt.Println("error:", err)
	}
	writeFile("./development/city.json", string(b))
}
