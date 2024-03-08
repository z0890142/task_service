package controller

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"task_service/c"
	"task_service/internal/data"
	"task_service/pkg/logger"
	"task_service/pkg/models"

	"github.com/gin-gonic/gin"
	"google.golang.org/genproto/googleapis/rpc/code"
)

type Controller struct {
	mysqlMgr      data.DataManager
	cacheMgr      data.DataManager
	enableCache   bool
	shuntDownOnce sync.Once
}

func NewController(mysqlMgr, cacheMgr data.DataManager) *Controller {
	return &Controller{
		mysqlMgr:      mysqlMgr,
		cacheMgr:      cacheMgr,
		enableCache:   false,
		shuntDownOnce: sync.Once{},
	}
}

func (ctrl *Controller) ListTask(ginc *gin.Context) {

	if lock, err := ctrl.cacheMgr.Lock(ginc, c.LockKey, 60); err != nil || lock != true {
		ctrl.handleError(ginc, err, http.StatusLocked, code.Code_INTERNAL)
		return
	}
	defer ctrl.cacheMgr.ReleaseLock(ginc, c.LockKey)

	limitStr := ginc.Query("limit")
	offsetStr := ginc.Query("offset")
	order := ginc.Query("order")

	defaultLimit := 20
	defaultOffset := 0

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = defaultLimit
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = defaultOffset
	}

	if order == "" {
		order = "id"
	}

	if ctrl.enableCache {
		tasks, err := ctrl.cacheMgr.ListTask(ginc, limit, offset, order)
		if err != nil {
			logger.GetLoggerWithKeys(map[string]interface{}{
				"error": err,
			}).Error("ListTask fail")
			ctrl.handleError(ginc, err, http.StatusLocked, code.Code_INTERNAL)
			return
		}

		if len(tasks) != 0 {
			ginc.JSON(http.StatusOK, models.Response{
				Code:    code.Code_OK,
				Message: c.Success,
				Data:    tasks,
			})
			return
		}
	}

	tasks, err := ctrl.mysqlMgr.ListTask(ginc, limit, offset, order)
	if err != nil {
		logger.GetLoggerWithKeys(map[string]interface{}{
			"error": err,
		}).Error("ListTask fail")
		ctrl.handleError(ginc, err, http.StatusLocked, code.Code_INTERNAL)
		return
	}

	ginc.JSON(http.StatusOK, models.Response{
		Code:    code.Code_OK,
		Message: c.Success,
		Data:    tasks,
	})
}
func (ctrl *Controller) GetTask(ginc *gin.Context) {

	if lock, err := ctrl.cacheMgr.Lock(ginc, c.LockKey, 60); err != nil || lock != true {
		ctrl.handleError(ginc, err, http.StatusLocked, code.Code_INTERNAL)
		return
	}
	defer ctrl.cacheMgr.ReleaseLock(ginc, c.LockKey)

	taskIdStr := ginc.Param("taskId")

	taskId, err := strconv.ParseUint(taskIdStr, 10, 64)
	if err != nil {
		logger.GetLoggerWithKeys(map[string]interface{}{
			"error": err,
		}).Error("GetTask fail")
		ctrl.handleError(ginc, err, http.StatusLocked, code.Code_INTERNAL)
		return
	}

	if ctrl.enableCache {
		task, err := ctrl.cacheMgr.GetTaskById(ginc, taskId)
		if err != nil {
			logger.GetLoggerWithKeys(map[string]interface{}{
				"error": err,
			}).Error("GetTask fail")
			ctrl.handleError(ginc, err, http.StatusLocked, code.Code_INTERNAL)
			return
		}

		if task != nil {
			ginc.JSON(http.StatusOK, models.Response{
				Code:    code.Code_OK,
				Message: c.Success,
				Data:    task,
			})
			return
		}
	}

	task, err := ctrl.mysqlMgr.GetTaskById(ginc, taskId)
	if err != nil {
		logger.GetLoggerWithKeys(map[string]interface{}{
			"error": err,
		}).Error("GetTask fail")
		ctrl.handleError(ginc, err, http.StatusLocked, code.Code_INTERNAL)
		return
	}

	if err := ctrl.cacheMgr.CreateTask(ginc, task); err != nil {
		logger.GetLoggerWithKeys(map[string]interface{}{
			"error": err,
		}).Error("insert task into cache fail")
	}

	ginc.JSON(http.StatusOK, models.Response{
		Code:    code.Code_OK,
		Message: c.Success,
		Data:    task,
	})

}

func (ctrl *Controller) CreateTask(ginc *gin.Context) {

	if lock, err := ctrl.cacheMgr.Lock(ginc, c.LockKey, 60); err != nil || lock != true {
		ctrl.handleError(ginc, err, http.StatusLocked, code.Code_INTERNAL)
		return
	}
	defer ctrl.cacheMgr.ReleaseLock(ginc, c.LockKey)

	task := models.Task{}
	if err := ginc.BindJSON(&task); err != nil {
		ctrl.handleError(ginc, err, http.StatusBadRequest, code.Code_INTERNAL)
		return
	}

	condition := map[string]interface{}{
		"name": task.Name,
	}
	if task.Tag != nil {
		condition["tag"] = task.Tag
	}

	if err := ctrl.mysqlMgr.CheckTaskExist(ginc, condition, &task); err != nil {
		logger.GetLoggerWithKeys(map[string]interface{}{
			"error": err,
		}).Error("ListTask fail")
		ctrl.handleError(ginc, err, http.StatusBadRequest, code.Code_INTERNAL)
		return
	}

	if task.ID != 0 {
		ctrl.handleError(ginc, fmt.Errorf("task is exist"), http.StatusBadRequest, code.Code_INVALID_ARGUMENT)
		return
	}

	if err := ctrl.mysqlMgr.CreateTask(ginc, &task); err != nil {
		logger.GetLoggerWithKeys(map[string]interface{}{
			"error": err,
		}).Error("ListTask fail")
		ctrl.handleError(ginc, err, http.StatusInternalServerError, code.Code_INTERNAL)
		return
	}

	targetTask, _ := ctrl.mysqlMgr.GetTaskById(ginc, task.ID)

	if ctrl.enableCache {
		if err := ctrl.cacheMgr.CreateTask(ginc, targetTask); err != nil {
			logger.GetLoggerWithKeys(map[string]interface{}{
				"error": err,
			}).Error("insert task into cache fail")
		} else {
			ctrl.checkTaskVersion(ginc, task.ID, task.Version)
		}
	}

	ginc.JSON(http.StatusOK, models.Response{
		Code:    code.Code_OK,
		Message: c.Success,
		Data:    targetTask,
	})
}

func (ctrl *Controller) DeleteTask(ginc *gin.Context) {

	if lock, err := ctrl.cacheMgr.Lock(ginc, c.LockKey, 60); err != nil || lock != true {
		ctrl.handleError(ginc, err, http.StatusLocked, code.Code_INTERNAL)
		return
	}
	defer ctrl.cacheMgr.ReleaseLock(ginc, c.LockKey)

	taskIdStr := ginc.Param("taskId")

	taskId, err := strconv.ParseUint(taskIdStr, 10, 64)
	if err != nil {
		ctrl.handleError(ginc, err, http.StatusBadRequest, code.Code_INTERNAL)
		return
	}

	if ctrl.enableCache {
		if err := ctrl.cacheMgr.DeleteTask(ginc, taskId); err != nil {
			ctrl.enableCache = false
		}
	}

	if err := ctrl.mysqlMgr.DeleteTask(ginc, taskId); err != nil {
		ctrl.handleError(ginc, err, http.StatusInternalServerError, code.Code_INTERNAL)
		return
	}

	ginc.JSON(http.StatusOK, models.Response{
		Code:    code.Code_OK,
		Message: c.Success,
	})
}

func (ctrl *Controller) UpdateTask(ginc *gin.Context) {

	if lock, err := ctrl.cacheMgr.Lock(ginc, c.LockKey, 60); err != nil || lock != true {
		ctrl.handleError(ginc, err, http.StatusLocked, code.Code_INTERNAL)
		return
	}
	defer ctrl.cacheMgr.ReleaseLock(ginc, c.LockKey)

	taskIdStr := ginc.Param("taskId")

	taskId, err := strconv.ParseUint(taskIdStr, 10, 64)
	if err != nil {
		ctrl.handleError(ginc, err, http.StatusBadRequest, code.Code_INTERNAL)
		return
	}

	task := models.Task{
		ID: taskId,
	}
	if err := ginc.BindJSON(&task); err != nil {
		ctrl.handleError(ginc, err, http.StatusBadRequest, code.Code_INTERNAL)
		return
	}
	targetTask, err := ctrl.mysqlMgr.GetTaskById(ginc, taskId)
	if err != nil {
		ctrl.handleError(ginc, err, http.StatusBadRequest, code.Code_INTERNAL)
		return
	}

	targetTask.Name = task.Name
	targetTask.Content = task.Content
	targetTask.Tag = task.Tag

	if ctrl.enableCache {
		if err := ctrl.cacheMgr.UpdateTask(ginc, targetTask); err != nil {
			ctrl.enableCache = false
		}
		ginc.JSON(http.StatusOK, models.Response{
			Code:    code.Code_OK,
			Message: c.Success,
			Data:    task,
		})
		return
	}

	if err := ctrl.mysqlMgr.UpdateTask(ginc, targetTask); err != nil {
		ctrl.enableCache = false
	}
	ginc.JSON(http.StatusOK, models.Response{
		Code:    code.Code_OK,
		Message: c.Success,
		Data:    targetTask,
	})

}

func (ctrl *Controller) Shutdown() {}

func (ctrl *Controller) handleError(ginc *gin.Context, err error, httpCode int, errorCode code.Code) {
	ginc.JSON(httpCode, models.HttpError{
		Code:    errorCode,
		Message: err.Error(),
	})
}

func (ctrl *Controller) checkTaskVersion(ctx context.Context, taskId uint64, version int) {
	task, err := ctrl.cacheMgr.GetTaskById(ctx, taskId)
	if err == nil && version == task.Version {
		return
	}

	logger.GetLoggerWithKeys(map[string]interface{}{
		"error":  err,
		"taskId": taskId,
	}).Error("checkTaskVersion: version not match")

	if err := ctrl.cacheMgr.DeleteTask(ctx, taskId); err != nil {
		logger.GetLoggerWithKeys(map[string]interface{}{
			"error":  err,
			"taskId": taskId,
		}).Error("checkTaskVersion: delete cache task fail")
		ctrl.enableCache = false
	}

}
