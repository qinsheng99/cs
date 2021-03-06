package dao

import (
	"cve-sa-backend/iniconf"
	"cve-sa-backend/models"
	"cve-sa-backend/utils"
	cveSa "cve-sa-backend/utils/entity/cve_sa"

	"gorm.io/gorm"
)

var columnCveDatabase = []string{"id", "summary", "cve_id", "cvsss_corenvd", "cvsss_coreoe", "announcement_time", "status", "package_name", "update_time"}

type cveDatabase struct {
}

var DefaultCveDatabase = cveDatabase{}

func (c cveDatabase) DatabaseFindAll(req cveSa.RequestData) (datas []models.CveDatabase, total int64, err error) {
	q := iniconf.DB
	page, size := utils.GetPage(req.Pages)
	query := q.Model(&models.CveDatabase{})
	if req.KeyWord != "" {
		query = query.Where(
			q.Where("cve_id like ?", "%"+req.KeyWord+"%").
				Or("summary like ?", "%"+req.KeyWord+"%"),
		)
	}

	if req.Type != "" {
		query = query.Where("type = ?", req.Type)
	}

	year := utils.InterfaceToString(req.Year)
	if year != "" {
		query.Where("announcement_time like ?", year+"%")
	}

	if req.Status != "" {
		query.Where("status like ?", req.Status)
	}

	if err = query.Count(&total).Error; err != nil {
		iniconf.SLog.Error(err)
		return
	}
	if total == 0 {
		return
	}

	query = query.Select(columnCveDatabase).Order("id desc").Limit(size).Offset((page - 1) * size)
	if err = query.Find(&datas).Error; err != nil {
		iniconf.SLog.Error(err)
		return
	}
	return
}

func (c cveDatabase) GetOneDatabase(cveDatabase *models.CveDatabase, tx *gorm.DB) (*models.CveDatabase, error) {
	result := tx.Where(cveDatabase).First(cveDatabase)
	return cveDatabase, result.Error
}

func (c cveDatabase) GetOneDatabaseTypeTwo(cveDatabase *models.CveDatabase) (*models.CveDatabase, error) {
	result := iniconf.DB.Where(cveDatabase).First(cveDatabase)
	return cveDatabase, result.Error
}

func (c cveDatabase) GetCveDatabaseByCveIdList(cveDatabase *models.CveDatabase, tx *gorm.DB) ([]models.CveDatabase, int64, error) {
	var dataBaseList []models.CveDatabase
	result := tx.Where(cveDatabase).Find(&dataBaseList)
	return dataBaseList, result.RowsAffected, result.Error
}

func (c cveDatabase) InsertCveDatabaseList(cveDatabaseList []models.CveDatabase, tx *gorm.DB) error {
	return tx.Create(&cveDatabaseList).Error
}

func (c cveDatabase) InsertCveDatabase(cveDatabase *models.CveDatabase, tx *gorm.DB) error {
	return tx.Create(cveDatabase).Error
}

func (c cveDatabase) DeleteCveDatabase(id int64, tx *gorm.DB) error {
	sqlString := "delete from cve_database where id=?"
	return tx.Exec(sqlString, id).Error
}

func (c cveDatabase) FindAllCveDatabase() (datas []*models.CveDatabase, err error) {
	err = iniconf.DB.Model(&models.CveDatabase{}).Find(&datas).Error
	return
}

func (c cveDatabase) UpdateCve(data *models.CveDatabase, tx *gorm.DB) (err error) {
	err = tx.Model(&models.CveDatabase{}).Where("id = ?", data.Id).Updates(&data).Error
	return
}

func (c cveDatabase) DeleteCveDatabaseByCveIdAndPackageName(cveId, packageName string, tx *gorm.DB) error {
	sqlString := "delete from cve_database where cve_id = ? and package_name = ?"
	return tx.Exec(sqlString, cveId, packageName).Error
}
