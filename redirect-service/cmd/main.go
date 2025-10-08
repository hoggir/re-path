package main

import (
	"log"
)

func main() {
	router := InitializeApp()
	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
