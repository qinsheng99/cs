package dao

import (
	"cve-sa-backend/models"

	"gorm.io/gorm"
)

func GetOneParser(parser *models.CveParser, tx *gorm.DB) (*models.CveParser, error) {
	result := tx.Where(parser).First(parser)
	return parser, result.Error
}

func DeleteParser(id int64, tx *gorm.DB) error {
	sqlString := "delete from cve_parser where id=?"
	return tx.Exec(sqlString, id).Error
}

func DeleteByCveAndPackageName(cve, packageName string, tx *gorm.DB) error {
	sqlString := "delete from cve_parser where cve = ? and package_name = ?"
	return tx.Exec(sqlString, cve, packageName).Error
}

func InsertParser(parser *models.CveParser, tx *gorm.DB) error {
	return tx.Create(parser).Error
}
