package transport

import "github.com/gin-gonic/gin"

func InitializeRouters(router *gin.Engine) {
	router.GET("/", showIndexPage)

	router.GET("/card/:card_id", getCallCard)
	router.GET("/mycalls", getMyCallCards)

}
