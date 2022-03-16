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
