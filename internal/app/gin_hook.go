package app

import (
	"fmt"
	"task_service/internal/data"
	"task_service/internal/service/controller"
	"task_service/pkg/database"

	"github.com/gin-gonic/gin"
)

var ctrl *controller.Controller

func initCtrl(app *Application, r *gin.Engine) error {

	gormCli, err := database.InitGormClient(app.GetDatabase())
	if err != nil {
		fmt.Errorf("initCtrl: %s", err.Error())
	}

	dataMgr := data.NewDataManager(gormCli)
	cacheMgr := data.NewCacheMgr(app.cacheClient)

	ctrl := controller.NewController(dataMgr, cacheMgr)

	v1Group := r.Group("task-service/api/v1")
	v1Group.GET("/tasks/:taskId", ctrl.GetTask)
	v1Group.GET("/tasks", ctrl.ListTask)
	v1Group.POST("/tasks", ctrl.CreateTask)
	v1Group.PUT("/tasks/:taskId", ctrl.UpdateTask)
	v1Group.DELETE("/tasks/:taskId", ctrl.DeleteTask)

	return nil
}

func InitGinApplicationHook(app *Application) error {
	if app.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	gin.EnableJsonDecoderUseNumber()

	r := gin.New()
	r.Use(gin.Recovery())

	initCtrl(app, r)
	addr := fmt.Sprintf("%s:%s", app.GetConfig().Service.Host, app.GetConfig().Service.Port)

	app.SetAddr(addr)
	app.SetSrv(r)

	return nil
}

func DestroyGinApplicationHook(app *Application) error {
	ctrl.Shutdown()
	return nil
}
