package main

import (
	"fmt"

	"notebook_app/internal/app"
)

func main() {
	err := app.RunBackend()
	if err == nil {
		fmt.Println("\n\n\nВсе ок")
	} else {
		fmt.Printf("\n\n\nЧто-то пошло не так: %s\n", err.Error())
	}
}
