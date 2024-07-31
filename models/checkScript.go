package models

import (
	"gorm.io/gorm"
	"inspection/database"
)

type CheckScript struct {
	gorm.Model
	Name        string `json:"name" gorm:"varchar(10);not null;index:idx_check_script,unique"` //脚本名称
	ContentJson string `json:"content_json" gorm:"text;not null"`                              //脚本内容
}

func (this *CheckScript) TableName() string {
	return "check_script_t"
}

func (this *CheckScript) CreateOneOrUpdate() (uint, error) {
	table := database.DB.Table(this.TableName())
	if this.Name != "" {
		table = table.Where("name = ?", this.Name)
	}
	var checkScript CheckScript
	if err := table.Debug().First(&checkScript).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			if err = table.Debug().Create(this).Error; err != nil {
				return this.ID, err
			}
		}
		return this.ID, err
	} else {

		if err := table.Updates(this).Error; err != nil {
			return this.ID, err
		}
	}

	return this.ID, nil
}

func (this *CheckScript) CreateOne() (uint, error) {
	table := database.DB.Table(this.TableName())

	if err := table.Debug().Create(this).Error; err != nil {
		return this.ID, err
	}

	return this.ID, nil
}

func (this *CheckScript) Update() error {

	table := database.DB.Table(this.TableName())
	if this.Name != "" {
		table = table.Where("name = ?", this.Name)
	}
	var checkScript CheckScript
	if err := table.Debug().First(&checkScript).Error; err != nil {
		return err
	}
	if err := table.Debug().Updates(this).Error; err != nil {
		return err
	}
	return nil
}

func (this *CheckScript) GetOne() error {
	table := database.DB.Table(this.TableName())
	if this.Name != "" {
		table = table.Where("name = ?", this.Name)
	}
	if err := table.Debug().First(this).Error; err != nil {
		return err
	}
	return nil
}

func (this *CheckScript) GetList() ([]CheckScript, error) {
	var checkScriptList []CheckScript
	table := database.DB.Table(this.TableName())
	if err := table.Debug().Find(&checkScriptList).Error; err != nil {
		return checkScriptList, err
	}
	return checkScriptList, nil
}

func (this *CheckScript) Delete() error {
	table := database.DB.Table(this.TableName())
	if this.Name != "" {
		table = table.Where("name = ?", this.Name)
	}
	if err := table.Debug().Delete(this).Error; err != nil {
		return err
	}
	return nil
}

func (this *CheckScript) CheckExist() (bool, error) {
	var (
		checkScript CheckScript
	)
	table := database.DB.Table(this.TableName())
	if this.Name != "" {
		table = table.Where("name = ?", this.Name)
	}

	if err := table.Debug().First(&checkScript).Error; err != nil {
		return false, err
	}

	return true, nil
}
