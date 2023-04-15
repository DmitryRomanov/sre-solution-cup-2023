package models

type TaskType string

const (
	TASK_TYPE_MANUAL TaskType = "manual"
	TASK_TYPE_AUTO   TaskType = "auto"
)

var ToTaskType = map[string]TaskType{
	"manual": TASK_TYPE_MANUAL,
	"auto":   TASK_TYPE_AUTO,
}
