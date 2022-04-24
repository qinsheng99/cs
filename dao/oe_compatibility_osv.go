package dao

import (
	"cve-sa-backend/iniconf"
	"cve-sa-backend/models"
	"cve-sa-backend/utils"
	cveSa "cve-sa-backend/utils/entity/cve_sa"

	"gorm.io/gorm"
)

var columnOsv = []string{"id", "architecture", "os_version", "osv_name", "date", "os_download_link", "type", "details", "friendly_link", "total_result", "checksum", "base_openeuler_version", "tools_result", "platform_result", "update_time"}

type compatibilityOsv struct{}

var DefaultCompatibilityOsv = compatibilityOsv{}

func (c compatibilityOsv) OSVFindAll(req cveSa.RequestOsv) (datas []models.OeCompatibilityOsv, total int64, err error) {
	q := iniconf.DB
	page, size := utils.GetPage(req.Pages)
	query := q.Model(&models.OeCompatibilityOsv{})
	if req.KeyWord != "" {
		query = query.Where(
			q.Where("osv_name like ?", "%"+req.KeyWord+"%").
				Or("os_version like ?", "%"+req.KeyWord+"%").
				Or("type like ?", "%"+req.KeyWord+"%"),
		)
	}
	if req.OsvName != "" {
		query.Where("osv_name like ?", req.OsvName)
	}

	if req.Type != "" {
		query = query.Where("type = ?", req.Type)
	}

	if err = query.Count(&total).Error; err != nil {
		iniconf.SLog.Error(err)
		return
	}

	if total == 0 {
		return
	}

	query = query.Select(columnOsv).Order("id desc").Limit(size).Offset((page - 1) * size)
	if err = query.Find(&datas).Error; err != nil {
		iniconf.SLog.Error(err)
		return
	}
	return
}

func (c compatibilityOsv) GetOsvName() (data []string, err error) {
	if err = iniconf.DB.
		Model(&models.OeCompatibilityOsv{}).
		Select("distinct(osv_name) as osvName").
		Order("osv_name asc").
		Pluck("osvName", &data).Error; err != nil {
		return nil, err
	}
	return
}

func (c compatibilityOsv) GetType() (data []string, err error) {
	if err = iniconf.DB.
		Model(&models.OeCompatibilityOsv{}).
		Select("distinct(type) as type").
		Order("type asc").
		Pluck("type", &data).Error; err != nil {
		return nil, err
	}
	return
}

func (c compatibilityOsv) CreateOsv(data models.OeCompatibilityOsv, tx *gorm.DB) error {
	return tx.Model(&models.OeCompatibilityOsv{}).Create(&data).Error
}

func (c compatibilityOsv) GetOneOSV(osv *models.OeCompatibilityOsv) (*models.OeCompatibilityOsv, error) {
	result := iniconf.DB.Where(osv).First(osv)
	return osv, result.Error
}

func (c compatibilityOsv) ExistsOsv(version string, tx *gorm.DB) (bool, error) {
	var exists models.OeCompatibilityOsv
	err := tx.Where("os_version = ?", version).First(&exists).Error
	if err != nil {
		if utils.ErrorNotFound(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (c compatibilityOsv) UpdateOsv(data models.OeCompatibilityOsv, tx *gorm.DB) error {
	return tx.Model(&models.OeCompatibilityOsv{}).Where("os_version = ?", data.OsVersion).Updates(&data).Error
}

func (c compatibilityOsv) DeleteOsv(version string, tx *gorm.DB) error {
	return tx.Exec("delete from oe_compatibility_osv where os_version = ?", version).Error
}
