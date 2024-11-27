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

	router.GET(DeliveryDomain, h.RoleMiddleware1(AdminRole, UserRole, GuestRole), h.GetAllDelivery)
	router.GET(DeliveryDomain+"/:id", h.GetDelivery)
	router.POST(DeliveryDomain, h.RoleMiddleware1(AdminRole), h.CreateDelivery) // without img
	router.POST(DeliveryDomain+"/img/:id", h.RoleMiddleware1(AdminRole), h.UploadImage)
	router.PUT(DeliveryDomain+"/:id", h.RoleMiddleware1(AdminRole), h.UpdateDelivery)
	router.DELETE(DeliveryDomain+"/:id", h.RoleMiddleware1(AdminRole), h.DeleteDelivery)
	router.POST(DeliveryDomain+"/add/:id", h.RoleMiddleware1(AdminRole, UserRole), h.AddDeliveryToCall)

	// домен заявки /call
	router.GET(CallDomain, h.RoleMiddleware1(AdminRole, UserRole), h.GetCalls)
	router.GET(CallDomain+"/:id", h.RoleMiddleware1(AdminRole, UserRole), h.GetCall)
	router.PUT(CallDomain+"/:id", h.RoleMiddleware1(AdminRole, UserRole), h.UpdateCall)
	router.PUT(CallDomain+"/form/:id", h.RoleMiddleware1(AdminRole), h.FormCall)
	router.PUT(CallDomain+"/complete/:id", h.RoleMiddleware1(AdminRole), h.CompleteOrRejectCall)
	router.DELETE(CallDomain+"/:id", h.RoleMiddleware1(AdminRole), h.DeleteCall)

	// домен м-м
	router.DELETE(RiDomain+"/delete/:id", h.RoleMiddleware1(AdminRole, UserRole), h.DeleteDC)
	router.PUT(RiDomain+"/count/:id", h.RoleMiddleware1(AdminRole, UserRole), h.UpdateDCCount)

	// домен пользователя
	router.POST(UserDomain+"/register", h.RegUser)
	router.PUT(UserDomain+"/update", h.UpdateUser)
	router.POST(UserDomain+"/login", h.AuthUser)
	router.POST(UserDomain+"/logout", h.LogoutUser)

	// для админа
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
