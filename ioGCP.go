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
func untarFileAndUpload(post postAPI, r io.Reader) error {
	// format checker
	jpgCount := 0
	m3u8Count := 0
	m3u8uploaded := false
	tsCount := 0

	foldername := makeBucketFolderName(post.Type, post.BlobID)

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
			if err := writeAndMakePublicCloudStorageGCP(client, bucketRootCloudStorage, foldername+filename+"."+fileFormat, tr); err != nil {
				log.Println("upload failed: ", err)
				return errors.New("upload failed")
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
