package render

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Render(c *gin.Context, templateName string, data gin.H) {
	c.HTML(http.StatusOK, templateName, data)
}
