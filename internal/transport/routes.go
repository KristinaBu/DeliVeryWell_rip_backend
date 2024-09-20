package transport

import "github.com/gin-gonic/gin"

func InitializeRouters(router *gin.Engine) {
	router.GET("/", showIndexPage)

	router.GET("/aboutcall/:card_id", getCallCard)
	router.GET("/mycalls/:callrequest_id", getMyCallCards)

}
