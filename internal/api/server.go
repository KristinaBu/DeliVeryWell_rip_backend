package api

import (
	"BMSTU_IU5_53B_rip/internal/transport"
	"github.com/gin-gonic/gin"
	"log"
)

func StartServer() {
	log.Println("Server start up")

	router := gin.Default()

	// считывание шаблонов
	router.LoadHTMLGlob("templates/*")
	router.Static("/static", "./static/")

	transport.InitializeRouters(router)

	router.Run() //  0.0.0.0:8080 ("localhost:8080")

	log.Println("Server down")
}
