package models

import (
	"gorm.io/gorm"
	"inspection/database"
)

type CheckJob struct {
	gorm.Model
	Name           string   `json:"name" gorm:"varchar(10);not null;index:idx_check_job,unique"` //任务名称
	ScriptName     string   `json:"script_name" gorm:"varchar(10);not null"`                     //脚本名称
	ClusterName    string   `json:"cluster_name" gorm:"varchar(20)"`                             //集群名称
	DesiredName    string   `json:"desired_name" gorm:"varchar(20);not null"`                    //基线名称
	IpString       string   `json:"ip_string" gorm:"text"`                                       //机器的列表
	JobHasSynced   int      `json:"job_has_synced" gorm:"job_has_synced"`                        //任务是否被同步
	JobHasComplete int      `json:"job_has_complete" gorm:"job_has_complete"`                    //任务是否完成
	IpList         []string `gorm:"-" json:"ip_list"`

	JobWaitCompleteMinutes int64 `gorm:"job_wait_complete_minutes" json:"job_wait_complete_minutes"`

	AllNum     int64 `json:"all_num" gorm:"all_num"`         //任务数量
	SuccessNum int64 `json:"success_num" gorm:"success_num"` //成功数量
	FailedNum  int64 `json:"failed_num" gorm:"failed_num"`   //失败数量
	MissNum    int64 `json:"miss_num" gorm:"miss_num"`       //失联数量
}

func (this *CheckJob) TableName() string {
	return "check_job_t"
}

func (this *CheckJob) CreateOrUpdate() error {

	table := database.DB.Table(this.TableName())
	if this.Name != "" {
		table = table.Where("name = ?", this.Name)
	}
	var checkJob CheckJob
	if err := table.Debug().First(&checkJob).Error; err != nil {
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

func (this *CheckJob) CreateOne() (uint, error) {
	table := database.DB.Table(this.TableName())

	if err := table.Debug().Create(this).Error; err != nil {
		return this.ID, err
	}

	return this.ID, nil
}

func (this *CheckJob) Update() error {
	table := database.DB.Table(this.TableName())
	if this.Name != "" {
		table = table.Where("name = ?", this.Name)
	}
	if this.ID != 0 {
		table = table.Where("id = ?", this.ID)
	}
	if err := table.Debug().Updates(this).Error; err != nil {
		return err
	}
	return nil
}

func (this *CheckJob) UpdateNodeStatus() error {
	table := database.DB.Table(this.TableName())
	if this.Name != "" {
		table = table.Where("name = ?", this.Name)
	}
	if this.ID != 0 {
		table = table.Where("id = ?", this.ID)
	}
	if err := table.Select("job_has_complete", "success_num", "failed_num", "miss_num").Debug().Updates(this).Error; err != nil {
		return err
	}
	return nil
}

func (this *CheckJob) UpdateStatus() error {
	table := database.DB.Table(this.TableName())
	if this.Name != "" {
		table = table.Where("name = ?", this.Name)
	}
	if this.ID != 0 {
		table = table.Where("id = ?", this.ID)
	}
	if err := table.Select("job_has_synced", "job_has_complete").Debug().Updates(this).Error; err != nil {
		return err
	}
	return nil
}

func (this *CheckJob) GetOne() (uint, error) {

	table := database.DB.Table(this.TableName())

	if err := table.Debug().Create(this).Error; err != nil {
		return this.ID, err
	}

	return this.ID, nil
}

func (this *CheckJob) GetList() ([]*CheckJob, error) {
	var checkJobList []*CheckJob
	table := database.DB.Table(this.TableName())
	if err := table.Debug().Find(&checkJobList).Error; err != nil {
		return checkJobList, err
	}
	return checkJobList, nil
}

func (this *CheckJob) GetNotSyncList() ([]*CheckJob, error) {
	var checkJobList []*CheckJob
	table := database.DB.Table(this.TableName())
	table.Where("job_has_synced = 0")

	if err := table.Debug().Find(&checkJobList).Error; err != nil {
		return checkJobList, err
	}
	return checkJobList, nil
}

func (this *CheckJob) GetSyncList() ([]*CheckJob, error) {
	var checkJobList []*CheckJob
	table := database.DB.Table(this.TableName())
	table.Where("job_has_synced = 1")

	if err := table.Debug().Find(&checkJobList).Error; err != nil {
		return checkJobList, err
	}
	return checkJobList, nil
}

func (this *CheckJob) GetNotCompleteList() ([]*CheckJob, error) {
	var checkJobList []*CheckJob
	table := database.DB.Table(this.TableName())
	table.Where("job_has_complete = 0")

	if err := table.Debug().Find(&checkJobList).Error; err != nil {
		return checkJobList, err
	}
	return checkJobList, nil
}

func (this *CheckJob) Delete() error {
	table := database.DB.Table(this.TableName())
	if this.Name != "" {
		table = table.Where("name = ?", this.Name)
	}
	if err := table.Debug().Delete(this).Error; err != nil {
		return err
	}
	return nil
}
