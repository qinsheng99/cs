package models

import (
	"time"
)

type CveSecurityNoticePackage struct {
	Id               int64     `gorm:"column:id" json:"id"`
	PackageLink      string    `gorm:"column:package_link" json:"packageLink"`
	PackageName      string    `gorm:"column:package_name" json:"packageName"`
	ProductName      string    `gorm:"column:product_name" json:"productName"`
	PackageType      string    `gorm:"column:package_type" json:"packageType"`
	SecurityNoticeNo string    `gorm:"column:security_notice_no" json:"securityNoticeNo"`
	Sha256           string    `gorm:"column:sha256" json:"sha256"`
	Updateime        time.Time `gorm:"column:update_time" description:"更新时间" json:"updateTime"`
}

func (c *CveSecurityNoticePackage) TableName() string {
	return "cve_security_notice_package"
}

type RCveSecurityNoticePackage struct {
	CveSecurityNoticePackage
	Updateime string `json:"updateTime"`
}