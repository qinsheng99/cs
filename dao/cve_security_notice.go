package dao

import (
	"cve-sa-backend/iniconf"
	"cve-sa-backend/models"
	cveSa "cve-sa-backend/utils/entity/cve_sa"
)

func FindAll(req cveSa.RequestData) {
	q := iniconf.DB
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

	}
}
