package main

import (
	"testing"
	"time"

	"github.com/DmitryRomanov/sre-solution-cup-2023/models"
	"github.com/stretchr/testify/assert"
)

func TestCheckWindowMaintenance1(t *testing.T) {
	assert := assert.New(t)

	windows := []models.MaintenanceWindows{
		{AviabilityZone: "msk-1a", Start: 0, End: 5},
		{AviabilityZone: "msk-1a", Start: 21, End: 24},
	}

	now := time.Now()
	task := new(models.Task)
	task.StartTime = time.Date(now.Year(), now.Month(), now.Day(), 1, 0, 0, now.Nanosecond(), now.Location())
	task.FinishTime = time.Date(now.Year(), now.Month(), now.Day(), 4, 30, 0, now.Nanosecond(), now.Location())

	result := checkWindowMaintenance(task, windows)
	assert.True(result)
}

func TestCheckWindowMaintenance2(t *testing.T) {
	assert := assert.New(t)

	windows := []models.MaintenanceWindows{
		{AviabilityZone: "msk-1a", Start: 0, End: 5},
		{AviabilityZone: "msk-1a", Start: 21, End: 24},
	}

	now := time.Now()
	task := new(models.Task)
	task.StartTime = time.Date(now.Year(), now.Month(), now.Day()-1, 22, 0, 0, now.Nanosecond(), now.Location())
	task.FinishTime = time.Date(now.Year(), now.Month(), now.Day(), 4, 30, 0, now.Nanosecond(), now.Location())

	result := checkWindowMaintenance(task, windows)
	assert.True(result)
}

func TestCheckWindowMaintenance3(t *testing.T) {
	assert := assert.New(t)

	windows := []models.MaintenanceWindows{
		{AviabilityZone: "msk-1a", Start: 0, End: 5},
		{AviabilityZone: "msk-1a", Start: 21, End: 23},
	}

	now := time.Now()
	task := new(models.Task)
	task.StartTime = time.Date(now.Year(), now.Month(), now.Day()-1, 22, 0, 0, now.Nanosecond(), now.Location())
	task.FinishTime = time.Date(now.Year(), now.Month(), now.Day(), 4, 30, 0, now.Nanosecond(), now.Location())

	result := checkWindowMaintenance(task, windows)
	assert.False(result)
}
