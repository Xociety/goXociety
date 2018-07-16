# goXociety


## golang dependency

`go get github.com/manveru/faker`

`go get github.com/lib/pq`

`go get github.com/graphql-go/graphql`

`go get github.com/chienfuchen32/goXHandler`

`go get -u cloud.google.com/go/storage`

`go get firebase.google.com/go`

`go get github.com/globalsign/mgo`

## config

reaction.csv, post_type.csv

## database

postgres 9.6 at least: user, post, config ...etc.

mongo 4.0: popular post, popular post read.

## image, video

upload tar.gz

sample [resource](https://videos.pexels.com/)

## cert

currently using cloudflare orgin ssl

## start develope this repo

`env=development go run $(ls -1 *.go | grep -v _test.go)`

