package models

import (
	"time"
)

type OeCompatibilityApplication struct {
	Id                  int64     `gorm:"column:id" json:"id"`
	ApplicationSoftware string    `gorm:"column:application_software" json:"applicationSoftware"`
	AsVersion           string    `gorm:"column:as_version" json:"asVersion"`
	Attribute           string    `gorm:"column:attribute" json:"attribute"`
	Category            string    `gorm:"column:category" json:"category"`
	CompatibilityLevel  string    `gorm:"column:compatibility_level" json:"compatibilityLevel"`
	Download            string    `gorm:"column:download" json:"download"`
	FriendlyLink        string    `gorm:"column:friendly_link" json:"friendlyLink"`
	Language            string    `gorm:"column:language" json:"language"`
	Os                  string    `gorm:"column:os" json:"os"`
	OsVersion           string    `gorm:"column:osVersion" json:"osVersion"`
	Remark              string    `gorm:"column:remark" json:"remark"`
	Source              string    `gorm:"column:source" json:"source"`
	SourceProtocol      string    `gorm:"column:source_protocol" json:"sourceProtocol"`
	Updateime           time.Time `gorm:"column:update_time" description:"更新时间" json:"updateTime"`
}

func (o *OeCompatibilityApplication) TableName() string {
	return "oe_compatibility_application"
}

type ROeCompatibilityApplication struct {
	OeCompatibilityApplication
	Updateime string `json:"updateTime"`
}
