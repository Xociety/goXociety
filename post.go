package main

func makeBlobURL(post postAPI) string {
	// restore url
	url := "http://" + bucketRootCloudStorage + "/" + makeBucketFolderName(post.Type, post.Blob.BlobID)
	switch postTypeMapID2Type[post.Type] {
	case mediaFormatJPG:
		url += "0." + mediaFormatJPG
	case mediaFormatHLS:
		url += "0." + mediaFormatM3U8
	}
	return url
}
