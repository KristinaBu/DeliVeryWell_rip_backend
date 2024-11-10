package handler

import (
	"BMSTU_IU5_53B_rip/docs"
	"BMSTU_IU5_53B_rip/internal/app/ds"
	"BMSTU_IU5_53B_rip/internal/app/repository"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
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
	RiDomain       = "/dc"
	UserDomain     = "/user"
)

func (h *Handler) RegisterHandler(router *gin.Engine) {

	docs.SwaggerInfo.Title = "DeliVeryWell"
	docs.SwaggerInfo.Description = "Delivery service"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/"

	// домен услуги /delivery
	router.GET(DeliveryDomain, h.GetAllDelivery)
	router.GET(DeliveryDomain+"/:id", h.GetDelivery)
	router.POST(DeliveryDomain, h.RoleMiddleware(ds.User{IsAdmin: true}),
		h.CreateDelivery) // without img
	router.POST(DeliveryDomain+"/img/:id", h.RoleMiddleware(ds.User{IsAdmin: true}),
		h.UploadImage)
	router.PUT(DeliveryDomain+"/:id", h.RoleMiddleware(ds.User{IsAdmin: true}),
		h.UpdateDelivery)
	router.DELETE(DeliveryDomain+"/:id", h.RoleMiddleware(ds.User{IsAdmin: true}),
		h.DeleteDelivery)
	router.POST(DeliveryDomain+"/add/:id", h.RoleMiddleware(ds.User{IsAdmin: true}, ds.User{IsAdmin: false}),
		h.AddDeliveryToCall)

	// домен заявки /call
	router.GET(CallDomain, h.RoleMiddleware(ds.User{IsAdmin: true}, ds.User{IsAdmin: false}),
		h.GetCalls)
	router.GET(CallDomain+"/:id", h.RoleMiddleware(ds.User{IsAdmin: true}, ds.User{IsAdmin: false}),
		h.GetCall)
	router.PUT(CallDomain+"/:id", h.RoleMiddleware(ds.User{IsAdmin: true}, ds.User{IsAdmin: false}),
		h.UpdateCall)
	router.PUT(CallDomain+"/form/:id", h.RoleMiddleware(ds.User{IsAdmin: true}),
		h.FormCall)
	router.PUT(CallDomain+"/complete/:id", h.RoleMiddleware(ds.User{IsAdmin: true}),
		h.CompleteOrRejectCall)
	router.DELETE(CallDomain+"/:id", h.RoleMiddleware(ds.User{IsAdmin: true}),
		h.DeleteCall)

	// домен м-м
	router.DELETE(RiDomain+"/delete/:id", h.RoleMiddleware(ds.User{IsAdmin: true}, ds.User{IsAdmin: false}),
		h.DeleteDC)
	router.PUT(RiDomain+"/count/:id", h.RoleMiddleware(ds.User{IsAdmin: true}, ds.User{IsAdmin: false}),
		h.UpdateDCCount)

	// домен пользователя
	router.POST(UserDomain+"/register", h.RegUser)
	router.PUT(UserDomain+"/update", h.UpdateUser)
	router.POST(UserDomain+"/login", h.AuthUser)
	router.POST(UserDomain+"/logout", h.LogoutUser)

	router.GET(UserDomain+"/protected", h.RoleMiddleware(ds.User{IsAdmin: true}), func(ctx *gin.Context) {
		userID := ctx.MustGet("user_id").(uint)

		ctx.JSON(http.StatusOK, gin.H{
			"message": "user is autorized",
			"user_id": userID,
		})
	})

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
