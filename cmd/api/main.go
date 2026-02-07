package main

import (
	"job-tracker/internal/bootstrap"
	"log"
)

func main() {

	if err := bootstrap.Start(); err != nil {
		log.Fatal(err)
	}

}
