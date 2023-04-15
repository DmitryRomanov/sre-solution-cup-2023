package models

import "database/sql/driver"

type TaskType string

const (
	TASK_TYPE_MANUAL TaskType = "manual"
	TASK_TYPE_AUTO   TaskType = "auto"
)

var ToTaskType = map[string]TaskType{
	"manual": TASK_TYPE_MANUAL,
	"auto":   TASK_TYPE_AUTO,
}

func (ct *TaskType) Scan(value interface{}) error {
	*ct = TaskType(value.([]byte))
	return nil
}

func (ct TaskType) Value() (driver.Value, error) {
	return string(ct), nil
}
