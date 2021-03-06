package dao

import (
	"cve-sa-backend/iniconf"
	"cve-sa-backend/models"
	"cve-sa-backend/utils"
	cveSa "cve-sa-backend/utils/entity/cve_sa"

	"gorm.io/gorm"
)

type compatibilityHardware struct{}

var DefaultCompatibilityHardware = compatibilityHardware{}

func (c compatibilityHardware) GetOsForHardware(lang string) (data []string, err error) {
	err = iniconf.DB.
		Model(&models.OeCompatibilityHardware{}).
		Select("distinct(os_version) as os").
		Where("lang = ?", lang).
		Order("os_version asc").
		Pluck("os", &data).Error
	return
}

func (c compatibilityHardware) GetArchitectureListForHardware(lang string) (data []string, err error) {
	err = iniconf.DB.
		Model(&models.OeCompatibilityHardware{}).
		Select("distinct(architecture) as architecture").
		Where("lang = ?", lang).
		Order("architecture asc").
		Pluck("architecture", &data).Error
	return
}

func (c compatibilityHardware) FindAllHardware(req cveSa.OeCompSearchRequest) (datas []*models.OeCompatibilityHardware, total int64, err error) {
	q := iniconf.DB
	query := q.Model(&models.OeCompatibilityHardware{})

	page, size := utils.GetPage(req.Pages)
	if req.Architecture != "" {
		query = query.Where("architecture = ?", req.Architecture)
	}

	if req.Os != "" {
		query = query.Where("os_version = ?", req.Os)
	}

	if req.Lang != "" {
		query = query.Where("lang = ?", req.Lang)
	}

	if req.Cpu != "" {
		query = query.Where("cpu = ?", req.Cpu)
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

func (c compatibilityHardware) GetOneHardware(data *models.OeCompatibilityHardware) (*models.OeCompatibilityHardware, error) {
	err := iniconf.DB.Where("id = ?", data.Id).First(data).Error
	if utils.ErrorNotFound(err) {
		return nil, nil
	}
	return data, err
}

func (c compatibilityHardware) GetCpuList(lang string) (datas []string, err error) {
	err = iniconf.DB.
		Model(&models.OeCompatibilityHardware{}).
		Select("distinct(cpu) as cpu").
		Where("lang = ?", lang).
		Order("cpu asc").
		Pluck("cpu", &datas).Error
	return
}

func (c compatibilityHardware) DeleteHardwareForLang(lang string, tx *gorm.DB) (err error) {
	err = tx.Exec("delete from oe_compatibility_hardware where lang = ?", lang).Error
	return
}

func (c compatibilityHardware) CreateHardwares(datas []models.OeCompatibilityHardware, tx *gorm.DB) (err error) {
	err = tx.Model(&models.OeCompatibilityHardware{}).Create(&datas).Error
	return
}

func (c compatibilityHardware) CreateHardware(data *models.OeCompatibilityHardware, tx *gorm.DB) (err error) {
	err = tx.Model(&models.OeCompatibilityHardware{}).Create(&data).Error
	return
}
