package main

import "alimmun-api/internal/config"

func main() {
	_, err := config.LoadConfig()

	if err != nil {
		panic(err)
	}

}
