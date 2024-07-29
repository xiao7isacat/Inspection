package models

import (
	"gorm.io/gorm"
	"inspection/database"
)

type DesiredResult struct {
	gorm.Model
	Name       string `json:"name" gorm:"varchar(10);not null;index:idx_desired_result,unique"` //基线名称
	ResultJson string `json:"result_json" gorm:"text;not null"`                                 //基线的内容
}

func (this *DesiredResult) TableName() string {
	return "desired_result_t"
}

func (this *DesiredResult) CreateOrUpdate() error {
	var desiredResult DesiredResult
	table := database.DB.Table(this.TableName())
	if err := table.Debug().Where("name = ?", this.Name).First(&desiredResult).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			if err = table.Debug().Create(this).Error; err != nil {
				return err
			}
		}
		return err
	}
	if this.Name != "" {
		desiredResult.Name = this.Name
		table = table.Where("name = ?", this.Name)
	}
	if err := table.Updates(this).Error; err != nil {
		return err
	}
	return nil
}

func (this *DesiredResult) Update() error {
	table := database.DB.Table(this.TableName())
	var desiredResult DesiredResult
	if this.Name != "" {
		desiredResult.Name = this.Name
		table = table.Where("name = ?", this.Name)
	}

	if err := table.Updates(&desiredResult).Error; err != nil {
		return err
	}
	return nil
}

func (this *DesiredResult) GetOne() error {
	table := database.DB.Table(this.TableName())
	if this.Name != "" {
		table = table.Where("name = ?", this.Name)
	}
	if err := table.Debug().First(this).Error; err != nil {
		return err
	}
	return nil
}

func (this *DesiredResult) GetList() ([]DesiredResult, error) {
	var desiredResultList []DesiredResult
	table := database.DB.Table(this.TableName())
	if err := table.Debug().Find(&desiredResultList).Error; err != nil {
		return desiredResultList, err
	}
	return desiredResultList, nil
}

func (this *DesiredResult) Delete() error {
	table := database.DB.Table(this.TableName())
	if this.Name != "" {
		table = table.Where("name = ?", this.Name)
	}
	if err := table.Debug().Delete(this).Error; err != nil {
		return err
	}
	return nil
}
