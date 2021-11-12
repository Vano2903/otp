package main

import "log"

var c Config

func init() {
	if err := c.Load(); err != nil {
		log.Fatal(err)
	}
}

func main() {

}
