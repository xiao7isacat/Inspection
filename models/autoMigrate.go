package models

import "inspection/database"

func AutoMigrat() {
	database.DB.AutoMigrate(&CheckScript{})
	database.DB.AutoMigrate(&CheckJob{})
	database.DB.AutoMigrate(&FailedNodeResult{})
	database.DB.AutoMigrate(&DesiredResult{})
}
