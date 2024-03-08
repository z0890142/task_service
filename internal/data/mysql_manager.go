package data

import (
	"context"
	"fmt"
	"task_service/pkg/logger"
	"task_service/pkg/models"
	"time"

	"gorm.io/gorm"
)

type MysqlMgr struct {
	client *gorm.DB
}

func NewMysqlManager(gormClient *gorm.DB) DataManager {
	return &MysqlMgr{
		client: gormClient,
	}
}
func (mgr *MysqlMgr) ListTask(ctx context.Context, limit, offset int, order string) ([]models.Task, error) {
	var tasks []models.Task
	if err := mgr.client.
		Order(order).
		Offset(offset).Limit(limit).
		Find(&tasks).
		Error; err != nil {
		return nil, fmt.Errorf("ListTask: %s", err.Error())
	}
	return tasks, nil
}

func (mgr *MysqlMgr) GetTaskById(ctx context.Context, taskId uint64) (*models.Task, error) {
	task := models.Task{
		ID: taskId,
	}
	if err := mgr.client.First(&task).Error; err != nil {
		return nil, fmt.Errorf("GetTaskById: %s", err.Error())
	}
	return &task, nil
}

func (mgr *MysqlMgr) CheckTaskExist(ctx context.Context, condition map[string]interface{}, task *models.Task) error {
	if err := mgr.client.Where(condition).First(task).Error; err != nil && err != gorm.ErrRecordNotFound {
		return fmt.Errorf("CheckTaskExist: %s", err.Error())
	}

	return nil
}

func (mgr *MysqlMgr) CreateTask(ctx context.Context, task *models.Task) error {
	if err := mgr.client.Create(task).Error; err != nil {
		return fmt.Errorf("CreateTask: %s", err.Error())
	}
	return nil
}

func (mgr *MysqlMgr) DeleteTask(ctx context.Context, taskId uint64) error {
	if err := mgr.client.Delete(&models.Task{}, "id = ?", taskId).Error; err != nil {
		return fmt.Errorf("DeleteTask: %s", err.Error())
	}
	return nil
}

func (mgr *MysqlMgr) UpdateTask(ctx context.Context, task *models.Task) error {
	if err := mgr.client.Where("id=?", task.ID).Save(task).Error; err != nil {
		return fmt.Errorf("UpdateTask: %s", err.Error())
	}
	return nil
}

func (mgr *MysqlMgr) Lock(ctx context.Context, lockKey string, expiration time.Duration) (bool, error) {
	return false, nil
}
func (mgr *MysqlMgr) ReleaseLock(ctx context.Context, lockKey string) {

}

func (mgr *MysqlMgr) Close(ctx context.Context) {
	db, err := mgr.client.DB()
	if err != nil {
		logger.Errorf("Close DB : %s", err.Error())
	}
	if err = db.Close(); err != nil {
		logger.Errorf("Close DB : %s", err.Error())
	}

}
