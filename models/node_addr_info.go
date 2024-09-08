package models

import (
	"gorm.io/gorm"
	"inspection/database"
)

type IpAddrInfo struct {
	gorm.Model
	Ip string `json:"ip" gorm:"varchar(10);not null;unique"`
}

func (this *IpAddrInfo) TableName() string {
	return "ip_addr_info_t"
}

func (this *IpAddrInfo) CreateOne() (uint, error) {
	table := database.DB.Table(this.TableName())

	if err := table.Debug().Create(this).Error; err != nil {
		return this.ID, err
	}

	return this.ID, nil
}

func (this *IpAddrInfo) Update() error {

	table := database.DB.Table(this.TableName())
	if this.Ip != "" {
		table = table.Where("ip = ?", this.Ip)
	}
	var ipAddrInfo IpAddrInfo
	if err := table.Debug().First(&ipAddrInfo).Error; err != nil {
		return err
	}
	if err := table.Debug().Updates(this).Error; err != nil {
		return err
	}
	return nil
}

func (this *IpAddrInfo) GetOne() error {
	table := database.DB.Table(this.TableName())
	if this.Ip != "" {
		table = table.Where("ip = ?", this.Ip)
	}
	if err := table.Debug().First(this).Error; err != nil {
		return err
	}
	return nil
}

func (this *IpAddrInfo) GetList() ([]IpAddrInfo, error) {
	var ipAddrInfoList []IpAddrInfo
	table := database.DB.Table(this.TableName())
	if err := table.Debug().Find(&ipAddrInfoList).Error; err != nil {
		return ipAddrInfoList, err
	}
	return ipAddrInfoList, nil
}

func (this *IpAddrInfo) Delete() error {
	table := database.DB.Table(this.TableName())
	if this.Ip != "" {
		table = table.Where("ip = ?", this.Ip)
	}
	if err := table.Debug().Delete(this).Error; err != nil {
		return err
	}
	return nil
}

func (this *IpAddrInfo) CheckExist() (bool, error) {
	var (
		ipAddrInfo IpAddrInfo
	)
	table := database.DB.Table(this.TableName())
	if this.Ip != "" {
		table = table.Where("ip = ?", this.Ip)
	}

	if err := table.Debug().First(&ipAddrInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
