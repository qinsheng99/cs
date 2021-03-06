package dao

import (
	"cve-sa-backend/iniconf"
	"cve-sa-backend/models"
	"cve-sa-backend/utils"
	cveSa "cve-sa-backend/utils/entity/cve_sa"

	"gorm.io/gorm"
)

var column = []string{"id", "summary", "cve_id", "affected_product", "announcement_time", "affected_component", "security_notice_no", "update_time", "type"}

type securityNotice struct{}

var DefaultSecurityNotice = securityNotice{}

func (s securityNotice) SecurityFindAll(req cveSa.RequestData) (datas []models.CveSecurityNotice, total int64, err error) {
	q := iniconf.DB
	page, size := utils.GetPage(req.Pages)
	query := q.Model(&models.CveSecurityNotice{})
	if req.KeyWord != "" {
		query = query.Where(
			q.Where("security_notice_no like ?", "%"+req.KeyWord+"%").
				Or("summary like ?", "%"+req.KeyWord+"%").
				Or("cve_id like ?", "%"+req.KeyWord+"%").
				Or("description like ?", "%"+req.KeyWord+"%"),
		)
	}
	if req.Type != "" {
		query = query.Where("type = ?", req.Type)
	}

	year := utils.InterfaceToString(req.Year)
	if year != "" {
		query.Where("announcement_time like ?", year+"%")
	}

	if err = query.Count(&total).Error; err != nil {
		iniconf.SLog.Error(err)
		return
	}
	if total == 0 {
		return
	}
	query = query.Select(column).Order("announcement_time desc,security_notice_no desc").Limit(size).Offset((page - 1) * size)
	if err = query.Find(&datas).Error; err != nil {
		iniconf.SLog.Error(err)
		return
	}
	return
}

func (s securityNotice) NoticeByCveId(cveId string) (datas []models.CveSecurityNotice, err error) {
	err = iniconf.DB.
		Model(&models.CveSecurityNotice{}).
		Where("cve_id like ?", "%"+cveId+"%").
		Find(&datas).Error
	return
}

func (s securityNotice) NoticeByCveIdAndAffectedComponent(cveId, affectedComponent string) (datas []models.CveSecurityNotice, err error) {
	err = iniconf.DB.
		Model(&models.CveSecurityNotice{}).
		Where("cve_id like ?", "%"+cveId+"%").
		Where("affected_component = ?", affectedComponent).
		Find(&datas).Error
	return
}

func (s securityNotice) NoticeByNo(securityNoticeNo string) (*models.CveSecurityNotice, error) {
	var data = new(models.CveSecurityNotice)
	if err := iniconf.DB.
		Model(&models.CveSecurityNotice{}).
		Where("security_notice_no = ?", securityNoticeNo).
		First(data).Error; err != nil {
		if utils.ErrorNotFound(err) {
			err = nil
			return nil, err
		}
		return nil, err
	}
	return data, nil
}

func (s securityNotice) GetSecurityNotice(securityNotice *models.CveSecurityNotice, tx *gorm.DB) (*models.CveSecurityNotice, error) {
	result := tx.Where(securityNotice).First(securityNotice)
	return securityNotice, result.Error
}

func (s securityNotice) DeleteSecurityNotice(id int64, tx *gorm.DB) error {
	sqlString := "delete from cve_security_notice where id=?"
	return tx.Exec(sqlString, id).Error
}

func (s securityNotice) NoticeByCveIdComponent(cveId, affectedComponent string) (datas []models.CveSecurityNotice, err error) {
	err = iniconf.DB.
		Model(&models.CveSecurityNotice{}).
		Where("cve_id like ?", "%"+cveId+"%").
		Where("affected_component = ?", affectedComponent).
		Find(&datas).Error
	return
}

func (s securityNotice) FindAllSecurityNotice() (datas []models.CveSecurityNotice, err error) {
	err = iniconf.DB.Model(&models.CveSecurityNotice{}).Find(&datas).Error
	return
}

func (s securityNotice) GetOneSecurity(data *models.CveSecurityNotice, tx *gorm.DB) (*models.CveSecurityNotice, error) {
	err := tx.Where(data).First(data).Error
	if utils.ErrorNotFound(err) {
		return nil, nil
	}
	return data, err
}

func (s securityNotice) DeleteSecurity(id int64, tx *gorm.DB) (err error) {
	err = tx.Exec("delete from cve_security_notice where id = ?", id).Error
	return
}

func (s securityNotice) CreateSecurity(data models.CveSecurityNotice, tx *gorm.DB) (err error) {
	err = tx.Model(&models.CveSecurityNotice{}).Create(&data).Error
	return
}
