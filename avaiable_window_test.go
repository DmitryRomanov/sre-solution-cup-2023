package main

import (
	"testing"
	"time"

	"github.com/DmitryRomanov/sre-solution-cup-2023/models"
	"github.com/stretchr/testify/assert"
)

func TestGetAvaiableWindow1(t *testing.T) {
	assert := assert.New(t)

	windows := []models.MaintenanceWindows{
		{Start: 0, End: 8},
	}

	now := time.Now()
	tasks := []models.Task{
		{
			StartTime:  time.Date(now.Year(), now.Month(), now.Day(), 1, 30, 0, now.Nanosecond(), now.Location()),
			FinishTime: time.Date(now.Year(), now.Month(), now.Day(), 2, 0, 0, now.Nanosecond(), now.Location()),
		},
		{
			StartTime:  time.Date(now.Year(), now.Month(), now.Day(), 2, 5, 0, now.Nanosecond(), now.Location()),
			FinishTime: time.Date(now.Year(), now.Month(), now.Day(), 2, 30, 0, now.Nanosecond(), now.Location()),
		},
		{
			StartTime:  time.Date(now.Year(), now.Month(), now.Day(), 5, 0, 0, now.Nanosecond(), now.Location()),
			FinishTime: time.Date(now.Year(), now.Month(), now.Day(), 6, 0, 0, now.Nanosecond(), now.Location()),
		},
	}

	task := new(models.Task)
	task.StartTime = time.Date(now.Year(), now.Month(), now.Day(), 0, 10, 0, now.Nanosecond(), now.Location())
	task.Duration = 2000

	result := getAvaiableWindows(task, windows, tasks)
	assert.Equal(
		time.Date(now.Year(), now.Month(), now.Day(), 0, 00, 0, now.Nanosecond(), now.Location()),
		result[0],
	)
	assert.Equal(
		time.Date(now.Year(), now.Month(), now.Day(), 2, 30, 1, now.Nanosecond(), now.Location()),
		result[1],
	)
	assert.Equal(
		time.Date(now.Year(), now.Month(), now.Day(), 6, 00, 1, now.Nanosecond(), now.Location()),
		result[2],
	)
	assert.Equal(3, len(result))
}
