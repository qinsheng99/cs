package dao

import (
	"cve-sa-backend/iniconf"
	"cve-sa-backend/models"

	"gorm.io/gorm"
)

type compatibilityHardwareAdapter struct{}

var DefaultCompatibilityHardwareAdapter = compatibilityHardwareAdapter{}

func (c compatibilityHardwareAdapter) ByhardwareId(hardwareId int64) (datas []models.OeCompatibilityHardwareAdapter, err error) {
	err = iniconf.DB.
		Model(&models.OeCompatibilityHardwareAdapter{}).
		Where("hardware_id = ?", hardwareId).Find(&datas).Error
	return
}

func (c compatibilityHardwareAdapter) CreateAdapter(data models.OeCompatibilityHardwareAdapter, tx *gorm.DB) error {
	return tx.Model(&models.OeCompatibilityHardwareAdapter{}).Create(&data).Error
}
