package main

import (
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/DmitryRomanov/sre-solution-cup-2023/docs"
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
	r.Post("/add_task", handleAddTaskRequest)
	r.Get("/az/list", handleAzListRequest)
	r.Get("/*", httpSwagger.WrapHandler)

	initDB()
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
		{ID: 1, Name: "msk-1a", DataCenter: "msk-1", BlockedForAutomatedTask: false},
		{ID: 2, Name: "msk-1b", DataCenter: "msk-1", BlockedForAutomatedTask: false},
		{ID: 3, Name: "msk-1c", DataCenter: "msk-1", BlockedForAutomatedTask: true},

		{ID: 4, Name: "msk-2a", DataCenter: "msk-2", BlockedForAutomatedTask: true},
		{ID: 5, Name: "msk-2b", DataCenter: "msk-2", BlockedForAutomatedTask: false},
		{ID: 6, Name: "msk-2c", DataCenter: "msk-2", BlockedForAutomatedTask: false},

		{ID: 7, Name: "nsk-1a", DataCenter: "nsk-1", BlockedForAutomatedTask: false},
		{ID: 8, Name: "nsk-1b", DataCenter: "nsk-1", BlockedForAutomatedTask: false},
		{ID: 9, Name: "nsk-1c", DataCenter: "nsk-1", BlockedForAutomatedTask: false},
	}

	db.Create(azs)
}

// @Summary Добавить задачу
// @Produce  json
// @Param request body dto.AddTaskRequest true "task info"
// @Success 200 {object} dto.AddTaskResponse
// @Router /add_task [post]
func handleAddTaskRequest(w http.ResponseWriter, r *http.Request) {

}

// @Summary Список зон доступности
// @Produce  json
// @Success 200 {object} []models.AviabilityZone
// @Router /az/list [get]
func handleAzListRequest(w http.ResponseWriter, r *http.Request) {
	var azs []models.AviabilityZone
	db.Debug().Find(&azs)
	js, err := json.Marshal(azs)
	if nil != err {
		log.Panicf("Can not marshall response %v", azs)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func handleRootRequest(w http.ResponseWriter, r *http.Request) {
	r2 := new(http.Request)
	*r2 = *r
	r2.RequestURI = r.RequestURI + "index.html"

	httpSwagger.WrapHandler(w, r2)
}
