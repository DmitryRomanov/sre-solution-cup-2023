package dto

type AddTaskRequest struct {
	AviabilityZone string `json:"aviability_zone"  example:"msk-1a"`
	Type           string `json:"type" example:"manual"`
	Priority       string `json:"priority" example:"normal"`
	StartTime      string `json:"start_time" example:"2023-04-15 11:00:00"`
	Duration       int    `json:"duration" example:"1800"`
	Deadline       string `json:"deadline" example:"2023-04-16 11:00:00"`
}

type MessageResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
