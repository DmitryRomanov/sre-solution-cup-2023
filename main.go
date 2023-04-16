package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	_ "github.com/DmitryRomanov/sre-solution-cup-2023/docs"
	"github.com/DmitryRomanov/sre-solution-cup-2023/dto"
	"github.com/DmitryRomanov/sre-solution-cup-2023/models"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title API для планирования работ

// @BasePath /
func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", handleRootRequest)
	r.Post("/task/add", handleAddTaskRequest)
	r.Post("/task/cancel/{task_id}", handleCancelTaskRequest)
	r.Get("/task/list", handleTasksListRequest)
	r.Get("/az/list", handleAzListRequest)
	r.Get("/az/maintenance_windows", handleMaintenanceWindowsListRequest)
	r.Get("/*", httpSwagger.WrapHandler)

	initDB()
	go func() {
		for range time.Tick(1 * time.Second) {
			updateTaskStatuses()
		}
	}()

	fmt.Println("Open http://localhost:3000")
	http.ListenAndServe(":3000", r)
}

var (
	mu sync.Mutex
)

// @Summary Добавить задачу
// @Tags     tasks
// @Produce  json
// @Param request body dto.AddTaskRequest true "task info"
// @Success 200 {object} dto.MessageResponse
// @Success 400 {object} dto.MessageResponse
// @Router /task/add [post]
func handleAddTaskRequest(w http.ResponseWriter, r *http.Request) {
	var p dto.AddTaskRequest
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		response := new(dto.MessageResponse)
		response.Success = false
		response.Message = err.Error()
		writeResponse(w, response)
		return
	}

	task := new(models.Task)
	task.Status = models.TASK_STATUS_WAITING
	task.AviabilityZone = p.AviabilityZone
	task.Duration = p.Duration

	startTime, err := time.Parse(time.DateTime, p.StartTime)
	fmt.Println(err)
	task.StartTime = startTime

	deadline, err := time.Parse(time.DateTime, p.Deadline)
	fmt.Println(err)
	task.Deadline = deadline

	task.Type = p.Type
	task.Priority = p.Priority

	duration := time.Duration(task.Duration-1) * time.Second
	finishTime := task.StartTime.Add(duration)
	task.FinishTime = finishTime

	mu.Lock()
	defer mu.Unlock()

	if task.Type == models.TASK_TYPE_AUTO && haveTasks(task) {
		w.WriteHeader(http.StatusLocked)
		response := new(dto.MessageResponse)
		response.Success = false
		response.Message = "Task already exists"
		writeResponse(w, response)
		return
	}

	err = task.ValidateValues()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := new(dto.MessageResponse)
		response.Success = false
		response.Message = err.Error()
		writeResponse(w, response)
		return
	}

	if task.Type == string(models.TASK_TYPE_MANUAL) {
		// отменить автоматические работы
		cancelAutoTasks(task)

		// критическая
		if task.Priority == string(models.TASK_PRIORITY_CRITICAL) {
			cancelManualTasksWithNormalPriority(task)
		}
	}

	if !validAZ(task) {
		w.WriteHeader(http.StatusBadRequest)
		response := new(dto.MessageResponse)
		response.Success = false
		response.Message = "Зона доступности недоступна"
		writeResponse(w, response)
		return
	}

	var windows []models.MaintenanceWindows
	db.Debug().Model(models.MaintenanceWindows{}).Where("aviability_zone = ?", task.AviabilityZone).Find(&windows)

	if !checkWindowMaintenance(task, windows) {
		w.WriteHeader(http.StatusBadRequest)
		response := new(dto.MessageResponse)
		response.Success = false
		response.Message = "Задача не поподает в окно обслуживания"
		writeResponse(w, response)
		return
	}

	db.Create(task)

	response := new(dto.MessageResponse)
	response.Success = true
	response.Message = "Added"
	writeResponse(w, response)
}

func validAZ(task *models.Task) bool {
	var az models.AviabilityZone
	result := db.Model(models.AviabilityZone{}).Where("name = ?", task.AviabilityZone).First(&az)

	if result.RowsAffected > 0 {
		if az.BlockedForAutomatedTask && task.Priority != models.TASK_PRIORITY_CRITICAL {
			// запрещены работы кроме критичных
			return false
		}

		return true
	}

	return false
}

func checkWindowMaintenance(task *models.Task, windows []models.MaintenanceWindows) bool {
	for _, window := range windows {
		if task.StartTime.Format(time.DateOnly) == task.FinishTime.Format(time.DateOnly) {
			if task.StartTime.Hour() >= window.Start && task.FinishTime.Hour() <= window.End {
				return true
			}
		}
	}

	if task.FinishTime.Day()-task.StartTime.Day() == 1 && task.FinishTime.Sub(task.StartTime).Hours() < 12 {
		//разные сутки
		existsAtBegin := false
		existsAtEnd := false

		//проверить разрыв в сутках
		for _, window := range windows {
			if window.Start == 0 {
				existsAtBegin = true
			}
			if window.End == 24 {
				existsAtEnd = true
			}
		}

		return existsAtBegin && existsAtEnd
	}

	return false
}

func haveTasks(newTask *models.Task) bool {
	var tasks []models.Task
	duration := time.Duration(newTask.Duration-1) * time.Second
	finishTime := newTask.StartTime.Add(duration)
	db.Debug().Where(
		"aviability_zone = ? AND status = ? AND ((? BETWEEN start_time AND finish_time) OR (? BETWEEN start_time AND finish_time))",
		newTask.AviabilityZone,
		models.TASK_STATUS_WAITING,
		newTask.StartTime,
		finishTime,
	).Find(&tasks)
	return len(tasks) > 0
}

func writeResponse(w http.ResponseWriter, object interface{}) {
	js, _ := json.Marshal(object)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

// @Summary Список задач
// @Tags     tasks
// @Produce  json
// @Success 200 {object} []models.Task
// @Router /task/list [get]
func handleTasksListRequest(w http.ResponseWriter, r *http.Request) {
	var tasks []models.Task
	db.Debug().Where("status = ?", models.TASK_STATUS_WAITING).Find(&tasks)
	writeResponse(w, tasks)
}

// @Summary Отменить задачу
// @Tags     tasks
// @Produce  json
// @Param task_id path int true "id задачи"
// @Success 200 {object} []models.Task
// @Router /task/cancel/{task_id} [post]
func handleCancelTaskRequest(w http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "task_id")
	var task models.Task
	db.Debug().First(&task, taskID)

	cancelTask(task, "Ручная отмена")

	response := new(dto.MessageResponse)
	response.Success = true
	response.Message = "Canceled"
	writeResponse(w, response)
}

// @Summary Список зон доступности
// @Tags     az
// @Produce  json
// @Success 200 {object} []models.AviabilityZone
// @Router /az/list [get]
func handleAzListRequest(w http.ResponseWriter, r *http.Request) {
	var azs []models.AviabilityZone
	db.Debug().Find(&azs)
	writeResponse(w, azs)
}

// @Summary Список окон обслуживания в зонах доступности
// @Tags     az
// @Produce  json
// @Success 200 {object} []models.MaintenanceWindows
// @Router /az/maintenance_windows [get]
func handleMaintenanceWindowsListRequest(w http.ResponseWriter, r *http.Request) {
	var windows []models.MaintenanceWindows
	db.Debug().Find(&windows)
	writeResponse(w, windows)
}

func handleRootRequest(w http.ResponseWriter, r *http.Request) {
	r2 := new(http.Request)
	*r2 = *r
	r2.RequestURI = r.RequestURI + "index.html"

	httpSwagger.WrapHandler(w, r2)
}

func updateTaskStatuses() {
	db.Debug().Model(&models.Task{}).Where(
		"status = ? AND start_time <= ?",
		models.TASK_STATUS_WAITING,
		time.Now(),
	).Update("status", models.TASK_STATUS_RUNNING)

	db.Debug().Model(&models.Task{}).Where(
		"status = ? AND finish_time <= ?",
		models.TASK_STATUS_RUNNING,
		time.Now(),
	).Update("status", models.TASK_STATUS_FINISHED)
}
