FROM golang:1.10.3
RUN apt-get update
RUN apt-get install -y vim sudo 
RUN go get github.com/manveru/faker
RUN go get github.com/lib/pq
RUN go get github.com/graphql-go/graphql
RUN go get github.com/chienfuchen32/handler
RUN go get -u cloud.google.com/go/storage
RUN go get firebase.google.com/go
RUN go get github.com/graphql-go/graphql
RUN go get github.com/globalsign/mgo
WORKDIR /goXociety
ADD . .
RUN go build
