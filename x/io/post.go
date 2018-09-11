package io

import "github.com/chienfuchen32/goXociety/x/config"

func MakeBlobURL(post config.PostAPI) string {
	// restore url
	url := "https://" + config.GlobalConfig[config.Env].GCPBucketRootCloudStorage + "/" + makeBucketFolderName(post.Type, post.Blob.BlobID)
	switch config.PostTypeMapID2Type[post.Type] {
	case config.MediaFormatJPG:
		url += "0." + config.MediaFormatJPG
	case config.MediaFormatHLS:
		url += "0." + config.MediaFormatM3U8
	}
	return url
}
