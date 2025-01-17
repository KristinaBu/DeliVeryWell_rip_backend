package pkg

import (
	"BMSTU_IU5_53B_rip/internal/app/config"
	"BMSTU_IU5_53B_rip/internal/app/handler"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Application struct {
	Config  *config.Config
	Logger  *logrus.Logger
	Router  *gin.Engine
	Handler *handler.Handler
}

func NewApp(c *config.Config, r *gin.Engine, l *logrus.Logger, h *handler.Handler) *Application {
	return &Application{
		Config:  c,
		Logger:  l,
		Router:  r,
		Handler: h,
	}
}

func (a *Application) StartServer() {
	a.Logger.Info("Server start up :)")
	a.Handler.RegisterHandler(a.Router)
	a.Handler.RegisterStatic(a.Router)
	serverAddress := fmt.Sprintf("%s:%d", a.Config.ServiceHost, a.Config.ServicePort)

	// TODO: новые роуты

	if err := a.Router.Run(serverAddress); err != nil {
		a.Logger.Fatalln(err)
	}
	a.Logger.Info("Server down")
}
