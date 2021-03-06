package dao

import (
	"cve-sa-backend/iniconf"
	"cve-sa-backend/models"

	"gorm.io/gorm"
)

type cveProductPackage struct{}

var DefaultCveProductPackage = cveProductPackage{}

func (c cveProductPackage) GetProductPackageList(cveProductPackage *models.CveProductPackage, tx *gorm.DB) ([]models.CveProductPackage, int64, error) {
	var list []models.CveProductPackage
	result := tx.Where(cveProductPackage).Find(&list)
	return list, result.RowsAffected, result.Error
}

func (c cveProductPackage) GetProductPackageListTypeTwo(cveProductPackage *models.CveProductPackage) ([]models.CveProductPackage, int64, error) {
	var list []models.CveProductPackage
	result := iniconf.DB.Where(cveProductPackage).Find(&list)
	return list, result.RowsAffected, result.Error
}

func (c cveProductPackage) DeleteProductPackage(id int64, tx *gorm.DB) error {
	sqlString := "delete from cve_product_package where id=?"
	return tx.Exec(sqlString, id).Error
}

func (c cveProductPackage) DeleteByCveIdAndPackageNameAndProductName(cveId, packageName, productName string, tx *gorm.DB) error {
	sqlString := "delete from cve_product_package where cve_id = ? and package_name = ? and product_name = ?"
	return tx.Exec(sqlString, cveId, packageName, productName).Error
}

func (c cveProductPackage) InsertProductPackage(productPackage *models.CveProductPackage, tx *gorm.DB) error {
	return tx.Create(productPackage).Error
}
