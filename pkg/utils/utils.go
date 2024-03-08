package utils

import "fmt"

func GetRedisKey(taskId int) string {
	return fmt.Sprintf("task:%v", taskId)
}
