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