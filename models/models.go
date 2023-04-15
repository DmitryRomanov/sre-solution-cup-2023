package models

import (
	"time"
)

type Task struct {
	ID             uint   `gorm:"primarykey"`
	Name           string `json:"name"`
	AviabilityZone string `json:"aviability_zone"`
	Type           string `json:"type"`
	Priority       string `json:"priority"`
	StartTime      time.Time
	Duration       int
	Deadline       time.Time
}

type AviabilityZone struct {
	Name                    string `json:"name" gorm:"unique"`
	DataCenter              string `json:"data_center"`
	BlockedForAutomatedTask bool   `json:"blocked_for_auto_tasks"`
}
