package main

import (
	"log"

	"github.com/gengjosh/windwall/internal/server"
)

func main() {
	log.Println("Starting Windwall...")

	log.Fatal(server.Start())
}
