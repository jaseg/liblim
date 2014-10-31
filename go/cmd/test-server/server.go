package main

import (
	"log"

	derrit "../../"
)

func main() {
	log.Println("Starting Test Server")
	log.Println(derrit.Listen("0.0.0.0:4242"))

}
