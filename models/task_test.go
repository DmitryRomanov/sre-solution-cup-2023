package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestValidateTask(t *testing.T) {
	assert := assert.New(t)
	task := new(Task)
	task.AviabilityZone = "msk-1a"
	task.Duration = 600
	task.StartTime = time.Now()
	task.Deadline = time.Now().Add(time.Minute * 10)
	task.Type = string(TASK_TYPE_AUTO)
	task.Priority = string(TASK_PRIORITY_NORMAL)

	err := task.ValidateValues()
	assert.Nil(err)
}
