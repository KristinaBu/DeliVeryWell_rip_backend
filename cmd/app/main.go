package main

import (
	"BMSTU_IU5_53B_rip/internal/api"

	"log"
)

//TODO
// развертывание Minio

func main() {
	log.Println("Application start!")
	api.StartServer()
	log.Println("Application terminated!")
}
