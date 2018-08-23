package main

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"errors"
	"io"
	"log"
	"strconv"
	"strings"

	"cloud.google.com/go/storage"
	"googlemaps.github.io/maps"
)

func makeBucketFolderName(postType int, blobID string) (foldername string) {
	// define cloud storage folder under bucket
	switch postTypeMapID2Type[postType] {
	case mediaFormatJPG:
		foldername += bucketImagesCloudStorage
	case mediaFormatHLS:
		foldername += bucketVideosCloudStorage
	}
	foldername += "/" + blobID + "/"
	return foldername
}
func untarFileAndUpload(post postAPI, r io.Reader, isGCP bool) error {
	// format checker
	jpgCount := 0
	m3u8Count := 0
	m3u8uploaded := false
	tsCount := 0

	foldername := makeBucketFolderName(post.Type, post.Blob.BlobID)

	// define gcp client
	client, err := storage.NewClient(context.Background(), clientOptionGoogleAPI)
	defer client.Close()
	if err != nil {
		log.Println("gcp sdk client", err)
	}

	// untar with gz
	gzr, err := gzip.NewReader(r)
	defer gzr.Close()
	if err != nil {
		log.Panicln("gz err", err)
		return err
	}

	// reader
	tr := tar.NewReader(gzr)
	for {
		header, err := tr.Next()

		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		if header == nil {
			continue
		}

		formatStrSplitter := strings.Split(header.Name, ".")
		if len(formatStrSplitter) == 1 {
			continue
		}
		fileFormat := strings.ToLower(formatStrSplitter[len(formatStrSplitter)-1])
		if !(fileFormat == mediaFormatJPG || fileFormat == mediaFormatM3U8 || fileFormat == mediaFormatTS) {
			continue
		}

		filename := ""
		directoryStrSplitter := strings.Split(header.Name, "/")
		if len(directoryStrSplitter) > 1 {
			filename = directoryStrSplitter[len(directoryStrSplitter)-1]
		}
		if len(filename) >= 1 {
			if filename[:1] == "." {
				continue
			}
		}

		isUpload := false
		switch postTypeMapID2Type[post.Type] {
		case mediaFormatJPG:
			if fileFormat == mediaFormatJPG {
				filename = strconv.Itoa(jpgCount) + "." + fileFormat
				jpgCount++
				isUpload = true
			}
		case mediaFormatHLS:
			if !m3u8uploaded && fileFormat == mediaFormatM3U8 {
				filename = strconv.Itoa(m3u8Count) + "." + fileFormat
				isUpload = true
				m3u8uploaded = true
			} else if fileFormat == mediaFormatTS {
				tsCount++
				isUpload = true
			}
		}
		if isUpload {
			if isGCP {
				if err := writeAndMakePublicCloudStorageGCP(client, globalConfig[env].GCPBucketRootCloudStorage, foldername+filename, tr); err != nil {
					log.Println("upload failed: ", err)
					return errors.New("upload failed")
				}
			} else {
				if err := writeDiscardIOGCP(tr); err != nil {
					log.Println("fake upload failed: ", err)
					return errors.New("fake upload failed")
				}
			}
		}
	}

	switch postTypeMapID2Type[post.Type] {
	case mediaFormatJPG:
		if jpgCount == 0 {
			return errors.New("jpg format wrong")
		}
	case mediaFormatHLS:
		if !m3u8uploaded || tsCount == 0 {
			return errors.New("hls format wrong")
		}
	}
	return nil
}

func getPlaceByLocationGCP(lat, lon float64, keyword, pageToken string) (places []placeAPI, nextPageToken string, err error) {
	c, err := maps.NewClient(maps.WithAPIKey(googleMapKey))
	if err != nil {
		log.Println("google map client ", err)
		return places, "", err
	}
	r := &maps.NearbySearchRequest{
		Location: &maps.LatLng{Lat: lat, Lng: lon},
		Radius:   radiusGoogleMap,
		Language: "en",
	}
	if keyword != "" {
		r.Keyword = keyword
	}
	if pageToken != "" {
		r.PageToken = pageToken
	}
	resp, err := c.NearbySearch(context.Background(), r)
	if err != nil {
		log.Println("getPlaceByLocationGCP", err)
		return places, "", err
	}
	for i := 0; i < len(resp.Results); i++ {
		place := placeAPI{
			Name: resp.Results[i].Name,
			Lat:  resp.Results[i].Geometry.Location.Lat,
			Lon:  resp.Results[i].Geometry.Location.Lng,
		}
		places = append(places, place)
	}
	return places, resp.NextPageToken, nil
}

func getPlaceByNameGCP(keyword, pageToken string) (places []placeAPI, nextPageToken string, err error) {
	c, err := maps.NewClient(maps.WithAPIKey(googleMapKey))
	if err != nil {
		log.Println("google map client ", err)
		return places, "", err
	}
	r := &maps.TextSearchRequest{
		Query:  "restaurants+in+Sydney",
		Radius: radiusGoogleMap,
	}
	if pageToken != "" {
		r.PageToken = pageToken
	}
	resp, err := c.TextSearch(context.Background(), r)
	if err != nil {
		log.Println("getPlaceByNameGCP", err)
		return places, "", err
	}
	for i := 0; i < len(resp.Results); i++ {
		place := placeAPI{
			Name: resp.Results[i].Name,
			Lat:  resp.Results[i].Geometry.Location.Lat,
			Lon:  resp.Results[i].Geometry.Location.Lng,
		}
		places = append(places, place)
	}
	return places, resp.NextPageToken, nil
}
