package main

import (
	"fmt"
	"winston"
)

func main() {
	var w winston.Winston

	w.FetchUrl("http://www.google.com")
	fmt.Println(w.Text)
}
