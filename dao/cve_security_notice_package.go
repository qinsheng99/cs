package dao

import (
	"cve-sa-backend/iniconf"
	"cve-sa-backend/models"

	"gorm.io/gorm"
)

type securityNoticePackage struct {
}

var DefaultSecurityNoticePackage = securityNoticePackage{}

func (s securityNoticePackage) ByPackageName(pnames []string) (datas []models.CveSecurityNoticePackage, err error) {
	err = iniconf.DB.
		Model(&models.CveSecurityNoticePackage{}).
		Where("package_name in ?", pnames).
		Find(&datas).Error
	return
}

func (s securityNoticePackage) NoticePackageByNoProduct(securityNoticeNo, product string) (datas []models.CveSecurityNoticePackage, err error) {
	err = iniconf.DB.
		Model(&models.CveSecurityNoticePackage{}).
		Where("security_notice_no = ?", securityNoticeNo).
		Where("product_name = ?", product).
		Find(&datas).Error
	return
}

func (s securityNoticePackage) CreatePackage(datas []models.CveSecurityNoticePackage, tx *gorm.DB) (err error) {
	err = tx.Model(&models.CveSecurityNoticePackage{}).Create(&datas).Error
	return
}

func (s securityNoticePackage) GetPackageList(data *models.CveSecurityNoticePackage, tx *gorm.DB) ([]models.CveSecurityNoticePackage, error) {
	var list []models.CveSecurityNoticePackage
	err := tx.Where(data).Find(&list).Error
	return list, err
}

func (s securityNoticePackage) DeletePackages(datas []models.CveSecurityNoticePackage, tx *gorm.DB) error {
	err := tx.Model(&models.CveSecurityNoticePackage{}).Delete(&datas).Error
	return err
}
