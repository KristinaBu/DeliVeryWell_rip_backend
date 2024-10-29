package main

import (
	"BMSTU_IU5_53B_rip/internal/app/config"
	"BMSTU_IU5_53B_rip/internal/app/dsn"
	"BMSTU_IU5_53B_rip/internal/app/handler"
	"BMSTU_IU5_53B_rip/internal/app/pkg"
	"BMSTU_IU5_53B_rip/internal/app/repository"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// @title DeliVeryWell
// @version 1.0
// @description Delivery service

// @host 127.0.0.1
// @schemes http
// @BasePath /
func main() {
	logger := logrus.New()
	router := gin.Default()
	conf, errConf := config.NewConfig()
	if errConf != nil {
		logrus.Fatalln("Error with config reading: #{errConf}")
	}
	// через dsn парсим и помещаем в переменную
	postgresString := dsn.FromEnv()

	fmt.Println(postgresString)

	rep, err := repository.NewRepository(postgresString, logger)
	if err != nil {
		logrus.Fatalln("Error with repo: err", err)
	}

	hand := handler.NewHandler(logger, rep)
	application := pkg.NewApp(conf, router, logger, hand)
	application.StartServer()
}
