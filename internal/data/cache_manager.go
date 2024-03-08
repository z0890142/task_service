package data

import (
	"context"
	"fmt"
	"task_service/pkg/logger"
	"task_service/pkg/models"
	"time"

	"github.com/redis/go-redis/v9"
)

type CacheMgr struct {
	client *redis.Client
}

func NewCacheMgr(client *redis.Client) *CacheMgr {
	return &CacheMgr{
		client: client,
	}
}

func (mgr *CacheMgr) ListTask(ctx context.Context, limit, offset int, order string) ([]models.Task, error) {
	return nil, nil
}
func (mgr *CacheMgr) GetTaskById(ctx context.Context, taskId uint64) (*models.Task, error) {
	return nil, nil
}
func (mgr *CacheMgr) CheckTaskExist(ctx context.Context, condition map[string]interface{}, task *models.Task) error {
	return nil
}
func (mgr *CacheMgr) CreateTask(ctx context.Context, task *models.Task) error {
	return nil
}
func (mgr *CacheMgr) DeleteTask(ctx context.Context, taskId uint64) error {
	return nil
}
func (mgr *CacheMgr) UpdateTask(ctx context.Context, task *models.Task) error {
	return nil
}

func (mgr *CacheMgr) Lock(ctx context.Context, lockKey string, expiration time.Duration) (bool, error) {
	success, err := mgr.client.SetNX(ctx, lockKey, "locked", expiration).Result()
	if err != nil {
		return false, fmt.Errorf("LockTask: %v", err)
	}
	return success, nil
}

func (mgr *CacheMgr) ReleaseLock(ctx context.Context, lockKey string) {
	if _, err := mgr.client.Del(ctx, lockKey).Result(); err != nil {
		logger.GetLoggerWithKeys(map[string]interface{}{
			"error": err,
		}).Error("ReleaseLock Fail")
	}
}

func (mgr *CacheMgr) Close(ctx context.Context) {
	if err := mgr.client.Close(); err != nil {
		logger.Errorf("Close Redis : %s", err.Error())
	}
}
