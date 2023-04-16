package models

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Task struct {
	ID             uint   `gorm:"primarykey"`
	Status         string `json:"status"`
	AviabilityZone string `json:"aviability_zone"`
	Type           string `json:"type"`
	Priority       string `json:"priority"`
	StartTime      time.Time
	FinishTime     time.Time
	Duration       int
	Deadline       time.Time
}

func (task *Task) BeforeCreate(tx *gorm.DB) (err error) {
	duration := time.Duration(task.Duration-1) * time.Second
	finishTime := task.StartTime.Add(duration)
	task.FinishTime = finishTime
	return
}

func (task *Task) ValidateValues() error {
	if task.Priority == TASK_PRIORITY_CRITICAL && task.Type != TASK_TYPE_MANUAL {
		return errors.New("приоритет critical (может выставляться только для типа manual)")
	}

	manualTaskMinuteMultiplicity := 5
	if task.Type == TASK_TYPE_MANUAL &&
		task.StartTime.Minute()%manualTaskMinuteMultiplicity != 0 &&
		task.StartTime.Second() == 0 {
		return fmt.Errorf("для ручных работ время кратно %v минутам", manualTaskMinuteMultiplicity)
	}

	autoTaskMinuteMultiplicity := 1
	if task.Type == TASK_TYPE_AUTO &&
		task.StartTime.Minute()%autoTaskMinuteMultiplicity != 0 &&
		task.StartTime.Second() == 0 {
		return fmt.Errorf("для автоматических работ время кратно %v минутам", autoTaskMinuteMultiplicity)
	}

	manualTaskMinDurationSeconds := 30 * 60
	if task.Type == TASK_TYPE_MANUAL && task.Duration < manualTaskMinDurationSeconds {
		return fmt.Errorf("для ручных работ длительность не меньше %v секунд", manualTaskMinDurationSeconds)
	}

	autoTaskMinDurationSeconds := 5 * 60
	if task.Type == TASK_TYPE_MANUAL && task.Duration < autoTaskMinDurationSeconds {
		return fmt.Errorf("для автоматических работ длительность не меньше %v секунд", autoTaskMinDurationSeconds)
	}

	criticalTaskMaxDurationSeconds := 6 * 60 * 60
	if task.Priority == TASK_PRIORITY_NORMAL && task.Duration > criticalTaskMaxDurationSeconds {
		return fmt.Errorf("длительность не больше %v секунд", criticalTaskMaxDurationSeconds)
	}

	return nil
}

type AviabilityZone struct {
	Name                    string `json:"name" gorm:"unique"`
	DataCenter              string `json:"data_center"`
	BlockedForAutomatedTask bool   `json:"blocked_for_auto_tasks"`
}

type CancelReason struct {
	TaskID     uint      `json:"task_id" gorm:"unique"`
	CancelTime time.Time `json:"cancel_time"`
	Reason     string    `json:"reason"`
}

type MaintenanceWindows struct {
	Start          int
	End            int
	AviabilityZone string `json:"aviability_zone"`
}
