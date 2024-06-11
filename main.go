package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"webscraping/discord"

	"github.com/joho/godotenv"
)

func main() {
	loadEnv()

	discord.Run()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
	log.Println("Interupted")
}

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	log.Println("[👍] Loaded .env")
}
