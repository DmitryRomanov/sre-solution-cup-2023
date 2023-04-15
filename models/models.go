package models

import (
	"errors"
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

func (task *Task) ValidateValues() error {
	if task.Priority == string(TASK_PRIORITY_CRITICAL) && task.Type == string(TASK_TYPE_MANUAL) {
		return errors.New("приоритет critical (может выставляться только для типа manual)")
	}

	return nil
}

type AviabilityZone struct {
	Name                    string `json:"name" gorm:"unique"`
	DataCenter              string `json:"data_center"`
	BlockedForAutomatedTask bool   `json:"blocked_for_auto_tasks"`
}
