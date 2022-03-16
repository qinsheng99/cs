package manage

import (
	"encoding/xml"
	"errors"
	"fmt"
	"time"

	"cve-sa-backend/dao"
	. "cve-sa-backend/iniconf"
	"cve-sa-backend/models"
	"cve-sa-backend/utils"
	"cve-sa-backend/utils/parsexml"
	"gorm.io/gorm"
)

func DeleteCVE(cveId, packageName string) (string, error) {
	tx := DB.Begin()
	if packageName == "" {
		//Delete the cve information corresponding to the specified cveId
		delCves, rowsAffected, err := dao.GetCveDatabaseByCveIdList(&models.CveDatabase{
			CveId: cveId,
		}, tx)
		if err != nil || rowsAffected == 0 {
			return "", errors.New("CVE error")
		}
		for _, delCve := range delCves {
			err = DeleteOneCVE(delCve, tx)
			if err != nil {
				tx.Rollback()
				return "", err
			}
		}
		tx.Commit()
		return fmt.Sprintf("Delete CVE %d record.", len(delCves)), nil
	} else {
		delCve, err := dao.GetOneDatabase(&models.CveDatabase{
			CveId:       cveId,
			PackageName: packageName,
		}, tx)
		if err == nil {
			err = DeleteOneCVE(*delCve, tx)
			if err != nil {
				tx.Rollback()
				return "", err
			}
			tx.Commit()
			return "Delete CVE 1 record.", nil
		} else {
			tx.Rollback()
			return "", err
		}
	}
}

func DeleteOneCVE(delCve models.CveDatabase, tx *gorm.DB) error {
	cveCvrf, err := dao.GetOneCvrf(&models.CveCvrf{
		CveId:       delCve.CveId,
		PackageName: delCve.PackageName,
	}, tx)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if err == nil {
		err = dao.DeleteCvrf(cveCvrf.Id, tx)
		if err != nil {
			return err
		}
	}

	packageList, rowsAffected, err := dao.GetProductPackageList(&models.CveProductPackage{
		CveId:       delCve.CveId,
		PackageName: delCve.PackageName,
	}, tx)
	if err != nil {
		return err
	}
	if rowsAffected > 0 {
		for _, v := range packageList {
			err = dao.DeleteProductPackage(v.Id, tx)
			if err != nil {
				return err
			}
		}
	}

	parserBean, err := dao.GetOneParser(&models.CveParser{
		Cve:         delCve.CveId,
		PackageName: delCve.PackageName,
	}, tx)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if err == nil {
		err = dao.DeleteParser(parserBean.Id, tx)
		if err != nil {
			return err
		}
	}

	err = dao.DeleteCveDatabase(delCve.Id, tx)
	if err != nil {
		return err
	}
	return nil
}

func SyncCve(cveFileName string) error {
	fileByte, err := utils.GetCvrfFile(cveFileName)
	if err != nil {
		return err
	}

	Element := utils.FixedCveXml{}
	err = xml.Unmarshal(fileByte, &Element)
	if err != nil {
		return err
	}
	updateTime := time.Now()
	for _, v := range Element.Vulnerability {
		cve := parsexml.GetCVEDatabase("", "", v, updateTime)

		fmt.Println(cve)
	}

	return nil
}