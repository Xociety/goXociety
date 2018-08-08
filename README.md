# goXociety

## golang dependency

`go get github.com/manveru/faker`

`go get github.com/lib/pq`

`go get github.com/graphql-go/graphql`

`go get github.com/chienfuchen32/goXHandler`

`go get -u cloud.google.com/go/storage`

`go get firebase.google.com/go`

`go get github.com/globalsign/mgo`

## dns

| Type   |      Name      |  Value |
|----------|:-------------:|------:|
| txt |  @ | google site ownership verify string (deletable) |
| cname |  storage | c.storage.googleapis.com |
| a |    api   |   kubernetes cluster service api extenal ip |
| a |    cron   |   mongo compute engine extenal ip |
| a |    status   |   cloud sql extenal ip |

ps: google site ownership verify with dns TXT should be easiest way.

## goXociety repo 

* config: `./config/xocitetyConfig.json`

* cert: currently using cloudflare orgin ssl `./secret/cloudflare/server.cert, server.key`

* graphql api: please check graphiql url

## database

### postgres 9.6: stores user, post, config ...etc.

* table definition: `./database/postgres/create_table.sql`

* table config data

`./config/postgres/category.csv, country.csv, gender.csv, language.csv, reaction.csv, post_type.csv`

* data import

init: `./config/postgres/init`

ps: development with psql shell, deployment with cloud sql import

patch: `./config/postgres/patch`

- development (no ssl)

brew services start postgresql@9.6

log: `/usr/local/var/log/postgresql@9.6.log`

- production

on google compute engine

google cloud signed cert: `./secret/postgres/client-cert.pem, client-key.pem, server-ca.pem`

### mongo 4.0.0:

* collection init

`./database/mongo/create_collection.js`

* data, config, ssl, log path

development (no ssl)

brew services start mongodb

`/usr/local/var/mongodb, /usr/local/mongod.conf, /usr/local/etc/ssl, /usr/local/var/log/mongodb/mongo.log`

- production

* install [tutorial](https://docs.mongodb.com/manual/tutorial/install-mongodb-on-debian/) on google compute engine

`./config/mongo/production/mongod.conf`

`/var/lib/mongodb, /etc/mongod.conf, /etc/ssl, /var/log/mongodb`

* self sign cert and config sample

[tutorial](https://gist.github.com/chienfuchen32/73e451039be37b86deec97572fb3d4fd)

`./secret/mongo/*`

* gcp firewall => [mongo] map to mongo compute engine network label

## image, video sample

image: `deployment/upload/sample/image.tar.gz`

video: `deployment/upload/sample/playlist.tar.gz`

sample [resource](https://videos.pexels.com/)

## start develope this repo

`env=development go run $(ls -1 *.go | grep -v _test.go)`

## deploy goXociety on kubernetes

* configmap: `./config`

* secret `./secret`

## deploy on gke

local:

`sh build.sh`

gcp:

`gcloud container builds submit --config cloudbuild.yaml .`