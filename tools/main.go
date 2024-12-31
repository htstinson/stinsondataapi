package main

import (
	"log"
	"os"

	salesforce "github.com/htstinson/stinsondataapi/api/internal/salesforce"
)

func main() {
	logger := log.New(os.Stdout, "[API] ", log.LstdFlags)

	logger.Println("initialize salesforce")
	sf, err := salesforce.New()
	if err != nil {
		logger.Println(err.Error())
		return
	}

}
