package models

import (
	"time"
)

type CveParser struct {
	Id             int64     `gorm:"column:id" json:"id"`
	Cve            string    `gorm:"column:cve" json:"cve"`
	Cvss           string    `gorm:"column:cvss" json:"cvss"`
	Exception      string    `gorm:"column:exception" json:"exception"`
	PackageName    string    `gorm:"column:package_name" json:"packageName"`
	Score          string    `gorm:"column:score" json:"score"`
	SeverityDetail string    `gorm:"column:severity_detail" json:"severityDetail"`
	Vector         string    `gorm:"column:vector" json:"vector"`
	Updateime      time.Time `gorm:"column:update_time" description:"更新时间" json:"updateTime"`
}

func (c *CveParser) TableName() string {
	return "cve_parser"
}

type RCveParser struct {
	CveParser
	Updateime string `json:"updateTime"`
}
