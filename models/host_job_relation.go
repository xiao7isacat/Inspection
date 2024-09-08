package models

import (
	"gorm.io/gorm"
	"inspection/database"
)

type HostJobRelation struct {
	gorm.Model
	HostId       string `json:"host_id" gorm:"host_id"`
	JobId        uint   `json:"job_id" gorm:"job_id"`
	JobHasSynced int    `json:"job_has_synced" gorm:"job_has_synced"`
}

func (this *HostJobRelation) TableName() string {
	return "host_job_relation"
}

func (this *HostJobRelation) CreateOne() (uint, error) {
	table := database.DB.Table(this.TableName())

	if err := table.Debug().Create(this).Error; err != nil {
		return this.ID, err
	}

	return this.ID, nil
}

func (this *HostJobRelation) Update() error {

	table := database.DB.Table(this.TableName())
	if this.JobId != 0 {
		table = table.Where("job_id = ?", this.JobId)
	}

	var hostJobRelation HostJobRelation
	if err := table.Debug().First(&hostJobRelation).Error; err != nil {
		return err
	}
	if err := table.Debug().Updates(this).Error; err != nil {
		return err
	}
	return nil
}

func (this *HostJobRelation) GetOne() error {
	table := database.DB.Table(this.TableName())
	if this.JobId != 0 {
		table = table.Where("job_id = ?", this.JobId)
	}
	if err := table.Debug().First(this).Error; err != nil {
		return err
	}
	return nil
}

func (this *HostJobRelation) GetList() ([]HostJobRelation, error) {
	var hostJobRelationList []HostJobRelation
	table := database.DB.Table(this.TableName())
	if err := table.Debug().Find(&hostJobRelationList).Error; err != nil {
		return hostJobRelationList, err
	}
	return hostJobRelationList, nil
}

func (this *HostJobRelation) Delete() error {
	table := database.DB.Table(this.TableName())
	if this.JobId != 0 {
		table = table.Where("job_id = ?", this.JobId)
	}
	if err := table.Debug().Delete(this).Error; err != nil {
		return err
	}
	return nil
}

func (this *HostJobRelation) CheckExist() (bool, error) {
	var (
		hostJobRelation HostJobRelation
	)
	table := database.DB.Table(this.TableName())
	if this.JobId != 0 {
		table = table.Where("job_id = ?", this.JobId)
	}

	if err := table.Debug().First(&hostJobRelation).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
