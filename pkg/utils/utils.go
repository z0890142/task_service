package utils

import (
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"task_service/pkg/models"
	"time"
)

var fildMap = map[string]string{
	"id":         "ID",
	"name":       "Name",
	"status":     "Stauts",
	"content":    "Content",
	"tag":        "Tag",
	"created_at": "CreatedAt",
	"updated_at": "UpdatedAt",
}

func ConvertTask(result map[string]string) (models.Task, error) {
	task := models.Task{}
	for field, value := range result {
		switch field {
		case "id":
			id, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				return task, fmt.Errorf("ConvertTask: %v", err)
			}
			task.ID = id
		case "name":
			task.Name = value
		case "status":
			status, err := strconv.Atoi(value)
			if err != nil {
				return task, fmt.Errorf("ConvertTask: %v", err)
			}
			task.Status = status
		case "content":
			task.Content = value
		case "tag":
			task.Tag = value
		case "version":
			version, err := strconv.Atoi(value)
			if err != nil {
				return task, fmt.Errorf("ConvertTask: %v", err)
			}
			task.Version = version
		case "created_at":
			createdAt, err := time.Parse(time.RFC3339, value)
			if err != nil {
				return task, fmt.Errorf("ConvertTask: %v", err)
			}
			task.CreatedAt = createdAt
		case "updated_at":
			updatedAt, err := time.Parse(time.RFC3339, value)
			if err != nil {
				return task, fmt.Errorf("ConvertTask: %v", err)
			}
			task.UpdatedAt = updatedAt
		}
	}
	return task, nil
}

func SortByField(tasks []models.Task, fieldName string, desc bool) {
	if _, ok := fildMap[fieldName]; !ok {
		return
	}

	if desc {
		sort.Slice(tasks, func(i, j int) bool {
			fieldI := reflect.ValueOf(tasks[i]).FieldByName(fildMap[fieldName])
			fieldJ := reflect.ValueOf(tasks[j]).FieldByName(fildMap[fieldName])

			switch fieldI.Interface().(type) {
			case string:
				valI := fieldI.Interface().(string)
				valJ := fieldJ.Interface().(string)
				if valI == "" && valJ != "" {
					return false // 如果 fieldI 為空字串，則將其排在後面
				} else if valI != "" && valJ == "" {
					return true // 如果 fieldJ 為空字串，則將其排在前面
				} else if valI == "" && valJ == "" {
					return false // 如果兩者都為空字串，則保持原有順序
				}
				return fieldI.Interface().(string) > fieldJ.Interface().(string)

			case int:
				return fieldI.Interface().(int) > fieldJ.Interface().(int)
			case uint64:
				return fieldI.Interface().(uint64) > fieldJ.Interface().(uint64)
			default:
				// Handle other types if needed
				return false
			}
		})
	} else {
		sort.Slice(tasks, func(i, j int) bool {
			fieldI := reflect.ValueOf(tasks[i]).FieldByName(fildMap[fieldName])
			fieldJ := reflect.ValueOf(tasks[j]).FieldByName(fildMap[fieldName])

			switch fieldI.Interface().(type) {
			case string:
				valI := fieldI.Interface().(string)
				valJ := fieldJ.Interface().(string)
				if valI == "" && valJ != "" {
					return true // 如果 fieldI 為空字串，則將其排在後面
				} else if valI != "" && valJ == "" {
					return false // 如果 fieldJ 為空字串，則將其排在前面
				} else if valI == "" && valJ == "" {
					return false // 如果兩者都為空字串，則保持原有順序
				}
				return fieldI.Interface().(string) < fieldJ.Interface().(string)

			case int:
				return fieldI.Interface().(int) < fieldJ.Interface().(int)
			case uint64:
				return fieldI.Interface().(uint64) < fieldJ.Interface().(uint64)
			default:
				// Handle other types if needed
				return false
			}
		})
	}

}
