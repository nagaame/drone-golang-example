package main

import (
	"demo/util"
	"net/http"
)

func main() {
	http.HandleFunc("/", util.SayHello)

	if err := http.ListenAndServe(":8088", nil); err != nil {
		panic(err)
	}
}

