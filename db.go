package main

import (
	"github.com/DmitryRomanov/sre-solution-cup-2023/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	db *gorm.DB
)

func initDB() {
	var err error
	db, err = gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	if err != nil {
		panic("Failed to open the SQLite database.")
	}

	db.AutoMigrate(
		&models.Task{}, &models.AviabilityZone{}, &models.CancelReason{},
		&models.MaintenanceWindows{},
	)

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

	maintenanceWindows := []models.MaintenanceWindows{
		{AviabilityZone: "msk-1a", Start: 0, End: 5},
		{AviabilityZone: "msk-1a", Start: 21, End: 24},

		{AviabilityZone: "msk-1b", Start: 0, End: 24},
		{AviabilityZone: "msk-1c", Start: 23, End: 24},

		{AviabilityZone: "msk-2a", Start: 0, End: 6},
		{AviabilityZone: "msk-2b", Start: 0, End: 6},
		{AviabilityZone: "msk-2c", Start: 0, End: 6},

		{AviabilityZone: "nsk-1a", Start: 4, End: 10},
		{AviabilityZone: "nsk-1b", Start: 4, End: 10},
		{AviabilityZone: "nsk-1c", Start: 4, End: 10},
	}
	db.Create(maintenanceWindows)
}
