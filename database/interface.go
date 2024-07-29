package database

import (
	"gorm.io/gorm"
)

var DB *gorm.DB

type DbConnect interface {
	ConnectDb() error
}

func ConnectDb(databse string) error {
	var dbConnect DbConnect
	switch databse {
	case "sqlite":
		dbConnect = &Sqlite{}

	default:
		dbConnect = &None{}
	}
	if err := dbConnect.ConnectDb(); err != nil {
		return err
	}
	return nil
}
