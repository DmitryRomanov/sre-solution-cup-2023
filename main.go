package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	_ "github.com/DmitryRomanov/sre-solution-cup-2023/docs"
	"github.com/DmitryRomanov/sre-solution-cup-2023/dto"
	"github.com/DmitryRomanov/sre-solution-cup-2023/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

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
	r.Get("/task/list", handleTasksListRequest)
	r.Get("/az/list", handleAzListRequest)
	r.Get("/*", httpSwagger.WrapHandler)

	initDB()
	fmt.Println("Open http://localhost:3000")
	http.ListenAndServe(":3000", r)
}

var (
	db *gorm.DB
)

func initDB() {
	var err error
	db, err = gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	if err != nil {
		panic("Failed to open the SQLite database.")
	}

	db.AutoMigrate(&models.Task{}, &models.AviabilityZone{})

	azs := []models.AviabilityZone{
		{Name: "msk-1a", DataCenter: "msk-1", BlockedForAutomatedTask: false},
		{Name: "msk-1b", DataCenter: "msk-1", BlockedForAutomatedTask: false},
		{Name: "msk-1c", DataCenter: "msk-1", BlockedForAutomatedTask: true},

		{Name: "msk-2a", DataCenter: "msk-2", BlockedForAutomatedTask: true},
		{Name: "msk-2b", DataCenter: "msk-2", BlockedForAutomatedTask: false},
		{Name: "msk-2c", DataCenter: "msk-2", BlockedForAutomatedTask: false},

		{Name: "nsk-1a", DataCenter: "nsk-1", BlockedForAutomatedTask: false},
		{Name: "nsk-1b", DataCenter: "nsk-1", BlockedForAutomatedTask: false},
		{Name: "nsk-1c", DataCenter: "nsk-1", BlockedForAutomatedTask: false},
	}

	db.Create(azs)
}

// @Summary Добавить задачу
// @Tags     tasks
// @Produce  json
// @Param request body dto.AddTaskRequest true "task info"
// @Success 200 {object} dto.AddTaskResponse
// @Success 400 {object} dto.AddTaskResponse
// @Router /task/add [post]
func handleAddTaskRequest(w http.ResponseWriter, r *http.Request) {
	var p dto.AddTaskRequest
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		response := new(dto.AddTaskResponse)
		response.Success = false
		response.Message = err.Error()
		writeResponse(w, response)
		return
	}

	task := new(models.Task)
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

	if task.Type == string(models.TASK_TYPE_AUTO) && haveTasks(task) {
		w.WriteHeader(http.StatusLocked)
		response := new(dto.AddTaskResponse)
		response.Success = false
		response.Message = "Task already exists"
		writeResponse(w, response)
		return
	}

	err = task.ValidateValues()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := new(dto.AddTaskResponse)
		response.Success = false
		response.Message = err.Error()
		writeResponse(w, response)
		return
	}

	if task.Type == string(models.TASK_TYPE_MANUAL) {
		// отменить автоматические работы
		cancelAutoTasks(task)
	}

	db.Create(task)
	response := new(dto.AddTaskResponse)
	response.Success = true
	response.Message = "Added"
	writeResponse(w, response)
}

func cancelAutoTasks(newTask *models.Task) {
	duration := time.Duration(newTask.Duration-1) * time.Second
	finishTime := newTask.StartTime.Add(duration)
	db.Debug().Delete("aviability_zone = ? AND type = ? AND ((? BETWEEN start_time AND finish_time) OR (? BETWEEN start_time AND finish_time))", newTask.AviabilityZone, models.TASK_TYPE_AUTO, newTask.StartTime, finishTime)
}

func haveTasks(newTask *models.Task) bool {
	var tasks []models.Task
	duration := time.Duration(newTask.Duration-1) * time.Second
	finishTime := newTask.StartTime.Add(duration)
	db.Debug().Where("aviability_zone = ? AND ((? BETWEEN start_time AND finish_time) OR (? BETWEEN start_time AND finish_time))", newTask.AviabilityZone, newTask.StartTime, finishTime).Find(&tasks)
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
	db.Debug().Find(&tasks)
	writeResponse(w, tasks)
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

func handleRootRequest(w http.ResponseWriter, r *http.Request) {
	r2 := new(http.Request)
	*r2 = *r
	r2.RequestURI = r.RequestURI + "index.html"

	httpSwagger.WrapHandler(w, r2)
}
