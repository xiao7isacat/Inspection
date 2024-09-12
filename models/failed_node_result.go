package models

import (
	"gorm.io/gorm"
	"inspection/database"
)

type FailedNodeResult struct {
	gorm.Model
	JobId        int64  `json:"job_id" gorm:"varchar(10);not null"`  //任务名称
	NodeIp       string `json:"node_ip" gorm:"varchar(20);not null"` //节点名称
	ResultJson   string `json:"result_json" gorm:"text;not null"`    //执行脚本的结果
	Succeed      bool   `json:"succeed" gorm:"-"`
	FinalSucceed int    `json:"final_succeed" gorm:"final_succeed"`
	ErrMsg       string `json:"err_msg" gorm:"text"` //报错信息
}

func (this *FailedNodeResult) TableName() string {
	return "failed_name_result_t"
}

func (this *FailedNodeResult) CreateOrUpdate() (uint, error) {
	table := database.DB.Table(this.TableName())
	if this.JobId != 0 {
		table = table.Where("job_id = ?", this.JobId)
	}

	var FailedNodeResul FailedNodeResult
	if err := table.Debug().First(&FailedNodeResul).Error; err != nil {
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

func (this *FailedNodeResult) CreateOne() (uint, error) {
	table := database.DB.Table(this.TableName())

	if err := table.Debug().Create(this).Error; err != nil {
		return this.ID, err
	}

	return this.ID, nil
}

func (this *FailedNodeResult) Update() error {
	table := database.DB.Table(this.TableName())
	if this.JobId != 0 {
		table = table.Where("job_id = ?", this.JobId)
	}
	if err := table.Updates(&this).Error; err != nil {
		return err
	}
	return nil
}

func (this *FailedNodeResult) GetOne() error {
	table := database.DB.Table(this.TableName())
	if this.JobId != 0 {
		table = table.Where("job_id = ?", this.JobId)
	}
	if this.NodeIp != "" {
		table = table.Where("node_ip = ?", this.NodeIp)
	}
	if err := table.Debug().First(this).Error; err != nil {
		return err
	}
	return nil
}

func (this *FailedNodeResult) GetList() ([]FailedNodeResult, error) {
	var FailedNodeResulList []FailedNodeResult
	table := database.DB.Table(this.TableName())
	if this.JobId != 0 {
		table = table.Where("job_id = ?", this.JobId)
	}

	if this.FinalSucceed != 0 {
		table = table.Where("final_succeed = ?", this.FinalSucceed)
	}

	if err := table.Debug().Find(&FailedNodeResulList).Error; err != nil {
		/*if err == gorm.err {
			return FailedNodeResulList, nil
		}*/
		return FailedNodeResulList, err
	}
	return FailedNodeResulList, nil
}

func (this *FailedNodeResult) Delete() error {
	table := database.DB.Table(this.TableName())
	if this.JobId != 0 {
		table = table.Where("job_id = ?", this.JobId)
	}
	if err := table.Debug().Delete(this).Error; err != nil {
		return err
	}
	return nil
}
