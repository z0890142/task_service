package data

import (
	"context"
	"fmt"
	"strings"
	"task_service/pkg/logger"
	"task_service/pkg/models"
	"task_service/pkg/utils"
	"time"

	"github.com/redis/go-redis/v9"
)

type CacheMgr struct {
	client *redis.Client
}

func newCacheMgr(client *redis.Client) *CacheMgr {
	return &CacheMgr{
		client: client,
	}
}

func (mgr *CacheMgr) ListTask(ctx context.Context, limit, offset int, order string) ([]models.Task, error) {
	orderSlice := strings.Split(order, " ")
	desc := false
	if len(orderSlice) == 2 && orderSlice[1] == "desc" {
		desc = true
	}

	keys, err := mgr.client.Keys(ctx, "task:*").Result()
	if err != nil {
		return nil, fmt.Errorf("ListTask:%v", err)
	}

	var tasks []models.Task
	tx := mgr.client.TxPipeline()

	for _, key := range keys {
		tx.HGetAll(ctx, key).Result()
	}

	cmds, err := tx.Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("ListTask: %v", err)
	}

	var results []map[string]string
	for _, cmd := range cmds {
		result, err := cmd.(*redis.MapStringStringCmd).Result()
		if err != nil {
			return nil, fmt.Errorf("ListTask: %v", err)
		}
		results = append(results, result)
	}

	for _, result := range results {
		task, err := utils.ConvertTask(result)
		if err != nil {
			return nil, fmt.Errorf("ListTask: %v", err)
		}
		tasks = append(tasks, task)
	}

	// 根據 orderSlice 中的排序條件對 tasks 進行排序
	utils.SortByField(tasks, orderSlice[0], desc)

	// 根據 limit 和 offset 截取 tasks 切片
	if offset > len(tasks) {
		offset = len(tasks)
	}
	if limit+offset > len(tasks) {
		limit = len(tasks) - offset
	}
	return tasks[offset : offset+limit], nil
}

func (mgr *CacheMgr) GetTaskById(ctx context.Context, taskId uint64) (models.Task, error) {
	key := getKey(taskId)

	result, err := mgr.client.HGetAll(ctx, key).Result()
	if err != nil {
		return models.Task{}, fmt.Errorf("GetTaskById: %v", err)
	}

	task, err := utils.ConvertTask(result)
	return task, nil
}

func (mgr *CacheMgr) CheckTaskExist(ctx context.Context, condition map[string]interface{}, task *models.Task) error {
	return nil
}

func (mgr *CacheMgr) CreateTask(ctx context.Context, tasks []models.Task) error {
	tx := mgr.client.TxPipeline()
	for _, task := range tasks {
		key := getKey(task.ID)

		tx.HSet(ctx, key, "id", task.ID)
		tx.HSet(ctx, key, "name", task.Name)
		tx.HSet(ctx, key, "content", task.Content)
		tx.HSet(ctx, key, "tag", task.Tag)
		tx.HSet(ctx, key, "status", task.Status)
		tx.HSet(ctx, key, "version", task.Version)
		tx.HSet(ctx, key, "created_at", task.CreatedAt)
		tx.HSet(ctx, key, "updated_at", task.UpdatedAt)
	}

	if _, err := tx.Exec(ctx); err != nil {
		return fmt.Errorf("UpdateTask: %v", err)
	}

	return nil
}

func (mgr *CacheMgr) DeleteTask(ctx context.Context, taskId uint64) error {
	key := getKey(taskId)

	keys, err := mgr.client.HKeys(ctx, key).Result()
	if err != nil {
		return fmt.Errorf("DeleteTask:%v", err)
	}
	for _, key := range keys {
		if err := mgr.client.HDel(ctx, key).Err(); err != nil {
			return fmt.Errorf("DeleteTask:%v", err)
		}
	}

	return nil
}

func (mgr *CacheMgr) UpdateTask(ctx context.Context, task *models.Task) error {
	key := getKey(task.ID)

	tx := mgr.client.TxPipeline()
	tx.HSet(ctx, key, "id", task.ID)
	tx.HSet(ctx, key, "name", task.Name)
	tx.HSet(ctx, key, "content", task.Content)
	tx.HSet(ctx, key, "tag", task.Tag)
	tx.HSet(ctx, key, "status", task.Status)
	tx.HSet(ctx, key, "version", task.Version)
	tx.HSet(ctx, key, "created_at", task.CreatedAt)
	tx.HSet(ctx, key, "updated_at", task.UpdatedAt)

	if _, err := tx.Exec(ctx); err != nil {
		return fmt.Errorf("UpdateTask: %v", err)
	}

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

func getKey(taskId uint64) string {
	return fmt.Sprintf("task:%d", taskId)
}
