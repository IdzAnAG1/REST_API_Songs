package main

import (
	"REST_API_Songs/internal/config"
	"REST_API_Songs/internal/db"
	"context"
	"github.com/sirupsen/logrus"
)

import (
	"fmt"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Sprintf("Ай ай")
	}
	dab, err := db.NewDataBase(cfg.DataBaseURL)
	if err != nil {
		fmt.Errorf("Провалилось подклюение к БД")
	}
	defer dab.CloseDatabase()

	err = dab.Cnct.Ping(context.Background())
	if err != nil {
		logrus.Fatal("Ошибка при проверке подключения к БД: %v", err)
	}
}
