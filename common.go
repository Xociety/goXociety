package main

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"log"
	"strings"
	"time"

	"cloud.google.com/go/storage"
)

func getNowUnixTimestamp() int {
	return int(time.Now().Unix())
}

func getSHA256Hash(input []byte) string {
	h := sha256.New()
	h.Write(input)
	return hex.EncodeToString(h.Sum(nil))
}

func untarFileAndUpload(post postDB, r io.Reader) error {

	// format checker
	jpgCount := 0
	m3u8uploaded := false
	tsCount := 0

	// define cloud storage folder under bucket
	foldername := ""
	switch post.Type {
	case 0:
		foldername += bucketImagesCloudStorage
	case 1:
		foldername += bucketVideosCloudStorage
	}
	foldername += "/" + post.BlobID + "/"

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
		if !(fileFormat == "jpg" || fileFormat == "m3u8" || fileFormat == "ts") {
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
		switch post.Type { // type load from config
		case 0:
			if fileFormat == "jpg" {
				jpgCount++
				isUpload = true
			}
		case 1:
			if !m3u8uploaded && fileFormat == "m3u8" {
				isUpload = true
				m3u8uploaded = true
			} else if fileFormat == "ts" {
				tsCount++
				isUpload = true
			}
		}
		if isUpload {
			if err := writeAndMakePublicCloudStorageGCP(client, bucketRootCloudStorage, foldername+filename, tr); err != nil {
				log.Println("upload failed: ", err)
				return errors.New("upload failed")
			}
		}
	}

	switch post.Type {
	case 0:
		if jpgCount == 0 {
			return errors.New("jpg format wrong")
		}
	case 1:
		if !m3u8uploaded || tsCount == 0 {
			return errors.New("hls format wrong")
		}
	}
	return nil
}
