package models

import (
	"time"
)

type CveCvrfSync struct {
	Id        int64     `gorm:"column:id" json:"id"`
	CvrfFile  string    `gorm:"column:cvrf_file" json:"cvrfFile"`
	Type      string    `gorm:"column:type" json:"type"`
	Updateime time.Time `gorm:"column:update_time" description:"更新时间" json:"updateTime"`
}

func (c *CveCvrfSync) TableName() string {
	return "cve_cvrf_sync"
}

type RCveCvrfSync struct {
	CveCvrfSync
	Updateime string `json:"updateTime"`
}
