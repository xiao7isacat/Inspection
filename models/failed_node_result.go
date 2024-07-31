package models

import (
	"gorm.io/gorm"
	"inspection/database"
)

type FailedNameResult struct {
	gorm.Model
	JobId      int    `json:"job_id" gorm:"varchar(10);not null;index:idx_failed_node_result,unique"`  //任务名称
	NodeIp     string `json:"node_ip" gorm:"varchar(20);not null;index:idx_failed_node_result,unique"` //节点名称
	ResultJson string `json:"result_json" gorm:"text;not null"`                                        //执行脚本的结果
	ErrMsg     string `json:"err_msg" gorm:"text"`                                                     //报错信息
}

func (this *FailedNameResult) TableName() string {
	return "failed_name_result_t"
}

func (this *FailedNameResult) CreateOrUpdate() error {
	table := database.DB.Table(this.TableName())
	if this.JobId != 0 {
		table = table.Where("name = ?", this.JobId)
	}

	var failedNameResul FailedNameResult
	if err := table.Debug().First(&failedNameResul).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			if err = table.Debug().Create(this).Error; err != nil {
				return err
			}
		}
		return err
	} else {
		if err := table.Updates(this).Error; err != nil {
			return err
		}
	}

	return nil
}

func (this *FailedNameResult) Update() error {
	table := database.DB.Table(this.TableName())
	var failedNameResul FailedNameResult
	if this.JobId != 0 {
		failedNameResul.JobId = this.JobId
		table = table.Where("name = ?", this.JobId)
	}

	if err := table.Updates(&failedNameResul).Error; err != nil {
		return err
	}
	return nil
}

func (this *FailedNameResult) GetOne() error {
	table := database.DB.Table(this.TableName())
	if this.JobId != 0 {
		table = table.Where("name = ?", this.JobId)
	}
	if err := table.Debug().First(this).Error; err != nil {
		return err
	}
	return nil
}

func (this *FailedNameResult) GetList() ([]FailedNameResult, error) {
	var failedNameResulList []FailedNameResult
	table := database.DB.Table(this.TableName())
	if err := table.Debug().Find(&failedNameResulList).Error; err != nil {
		return failedNameResulList, err
	}
	return failedNameResulList, nil
}

func (this *FailedNameResult) Delete() error {
	table := database.DB.Table(this.TableName())
	if this.JobId != 0 {
		table = table.Where("name = ?", this.JobId)
	}
	if err := table.Debug().Delete(this).Error; err != nil {
		return err
	}
	return nil
}
