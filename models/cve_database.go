package models

import (
	"time"
)

type CveDatabase struct {
	Id                           int64     `gorm:"column:id" json:"id" xml:"id"`
	AffectedProduct              string    `gorm:"column:affected_product" json:"affectedProduct" xml:"affectedProduct"`
	AnnouncementTime             string    `gorm:"column:announcement_tim" json:"announcementTime" xml:"announcementTime"`
	AttackComplexitynvd          string    `gorm:"column:attack_complexitynvd"  json:"attackComplexityNVD" xml:"attackComplexityNVD"`
	AttackComplexityoe           string    `gorm:"column:attack_complexityoe"  json:"attackComplexityOE" xml:"attackComplexityOE"`
	AttackVectornvd              string    `gorm:"column:attack_vectornvd"  json:"attackVectorNVD" xml:"attackVectorNVD"`
	AttackVectoroe               string    `gorm:"column:attack_vectoroe"  json:"attackVectorOE" xml:"attackVectorOE"`
	Availabilitynvd              string    `gorm:"column:availabilitynvd"  json:"availabilityNVD" xml:"availabilityNVD"`
	Availabilityoe               string    `gorm:"column:availabilityoe"  json:"availabilityOE" xml:"availabilityOE"`
	Confidentialitynvd           string    `gorm:"column:confidentialitynvd"  json:"confidentialityNVD" xml:"confidentialityNVD"`
	Confidentialityoe            string    `gorm:"column:confidentialityoe"  json:"confidentialityOE" xml:"confidentialityOE"`
	CveId                        string    `gorm:"column:cve_id" json:"cveId" xml:"cveId"`
	CvsssCorenvd                 string    `gorm:"column:cvsss_corenvd" json:"cvsssCoreNVD" xml:"cvsssCoreNVD"`
	CvsssCoreoe                  string    `gorm:"column:cvsss_coreoe" json:"cvsssCoreOE" xml:"cvsssCoreOE"`
	Integritynvd                 string    `gorm:"column:integritynvd" json:"integrityNVD" xml:"integrityNVD"`
	Integrityoe                  string    `gorm:"column:integrityoe" json:"integrityOE" xml:"integrityOE"`
	NationalCyberAwarenessSystem string    `gorm:"column:national_cyber_awareness_system" json:"nationalCyberAwarenessSystem" xml:"nationalCyberAwarenessSystem"`
	PackageName                  string    `gorm:"column:package_name" json:"packageName" xml:"packageName"`
	PrivilegesRequirednvd        string    `gorm:"column:privileges_requirednvd" json:"privilegesRequiredNVD" xml:"privilegesRequiredNVD"`
	PrivilegesRequiredoe         string    `gorm:"column:privileges_requiredoe" json:"privilegesRequiredOE" xml:"privilegesRequiredOE"`
	Scopenvd                     string    `gorm:"column:scopenvd" json:"scopeNVD" xml:"scopeNVD"`
	Scopeoe                      string    `gorm:"column:scopeoe" json:"scopeOE" xml:"scopeOE"`
	Status                       string    `gorm:"column:status" json:"status" xml:"status"`
	Summary                      string    `gorm:"column:summary" json:"summary" xml:"summary"`
	Type                         string    `gorm:"column:type" json:"type" xml:"type"`
	UserInteractionnvd           string    `gorm:"column:user_interactionnvd" json:"userInteractionNVD" xml:"userInteractionNVD"`
	UserInteractionoe            string    `gorm:"column:user_interactionoe" json:"userInteractionOE" xml:"userInteractionOE"`
	Updateime                    time.Time `gorm:"column:update_time" description:"更新时间" json:"updateTime" xml:"updateTime"`
}

func (c *CveDatabase) TableName() string {
	return "cve_database"
}

type RCveDatabase struct {
	CveDatabase
	Updateime string ` json:"updateTime" xml:"updateTime"`
}
