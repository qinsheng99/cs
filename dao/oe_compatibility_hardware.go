package dao

import (
	"cve-sa-backend/iniconf"
	"cve-sa-backend/models"
	"cve-sa-backend/utils"
	cveSa "cve-sa-backend/utils/entity/cve_sa"

	"gorm.io/gorm"
)

func GetOsForHardware(lang string) (data []string, err error) {
	err = iniconf.DB.
		Model(&models.OeCompatibilityHardware{}).
		Select("distinct(os_version) as os").
		Where("lang = ?", lang).
		Order("os_version asc").
		Pluck("os", &data).Error
	return
}

func GetArchitectureListForHardware(lang string) (data []string, err error) {
	err = iniconf.DB.
		Model(&models.OeCompatibilityHardware{}).
		Select("distinct(architecture) as architecture").
		Where("lang = ?", lang).
		Order("architecture asc").
		Pluck("architecture", &data).Error
	return
}

func FindAllHardware(req cveSa.OeCompSearchRequest) (datas []*models.OeCompatibilityHardware, total int64, err error) {
	q := iniconf.DB
	query := q.Model(&models.OeCompatibilityHardware{})

	page, size := getPage(req.Pages)
	if req.Architecture != "" {
		query = query.Where("architecture = ?", req.Architecture)
	}

	if req.Os != "" {
		query = query.Where("os_version = ?", req.Os)
	}

	if req.Lang != "" {
		query = query.Where("lang = ?", req.Lang)
	}

	if req.KeyWord != "" {
		query = query.Where(
			q.Where("hardware_factory like ?", "%"+req.KeyWord+"%").
				Or("hardware_model like ?", "%"+req.KeyWord+"%").
				Or("os_version like ?", "%"+req.KeyWord+"%").
				Or("cpu like ?", "%"+req.KeyWord+"%"),
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

func GetOneHardware(data *models.OeCompatibilityHardware) (*models.OeCompatibilityHardware, error) {
	err := iniconf.DB.Where(data).First(data).Error
	if utils.ErrorNotFound(err) {
		return nil, nil
	}
	return data, err
}

func GetCpuList(lang string) (datas []string, err error) {
	err = iniconf.DB.
		Model(&models.OeCompatibilityHardware{}).
		Select("distinct(cpu) as cpu").
		Where("lang = ?", lang).
		Order("cpu asc").
		Pluck("cpu", &datas).Error
	return
}

func DeleteHardwareForLang(lang string, tx *gorm.DB) (err error) {
	err = tx.Exec("delete from oe_compatibility_hardware where lang = ?", lang).Error
	return
}
