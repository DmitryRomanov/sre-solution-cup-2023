package dto

type AddTaskRequest struct {
	Name           string `json:"name"`
	AviabilityZone string `json:"aviability_zone"`
	Type           string `json:"type"`
	Priority       string `json:"priority"`
	StartTime      string `json:"strat_time"`
	Duration       int
	Deadline       string `json:"deadline"`
}

type AddTaskResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
