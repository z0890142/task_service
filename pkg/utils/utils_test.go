package utils

import (
	"task_service/pkg/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConvertTask(t *testing.T) {
	tests := []struct {
		result   map[string]string
		isErr    bool
		expected models.Task
	}{
		{
			map[string]string{
				"id":         "1",
				"name":       "test task",
				"status":     "1",
				"content":    "test content",
				"version":    "1",
				"Tag":        "123",
				"created_at": "2006-01-02T15:04:05+08:00",
				"updated_at": "2006-01-02T15:04:05+08:00",
			},
			false,
			models.Task{
				ID:      1,
				Name:    "test task",
				Status:  1,
				Content: "test content",

				Version: 1,
			},
		},
		{
			map[string]string{
				"id":         "1",
				"name":       "test task",
				"status":     "1",
				"content":    "test content",
				"version":    "1",
				"created_at": "2006-01-02T15:04:05",
				"updated_at": "2006-01-02T15:04:05",
			},
			true,
			models.Task{},
		},
	}

	for _, testItem := range tests {
		task, err := ConvertTask(testItem.result)
		if testItem.isErr {
			assert.NotNil(t, err)
			continue
		}
		createdAt, err := time.Parse(time.RFC3339, testItem.result["created_at"])
		assert.Nil(t, err)

		updatedAt, err := time.Parse(time.RFC3339, testItem.result["updated_at"])
		assert.Nil(t, err)

		testItem.expected.CreatedAt = createdAt
		testItem.expected.UpdatedAt = updatedAt
		testItem.expected.Tag = testItem.result["tag"]

		assert.Nil(t, err)
		assert.Equal(t, testItem.expected, task)
	}

}

func TestSortByField(t *testing.T) {

	createdAt1, _ := time.Parse(time.RFC3339, "2006-01-02T15:04:05+08:00")
	updateAt1, _ := time.Parse(time.RFC3339, "2006-01-02T15:04:05+08:00")

	createdAt2, _ := time.Parse(time.RFC3339, "2006-01-02T16:04:05+08:00")
	updateAt2, _ := time.Parse(time.RFC3339, "2006-01-02T16:04:05+08:00")

	createdAt3, _ := time.Parse(time.RFC3339, "2006-01-02T17:04:05+08:00")
	updateAt3, _ := time.Parse(time.RFC3339, "2006-01-02T17:04:05+08:00")

	tasks := []models.Task{
		models.Task{
			ID:        1,
			Name:      "b",
			Content:   "content1",
			Tag:       "tag2",
			CreatedAt: createdAt1,
			UpdatedAt: updateAt1,
		},
		models.Task{
			ID:        2,
			Name:      "a",
			Content:   "content2",
			Tag:       "tag1",
			CreatedAt: createdAt2,
			UpdatedAt: updateAt2,
		},
		models.Task{
			ID:        3,
			Name:      "c",
			Content:   "content1",
			Tag:       "",
			CreatedAt: createdAt3,
			UpdatedAt: updateAt3,
		},
	}

	tests := []struct {
		col      string
		desc     bool
		expected []models.Task
	}{
		{
			"",
			true,
			tasks,
		},
		{
			"id",
			true,
			[]models.Task{
				models.Task{
					ID:        3,
					Name:      "c",
					Content:   "content1",
					Tag:       "",
					CreatedAt: createdAt3,
					UpdatedAt: updateAt3,
				},
				models.Task{
					ID:        2,
					Name:      "a",
					Content:   "content2",
					Tag:       "tag1",
					CreatedAt: createdAt2,
					UpdatedAt: updateAt2,
				},
				models.Task{
					ID:        1,
					Name:      "b",
					Content:   "content1",
					Tag:       "tag2",
					CreatedAt: createdAt1,
					UpdatedAt: updateAt1,
				},
			},
		},
		{
			"id",
			false,
			[]models.Task{
				models.Task{
					ID:        1,
					Name:      "b",
					Content:   "content1",
					Tag:       "tag2",
					CreatedAt: createdAt1,
					UpdatedAt: updateAt1,
				},
				models.Task{
					ID:        2,
					Name:      "a",
					Content:   "content2",
					Tag:       "tag1",
					CreatedAt: createdAt2,
					UpdatedAt: updateAt2,
				},
				models.Task{
					ID:        3,
					Name:      "c",
					Content:   "content1",
					Tag:       "",
					CreatedAt: createdAt3,
					UpdatedAt: updateAt3,
				},
			},
		},
		{
			"name",
			false,
			[]models.Task{
				models.Task{
					ID:        2,
					Name:      "a",
					Content:   "content2",
					Tag:       "tag1",
					CreatedAt: createdAt2,
					UpdatedAt: updateAt2,
				},
				models.Task{
					ID:        1,
					Name:      "b",
					Content:   "content1",
					Tag:       "tag2",
					CreatedAt: createdAt1,
					UpdatedAt: updateAt1,
				},
				models.Task{
					ID:        3,
					Name:      "c",
					Content:   "content1",
					Tag:       "",
					CreatedAt: createdAt3,
					UpdatedAt: updateAt3,
				},
			},
		},
		{
			"tag",
			true,
			[]models.Task{
				models.Task{
					ID:        1,
					Name:      "b",
					Content:   "content1",
					Tag:       "tag2",
					CreatedAt: createdAt1,
					UpdatedAt: updateAt1,
				},
				models.Task{
					ID:        2,
					Name:      "a",
					Content:   "content2",
					Tag:       "tag1",
					CreatedAt: createdAt2,
					UpdatedAt: updateAt2,
				},
				models.Task{
					ID:        3,
					Name:      "c",
					Content:   "content1",
					Tag:       "",
					CreatedAt: createdAt3,
					UpdatedAt: updateAt3,
				},
			},
		},
		{
			"tag",
			false,
			[]models.Task{

				models.Task{
					ID:        3,
					Name:      "c",
					Content:   "content1",
					Tag:       "",
					CreatedAt: createdAt3,
					UpdatedAt: updateAt3,
				},
				models.Task{
					ID:        2,
					Name:      "a",
					Content:   "content2",
					Tag:       "tag1",
					CreatedAt: createdAt2,
					UpdatedAt: updateAt2,
				},
				models.Task{
					ID:        1,
					Name:      "b",
					Content:   "content1",
					Tag:       "tag2",
					CreatedAt: createdAt1,
					UpdatedAt: updateAt1,
				},
			},
		},
	}

	for _, testItem := range tests {
		SortByField(tasks, testItem.col, testItem.desc)
		assert.Equal(t, testItem.expected, tasks)
	}

}
