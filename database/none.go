package database

import "fmt"

type None struct {
}

func (this *None) ConnectDb() error {

	return fmt.Errorf("unsport")
}
