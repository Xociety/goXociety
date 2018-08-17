package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"strconv"
	"testing"
)

func TestLoadAdministrativeArea(t *testing.T) {
	type args struct {
		dataPath string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "init",
			args: args{dataPath: "./database/postgres/rg_cities1000.csv"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loadAdministrativeArea(tt.args.dataPath)
		})
	}
}

func TestParseCity2CityLevel(t *testing.T) { // write python version (use try catch to speed up map range)
	filePath := "/Volumes/Volume1/Projects/Python/shp2geojson/"
	countryCode := loadListFromFile(filePath + "countryList.txt")
	// log.Println(countryCode)
	for i := 0; i < len(countryCode); i++ {
		for j := 5; j >= 0; j-- {
			jj := strconv.Itoa(j)
			cities := []cityFromFile{}
			if file, err := ioutil.ReadFile(filePath + "geojson/gadm36_" + countryCode[i] + "_geojson/gadm36_" + countryCode[i] + "_" + jj + ".json"); err != nil {
				continue
			} else {
				if err := json.Unmarshal(file, &cities); err != nil {
					log.Panicln("gadm file parse fail", err)
				}
			}
			log.Println(i, countryCode[i], jj, "cities", len(cities))
			cls := []cityLevelAPI{}
			for k := 0; k < len(cities); k++ {
				cl := cityLevelAPI{}
				cl.GID = cities[k].Properties["GID_"+jj]
				cl.Name = cities[k].Properties["NAME_"+jj]
				cl.Type = cities[k].Properties["ENGTYPE_"+jj]
				cl.PostCount = 0
				cls = append(cls, cl)
			}
			b, err := json.Marshal(cls)
			if err != nil {
				log.Panicln("error:", err)
			}
			writeFile(filePath+"geojson/gadm36_"+countryCode[i]+"_geojson/city_level"+jj+".json", string(b))
		}
	}
}

func TestCityPostUpsertSample(t *testing.T) {
	// mongo upsert
	// c, err := connectMongoDB(globalConfig[env].MongoConStr)
	// if err != nil {
	// 	log.Panicln("mongo session", err)
	// }
	// collection := c.session.DB(mongoDBXociety).C("city_level" + jj)
	// log.Println(i, countryCode[i], jj, "cities", len(cities))
	// for k := 0; k < len(cities); k++ {
	// 	if _, err := collection.Upsert(
	// 		bson.M{"GID": cities[k].Properties["GID_"+jj]},
	// 		bson.M{"$set": bson.M{
	// 			"NAME":       cities[k].Properties["NAME_"+jj],
	// 			"TYPE":       cities[k].Properties["ENGTYPE_"+jj],
	// 			"POST_COUNT": 0,
	// 		}}); err != nil {
	// 		log.Println("upsert", err)
	// 	}
	// 	// log.Println(citi	es[k])
	// }
	// c.session.Close()
}
