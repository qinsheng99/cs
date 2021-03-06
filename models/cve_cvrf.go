package models

import (
	"time"
)

type CveCvrf struct {
	Id               int64     `gorm:"column:id" json:"id"`
	FileName         string    `gorm:"column:file_name"  json:"fileName"`
	CveId            string    `gorm:"column:cve_id"  json:"cveId"`
	Cvrf             string    `gorm:"column:cvrf"  json:"cvrf"`
	PackageName      string    `gorm:"column:package_name"  json:"packageName"`
	SecurityNoticeNo string    `gorm:"column:security_notice_no"  json:"securityNoticeNo"`
	Updateime        time.Time `gorm:"column:update_time"  json:"updateTime"`
}

func (c *CveCvrf) TableName() string {
	return "cve_cvrf"
}

type RCveCvrf struct {
	CveCvrf
	Updateime string `json:"updateTime,omitempty"`
}
