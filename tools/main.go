package main

import (
	"log"
	"os"

	"github.com/htstinson/stinsondataapi/tree/main/api/internal/salesforce"
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
