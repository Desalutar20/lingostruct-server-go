package main

import (
	"fmt"

	"github.com/Desalutar20/lingostruct-server-go/internal/config"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
}

func main() {
	config := config.New()

	fmt.Println(config)
}
