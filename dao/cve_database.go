package dao

import (
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

func DeleteCveDatabase(id int64, tx *gorm.DB) error {
	sqlString := "delete from cve_database where id=?"
	return tx.Exec(sqlString, id).Error
}
