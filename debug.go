package main

import (
	"log"
	"os"
)

var debug *log.Logger

func init() {
	file, err := os.OpenFile("/tmp/provider", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		panic(err)
	}
	debug = log.New(file, "", 0)
}
