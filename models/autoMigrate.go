package models

import "inspection/database"

func AutoMigrat() {
	database.DB.AutoMigrate(&CheckScript{})
	database.DB.AutoMigrate(&CheckJob{})
	database.DB.AutoMigrate(&FailedNameResult{})
	database.DB.AutoMigrate(&DesiredResult{})
}
