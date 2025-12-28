package main

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func main() {

	err := bcrypt.CompareHashAndPassword([]byte("$2a$10$EkIApwhf1kKCijeASFYU5OWeROGYtQtC3b0tg4j028QS7EvjAGcK2"), []byte("password"))
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("match")
	}

}
