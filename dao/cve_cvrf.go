package dao

import (
	"cve-sa-backend/models"
	"gorm.io/gorm"
)

func GetOneCvrf(cvrf *models.CveCvrf, tx *gorm.DB) (*models.CveCvrf, error) {
	result := tx.Where(cvrf).First(cvrf)
	return cvrf, result.Error
}

func DeleteCvrf(id int64, tx *gorm.DB) error {
	sqlString := "delete from cve_cvrf where id= ?"
	return tx.Exec(sqlString, id).Error
}

func DeleteByCveIdAndPackageName(cveId, packageName string, tx *gorm.DB) error {
	sqlString := "delete from cve_cvrf where cve_id = ? and package_name = ?"
	return tx.Exec(sqlString, cveId, packageName).Error
}

func InsertCvrf(cvrf *models.CveCvrf, tx *gorm.DB) error {
	return tx.Create(cvrf).Error
}
