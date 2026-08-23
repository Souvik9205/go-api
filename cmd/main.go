package main

import (
	"log"
	"os"
)

func main() {
	cfg := config{
		addr: ":8080",
		db:   dbConfig{},
	}

	api := applicattion{
		config: cfg,
	}

	h := api.mount()
	if err := api.run(h); err != nil {
		log.Printf("Server has failed to start, err: %s", err)
		os.Exit(1)
	}
}
