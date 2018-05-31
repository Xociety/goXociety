package main

import (
	"fmt"
)

const serverPort = 4000

const graphqlRoute = "/graphql"
const graphiqlRoute = "/graphiql"

const postgresConStr = "host=localhost port=31160 user=postgres password=mysecretpassword sslmode=disable"

func main() {
	fmt.Println("hello xcociety")
	// err := startInsertXuserfaker()
	// fmt.Println("finished", err)
	startServer()
}
