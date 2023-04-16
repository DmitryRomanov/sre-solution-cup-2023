package main

import (
	"fmt"
	"time"

	"github.com/DmitryRomanov/sre-solution-cup-2023/models"
)

func cancelAutoTasks(newTask *models.Task) {
	duration := time.Duration(newTask.Duration-1) * time.Second
	finishTime := newTask.StartTime.Add(duration)

	var tasks []models.Task

	db.Debug().Where(
		"aviability_zone = ? AND type = ? AND status = ? AND ((? BETWEEN start_time AND finish_time) OR (? BETWEEN start_time AND finish_time))",
		newTask.AviabilityZone,
		models.TASK_TYPE_AUTO,
		models.TASK_STATUS_WAITING,
		newTask.StartTime,
		finishTime,
	).Find(&tasks)
	for i := range tasks {
		cancelTask(tasks[i], fmt.Sprintf("cancelAutoTasks by task %v", newTask))
	}
}

func cancelManualTasksWithNormalPriority(newTask *models.Task) {
	duration := time.Duration(newTask.Duration-1) * time.Second
	finishTime := newTask.StartTime.Add(duration)

	var tasks []models.Task

	db.Debug().Where(
		"aviability_zone = ? AND type = ? AND priority = ? AND status = ? AND ((? BETWEEN start_time AND finish_time) OR (? BETWEEN start_time AND finish_time))",
		newTask.AviabilityZone,
		models.TASK_TYPE_MANUAL,
		models.TASK_PRIORITY_NORMAL,
		models.TASK_STATUS_WAITING,
		newTask.StartTime,
		finishTime,
	).Find(&tasks)
	for i := range tasks {
		cancelTask(tasks[i], fmt.Sprintf("cancelManualTasksWithNormalPriority by task %v", newTask))
	}
}

func cancelTask(task models.Task, reason string) {
	cancelReason := new(models.CancelReason)
	cancelReason.CancelTime = time.Now()
	cancelReason.Reason = reason
	cancelReason.TaskID = task.ID
	db.Create(cancelReason)

	db.Model(models.Task{}).Where("id = ?", task.ID).Update("status", models.TASK_STATUS_CANCELED)
}
