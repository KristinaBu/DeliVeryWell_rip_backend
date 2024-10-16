package handler

import (
	"BMSTU_IU5_53B_rip/internal/app/repository"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	Repository *repository.Repository
	Logger     *logrus.Logger
}

func NewHandler(l *logrus.Logger, r *repository.Repository) *Handler {
	return &Handler{
		Logger:     l,
		Repository: r,
	}
}

const (
	DeliveryDomain = "/delivery"
	CallDomain     = "/call"
	RiDomain       = "/ri"
	UserDomain     = "/user"
)

func (h *Handler) RegisterHandler(router *gin.Engine) {
	/*
		router.GET("/", h.GetAllDelivery)
		router.GET(DeliveryDomain+"/:id", h.GetDelivery)
		router.POST("/delete/:id", h.DeleteDeliveryReq)
		router.POST("/add/:id", h.AddDeliveryToCall)
		router.GET("/mycalls/:callrequest_id", h.GetMyCallCards)
	*/

	// домен услуги /delivery
	router.GET(DeliveryDomain, h.GetAllDelivery)
	router.GET(DeliveryDomain+"/:id", h.GetDelivery)
	router.POST(DeliveryDomain, h.CreateDelivery) // without img
	router.POST(DeliveryDomain+"/img/:id", h.UploadImage)
	router.PUT(DeliveryDomain+"/:id", h.UpdateDelivery)
	router.DELETE(DeliveryDomain+"/:id", h.DeleteDelivery)
	router.POST(DeliveryDomain+"/add/:id", h.AddDeliveryToCall)

	// домен заявки /call
	router.GET(CallDomain, h.GetCalls)
	/*


		router.GET(CallDomain, h.GetMyCallCards)
		router.PUT(CallDomain+"/:id", h.UpdateCall)
		router.PUT(CallDomain+"/:id", h.FinishCall)
		router.DELETE(CallDomain+"/:id", h.DeleteCall)

		// домен м-м
		router.DELETE(RiDomain+"/delete/:id", h.DeleteRi)
		router.PUT(RiDomain+"/count/:id", h.UpdateRiCount)

		// домен пользователя
		router.POST(UserDomain, h.CreateUser)
		router.PUT(UserDomain+"/:id", h.Update)
		router.POST("/auth", h.AuthenticateUser)
		router.POST("/logout", h.LogoutUser)
	*/
}

func (h *Handler) RegisterStatic(router *gin.Engine) {
	router.LoadHTMLGlob("templates/*")
	router.Static("/static", "./static")
	router.Static("/css", "./static")
	router.Static("/img", "./static")
}

func (h *Handler) errorHandler(ctx *gin.Context, errorStatusCode int, err error) {
	h.Logger.Error(err.Error())
	ctx.JSON(errorStatusCode, gin.H{
		"status":      "error",
		"description": err.Error(),
	})
}
