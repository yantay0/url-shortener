package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func initEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

func main() {
	initEnv()
	fmt.Printf("port is %s\n", os.Getenv("SERVER_PORT"))
}
