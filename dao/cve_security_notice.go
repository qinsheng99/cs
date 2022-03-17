package dao

import (
	"cve-sa-backend/iniconf"
	"cve-sa-backend/models"
	"cve-sa-backend/utils"
	cveSa "cve-sa-backend/utils/entity/cve_sa"
)

var column = []string{"id", "summary", "cve_id", "affected_product", "announcement_time", "affected_component", "security_notice_no", "update_time"}

const (
	Page = 1
	Size = 10
)

func SecurityFindAll(req cveSa.RequestData) (datas []models.CveSecurityNotice, total int64, err error) {
	q := iniconf.DB
	page, size := getPage(req.Pages)
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

	if req.Year != "" {
		query.Where("announcement_time like ?", req.Year+"%")
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

func NoticeByCveId(cveId string) (datas []models.CveSecurityNotice, err error) {
	err = iniconf.DB.
		Model(&models.CveSecurityNotice{}).
		Where("cve_id like ?", "%"+cveId+"%").
		Find(&datas).Error
	return
}

func NoticeByNo(securityNoticeNo string) (*models.CveSecurityNotice, error) {
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

func NoticeByCveIdComponent(cveId, affectedComponent string) (datas []models.CveSecurityNotice, err error) {
	err = iniconf.DB.
		Model(&models.CveSecurityNotice{}).
		Where("cve_id like ?", "%"+cveId+"%").
		Where("affected_component = ?", affectedComponent).
		Find(&datas).Error
	return
}

func FindAllSecurityNotice() (datas []models.CveSecurityNotice, err error) {
	err = iniconf.DB.Model(&models.CveSecurityNotice{}).Find(&datas).Error
	return
}

func getPage(req cveSa.Pages) (int, int) {
	var page, size int
	if req.Page == 0 {
		page = Page
	} else {
		page = req.Page
	}
	if req.Size == 0 {
		size = Size
	} else {
		size = req.Size
	}
	return page, size
}
