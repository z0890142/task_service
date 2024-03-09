package data

import (
	"context"
	"task_service/pkg/models"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type DataManager interface {
	ListTask(ctx context.Context, limit, offset int, order string) ([]models.Task, error)
	GetTaskById(ctx context.Context, taskId uint64) (models.Task, error)
	CheckTaskExist(ctx context.Context, condition map[string]interface{}, task *models.Task) error
	CreateTask(ctx context.Context, task []models.Task) error
	DeleteTask(ctx context.Context, taskId uint64) error
	UpdateTask(ctx context.Context, task *models.Task) error

	Lock(ctx context.Context, lockKey string, expiration time.Duration) (bool, error)
	ReleaseLock(ctx context.Context, lockKey string)
	Close(context.Context)
}

func NewDataManager(client interface{}) DataManager {
	switch client.(type) {
	case *gorm.DB:
		return newMysqlManager(client.(*gorm.DB))
	case *redis.Client:
		return newCacheMgr(client.(*redis.Client))
	}
	return nil
}
