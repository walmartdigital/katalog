package main

import (
	"fmt"

	"github.com/verdverm/frisby"
)

func main() {
	fmt.Println("Frisby!")

	frisby.Create("Test GET Go homepage").
		Get("http://golang.org").
		Send().
		ExpectStatus(200).
		ExpectContent("The Go Programming Language")

	frisby.Global.PrintReport()
}
