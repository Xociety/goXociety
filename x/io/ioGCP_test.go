package io

import (
	"context"
	"log"
	"testing"

	"googlemaps.github.io/maps"
)

func TestGooglePlaceAPI(t *testing.T) {
	c, err := maps.NewClient(maps.WithAPIKey(config.GoogleMapKey))
	if err != nil {
		log.Fatalln("client ", err)
	}
	// // FindPlaceFromTextRequest
	// r := &maps.FindPlaceFromTextRequest{
	// 	Input:     "Museum of Contemporary Art Australia",
	// 	InputType: maps.FindPlaceFromTextInputType("textquery"),
	// 	Fields:    []maps.PlaceSearchFieldMask{"formatted_address", "geometry", "id", "name", "place_id", "types"},
	// }
	// resp, err := c.FindPlaceFromText(context.Background(), r)
	// log.Println("candidate", resp.Candidates)
	// r := &maps.PlaceAutocompleteRequest{
	// 	Input: "SHINKONG MITSUKOSHI", // "taipei 101",
	// 	// Location:
	// 	SessionToken: maps.PlaceAutocompleteSessionToken(uuid.New()),
	// }
	// resp, err := c.PlaceAutocomplete(context.Background(), r)
	// if err != nil {
	// 	log.Fatalln("PlaceAutocomplete", err)
	// }
	// for i := 0; i < len(resp.Predictions); i++ {
	// 	log.Println("Des", resp.Predictions[i].Description)
	// 	log.Println("MatchedSubstrings", resp.Predictions[i].MatchedSubstrings)
	// 	log.Println("ID", resp.Predictions[i].PlaceID)
	// 	log.Println("str", resp.Predictions[i].StructuredFormatting)
	// 	log.Println("term", resp.Predictions[i].Terms)
	// 	log.Println("type", resp.Predictions[i].Types)
	// }
	// //	// NearbySearchRequest
	// r := &maps.NearbySearchRequest{
	// 	Location: &maps.LatLng{Lat: 25.033667, Lng: 121.564022},
	// 	Radius:   100,
	// 	// Keyword:  "cruise",
	// 	Type:     maps.PlaceType(""),
	// 	Language: "en", //"zh-TW",
	// 	// PageToken: "1333",
	// }
	// resp, err := c.NearbySearch(context.Background(), r)
	// if err != nil {
	// 	log.Fatalln("NearbySearch", err)
	// }
	// log.Println(len(resp.Results), " results")
	// for i := 0; i < len(resp.Results); i++ {
	// 	log.Println("AltIDs", resp.Results[i].AltIDs)
	// 	log.Println("FormattedAddress", resp.Results[i].FormattedAddress)
	// 	log.Println("Geometry", resp.Results[i].Geometry)
	// 	log.Println("Icon", resp.Results[i].Icon)
	// 	log.Println("ID", resp.Results[i].ID)
	// 	log.Println("Name", resp.Results[i].Name)
	// 	log.Println("OpeningHours", resp.Results[i].OpeningHours)
	// 	log.Println("PermanentlyClosed", resp.Results[i].PermanentlyClosed)
	// 	log.Println("Photos", resp.Results[i].Photos)
	// 	log.Println("PlaceID", resp.Results[i].PlaceID)
	// 	log.Println("PriceLevel", resp.Results[i].PriceLevel)
	// 	log.Println("Rating", resp.Results[i].Rating)
	// 	log.Println("Scope", resp.Results[i].Scope)
	// 	log.Println("Types", resp.Results[i].Types)
	// }
	// // TextSearchRequest
	r := &maps.TextSearchRequest{
		Query:  "restaurants+in+Sydney",
		Radius: 1000,
	}
	resp, err := c.TextSearch(context.Background(), r)
	if err != nil {
		log.Fatalln("TextSearch", err)
	}
	log.Println(len(resp.Results), " results")
	for i := 0; i < len(resp.Results); i++ {
		log.Println("FormattedAddress", resp.Results[i].FormattedAddress)
		log.Println("Geometry", resp.Results[i].Geometry)
		log.Println("Name", resp.Results[i].Name)
	}
}
