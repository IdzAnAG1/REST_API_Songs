package main

import "REST_API_Songs/internal/config"

import (
	"fmt"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Sprintf("Ай ай")
	}
	fmt.Println(cfg)
}
