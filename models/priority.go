package models

import "database/sql/driver"

type TaskPriority string

const (
	TASK_PRIORITY_NORMAL   TaskPriority = "normal"
	TASK_PRIORITY_CRITICAL TaskPriority = "critical"
)

var ToTaskPriority = map[string]TaskPriority{
	"normal":   TASK_PRIORITY_NORMAL,
	"critical": TASK_PRIORITY_CRITICAL,
}

func (ct *TaskPriority) Scan(value interface{}) error {
	*ct = TaskPriority(value.([]byte))
	return nil
}

func (ct TaskPriority) Value() (driver.Value, error) {
	return string(ct), nil
}
