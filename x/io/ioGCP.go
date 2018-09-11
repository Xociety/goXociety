package io

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"strconv"
	"strings"

	"cloud.google.com/go/storage"
	"github.com/chienfuchen32/goXociety/x/config"
	"googlemaps.github.io/maps"
)

func makeBucketFolderName(postType int, blobID string) (foldername string) {
	// define cloud storage folder under bucket
	switch config.PostTypeMapID2Type[postType] {
	case config.MediaFormatJPG:
		foldername += config.BucketImagesCloudStorage
	case config.MediaFormatHLS:
		foldername += config.BucketVideosCloudStorage
	}
	foldername += "/" + blobID + "/"
	return foldername
}
func untarFileAndUpload(post config.PostAPI, r io.Reader, isGCP bool) error {
	// format checker
	jpgCount := 0
	m3u8Count := 0
	m3u8uploaded := false
	tsCount := 0

	foldername := makeBucketFolderName(post.Type, post.Blob.BlobID)

	// define gcp client
	client, err := storage.NewClient(context.Background(), config.ClientOptionGoogleAPI)
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
		if !(fileFormat == config.MediaFormatJPG || fileFormat == config.MediaFormatM3U8 || fileFormat == config.MediaFormatTS) {
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
		switch config.PostTypeMapID2Type[post.Type] {
		case config.MediaFormatJPG:
			if fileFormat == config.MediaFormatJPG {
				filename = strconv.Itoa(jpgCount) + "." + fileFormat
				jpgCount++
				isUpload = true
			}
		case config.MediaFormatHLS:
			if !m3u8uploaded && fileFormat == config.MediaFormatM3U8 {
				filename = strconv.Itoa(m3u8Count) + "." + fileFormat
				isUpload = true
				m3u8uploaded = true
			} else if fileFormat == config.MediaFormatTS {
				tsCount++
				isUpload = true
			}
		}
		if isUpload {
			if isGCP {
				if err := WriteAndMakePublicCloudStorageGCP(client, config.GlobalConfig[config.Env].GCPBucketRootCloudStorage, foldername+filename, tr); err != nil {
					log.Println("upload failed: ", err)
					return errors.New("upload failed")
				}
			} else {
				if err := WriteDiscardIOGCP(tr); err != nil {
					log.Println("fake upload failed: ", err)
					return errors.New("fake upload failed")
				}
			}
		}
	}

	switch config.PostTypeMapID2Type[post.Type] {
	case config.MediaFormatJPG:
		if jpgCount == 0 {
			return errors.New("jpg format wrong")
		}
	case config.MediaFormatHLS:
		if !m3u8uploaded || tsCount == 0 {
			return errors.New("hls format wrong")
		}
	}
	return nil
}

func getPlaceByLocationGCP(lat, lon float64, keyword, pageToken string) (places []config.PlaceAPI, nextPageToken string, err error) {
	c, err := maps.NewClient(maps.WithAPIKey(config.GoogleMapKey))
	if err != nil {
		log.Println("google map client ", err)
		return places, "", err
	}
	r := &maps.NearbySearchRequest{
		Location: &maps.LatLng{Lat: lat, Lng: lon},
		Radius:   config.RadiusGoogleMap,
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
		place := config.PlaceAPI{
			Name: resp.Results[i].Name,
			Lat:  resp.Results[i].Geometry.Location.Lat,
			Lon:  resp.Results[i].Geometry.Location.Lng,
		}
		places = append(places, place)
	}
	return places, resp.NextPageToken, nil
}

func getPlaceByNameGCP(keyword, pageToken string) (places []config.PlaceAPI, nextPageToken string, err error) {
	c, err := maps.NewClient(maps.WithAPIKey(config.GoogleMapKey))
	if err != nil {
		log.Println("google map client ", err)
		return places, "", err
	}
	r := &maps.TextSearchRequest{
		Query:  keyword,
		Radius: config.RadiusGoogleMap,
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
		place := config.PlaceAPI{
			Name: resp.Results[i].Name,
			Lat:  resp.Results[i].Geometry.Location.Lat,
			Lon:  resp.Results[i].Geometry.Location.Lng,
		}
		places = append(places, place)
	}
	return places, resp.NextPageToken, nil
}

func WriteAndMakePublicCloudStorageGCP(client *storage.Client, bucket, object string, f io.Reader) error {
	ctx := context.Background()
	// [START upload_file]
	wc := client.Bucket(bucket).Object(object).NewWriter(ctx)
	if _, err := io.Copy(wc, f); err != nil {
		return err
	}
	if err := wc.Close(); err != nil {
		return err
	}
	acl := client.Bucket(bucket).Object(object).ACL()
	if err := acl.Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
		return err
	}
	// [END upload_file]
	return nil
}
func WriteDiscardIOGCP(f io.Reader) error {
	// time.Sleep(time.Duration(config.RandSeed.Intn(5)) * time.Millisecond)
	if _, err := io.Copy(ioutil.Discard, f); err != nil {
		return err
	}
	return nil
}
