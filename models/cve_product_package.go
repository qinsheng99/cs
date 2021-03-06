package models

import (
	"time"
)

type CveProductPackage struct {
	Id          int64     `gorm:"column:id" json:"id"`
	CveId       string    `gorm:"column:cve_id" json:"cveId"`
	PackageName string    `gorm:"column:package_name" json:"packageName"`
	ProductName string    `gorm:"column:product_name" json:"productName"`
	Status      string    `gorm:"column:status" json:"status"`
	Updateime   time.Time `gorm:"column:update_time" description:"更新时间" json:"updateTime"`
}

func (c *CveProductPackage) TableName() string {
	return "cve_product_package"
}

type RCveProductPackage struct {
	CveProductPackage
	Updateime string `json:"updateTime"`
}
