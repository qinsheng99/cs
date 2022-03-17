package dao

import (
	"cve-sa-backend/iniconf"
	"cve-sa-backend/models"

	"gorm.io/gorm"
)

func GetOneDatabase(cveDatabase *models.CveDatabase, tx *gorm.DB) (*models.CveDatabase, error) {
	result := tx.Where(cveDatabase).First(cveDatabase)
	return cveDatabase, result.Error
}

func GetCveDatabaseByCveIdList(cveDatabase *models.CveDatabase, tx *gorm.DB) ([]models.CveDatabase, int64, error) {
	var dataBaseList []models.CveDatabase
	result := tx.Where(cveDatabase).Find(&dataBaseList)
	return dataBaseList, result.RowsAffected, result.Error
}

func InsertCveDatabaseList(cveDatabaseList []models.CveDatabase, tx *gorm.DB) error {
	return tx.Create(&cveDatabaseList).Error
}

func DeleteCveDatabase(id int64, tx *gorm.DB) error {
	sqlString := "delete from cve_database where id=?"
	return tx.Exec(sqlString, id).Error
}

func FindAllCveDatabase() (datas []models.CveDatabase, err error) {
	err = iniconf.DB.Model(&models.CveDatabase{}).Find(&datas).Error
	return
}

func UpdateCve(data models.CveDatabase, tx *gorm.DB) (err error) {
	err = tx.Model(&models.CveDatabase{}).Where("id = ?", data.Id).Updates(&data).Error
	return
}
