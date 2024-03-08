package models

import (
	"time"
)

type Task struct {
	ID        uint64    `json:"id" gorm:"primaryKey;autoIncrement"`
	Name      string    `json:"name" gorm:"size:200;not null"`
	Status    int       `json:"status" gorm:"type:tinyint;not null;default:1"`
	Content   string    `json:"content" gorm:"size:500;not null"`
	Tag       *string   `json:"tag" gorm:"size:50;null"`
	Version   int       `json:"-,omitempty"; gorm:"version:int;null"`
	CreatedAt time.Time `json:"created_at,omitempty" gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `json:"updated_at,omitempty" gorm:"not null;default:CURRENT_TIMESTAMP"`
}

func (Task) TableName() string {
	return "Task"
}
