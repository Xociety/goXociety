package main

import (
	"fmt"
)

func main() {
	fmt.Println("hello xcociety")
	err := startInsertXuserfaker()
	fmt.Println("finished", err)
}
