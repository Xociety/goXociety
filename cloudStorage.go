package main

import (
	"context"
	"io"
	"io/ioutil"
	"time"

	"cloud.google.com/go/storage"
)

// gcp cloud storage
func writeAndMakePublicCloudStorageGCP(client *storage.Client, bucket, object string, f io.Reader) error {
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
func writeDiscardIOGCP(f io.Reader) error {
	time.Sleep(time.Duration(randSeed.Intn(5)) * time.Millisecond)
	if _, err := io.Copy(ioutil.Discard, f); err != nil {
		return err
	}
	return nil
}
