package dao

import (
	"cve-sa-backend/iniconf"
	"cve-sa-backend/models"
	"cve-sa-backend/utils"
	cveSa "cve-sa-backend/utils/entity/cve_sa"

	"gorm.io/gorm"
)

type compatibilityDriver struct{}

var DefaultCompatibilityDriver = compatibilityDriver{}

func (c compatibilityDriver) GetOsList(lang string) (data []string, err error) {
	if err = iniconf.DB.
		Model(&models.OeCompatibilityDriver{}).
		Select("distinct(os) as os").
		Where("lang = ?", lang).
		Order("os asc").
		Pluck("os", &data).Error; err != nil {
		return nil, err
	}
	return
}

func (c compatibilityDriver) GetArchitectureList(lang string) (data []string, err error) {
	if err = iniconf.DB.
		Model(&models.OeCompatibilityDriver{}).
		Select("distinct(architecture) as architecture").
		Where("lang = ?", lang).
		Order("architecture asc").
		Pluck("architecture", &data).Error; err != nil {
		return nil, err
	}
	return
}

func (c compatibilityDriver) FindAllDriver(req cveSa.OeCompSearchRequest) (datas []models.OeCompatibilityDriver, total int64, err error) {
	q := iniconf.DB
	query := q.Model(&models.OeCompatibilityDriver{})
	page, size := utils.GetPage(req.Pages)
	if req.Os != "" {
		query = query.Where("os = ?", req.Os)
	}

	if req.Architecture != "" {
		query = query.Where("architecture = ?", req.Architecture)
	}

	if req.Lang != "" {
		query = query.Where("lang = ?", req.Lang)
	}

	if req.KeyWord != "" {
		query = query.Where(
			q.Where("driver_name like ?", "%"+req.KeyWord+"%").
				Or("board_model like ?", "%"+req.KeyWord+"%").
				Or("chip_vendor like ?", "%"+req.KeyWord+"%"),
		)
	}
	if err = query.Count(&total).Error; err != nil {
		iniconf.SLog.Error(err)
		return
	}
	if total == 0 {
		return
	}
	query = query.Order("id desc").Limit(size).Offset((page - 1) * size)
	if err = query.Find(&datas).Error; err != nil {
		iniconf.SLog.Error(err)
		return
	}
	return
}

func (c compatibilityDriver) DeleteDriverForLang(lang string, tx *gorm.DB) (err error) {
	err = tx.Exec("delete from oe_compatibility_driver where lang = ?", lang).Error
	return
}

func (c compatibilityDriver) CreateDriver(datas []models.OeCompatibilityDriver, tx *gorm.DB) (err error) {
	err = tx.Model(&models.OeCompatibilityDriver{}).Create(&datas).Error
	return
}

func (c compatibilityDriver) GetAllDataForId(id, limit, page int) (data []*models.OeCompatibilityDriver, err error) {
	err = iniconf.GetDb().Model(&models.OeCompatibilityDriver{}).
		Where("id > ?", id).
		Limit(limit).
		Offset((page - 1) * limit).
		Find(&data).Error
	return
}
