package dao

import (
	"cve-sa-backend/iniconf"
	"cve-sa-backend/models"
)

func ByhardwareId(hardwareId int64) (datas []models.OeCompatibilityHardwareAdapter, err error) {
	err = iniconf.DB.
		Model(&models.OeCompatibilityHardwareAdapter{}).
		Where("hardware_id = ?", hardwareId).Find(&datas).Error
	return
}
