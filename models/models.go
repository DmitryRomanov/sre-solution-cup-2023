package models

import (
	"time"

	"gorm.io/gorm"
)

type Task struct {
	gorm.Model
	Name           string       `json:"name"`
	AviabilityZone string       `json:"aviability_zone"`
	Type           TaskType     `json:"type"`
	Priority       TaskPriority `json:"priority"`
	StartTime      time.Time
	Duration       int
	Deadline       time.Time
}

type AviabilityZone struct {
	ID                      uint   `gorm:"primarykey" json:"id"`
	Name                    string `json:"name" gorm:"unique"`
	DataCenter              string `json:"data_center"`
	BlockedForAutomatedTask bool   `json:"blocked_for_auto_tasks"`
}
