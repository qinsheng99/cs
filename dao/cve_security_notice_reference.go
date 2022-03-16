package dao

import (
	"cve-sa-backend/iniconf"
	"cve-sa-backend/models"
	_const "cve-sa-backend/utils/const"
)

func GetReferenceByNo(securityNoticeNo string) (datas []models.CveSecurityNoticeReference, err error) {
	err = iniconf.DB.
		Model(&models.CveSecurityNoticeReference{}).
		Where("security_notice_no = ?", securityNoticeNo).
		Where("type = ?", _const.TypeOther).
		Find(&datas).Error
	return
}