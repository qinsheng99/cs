package manage

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"strings"
	"time"

	"cve-sa-backend/dao"
	. "cve-sa-backend/iniconf"
	"cve-sa-backend/models"
	"cve-sa-backend/utils"
	_const "cve-sa-backend/utils/const"
	cveSa "cve-sa-backend/utils/entity/cve_sa"
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

func SyncHardwareCompatibility() (string, error) {
	var listZh, listEn []cveSa.HardwareCompatibility
	var err error
	listZh, err = parserHardwareCompatibility("zh")
	if err != nil {
		return fmt.Sprint("获取数据失败:", err.Error()), err
	}
	tx := DB.Begin()
	if listZh != nil && len(listZh) > 0 {
		err = dao.DeleteHardwareForLang("zh", tx)
		if err != nil {
			tx.Rollback()
			return "", err
		}
		hardwareds := assemblyHardware(listZh, "zh")
		err = dao.CreateHardware(hardwareds, tx)
		if err != nil {
			tx.Rollback()
			return "", err
		}
	}

	listEn, err = parserHardwareCompatibility("en")
	if err != nil {
		return fmt.Sprint("获取数据失败:", err.Error()), err
	}

	if listEn != nil && len(listEn) > 0 {
		err = dao.DeleteHardwareForLang("en", tx)
		if err != nil {
			tx.Rollback()
			return "", err
		}
		hardwareds := assemblyHardware(listEn, "en")
		err = dao.CreateHardware(hardwareds, tx)
		if err != nil {
			tx.Rollback()
			return "", err
		}
	}
	tx.Commit()
	return "success", nil
}

// parserHardwareCompatibility Request zh or en json data, and json unmarshal cveSa.HardwareCompatibility
func parserHardwareCompatibility(lang string) ([]cveSa.HardwareCompatibility, error) {
	var bytes []byte
	var err error
	switch lang {
	case "zh":
		bytes, err = utils.HTTPGet(_const.ParserHardwareFileZh)
		if err != nil {
			SLog.Error("http request failed, err :", err)
			return nil, err
		}
	case "en":
		bytes, err = utils.HTTPGet(_const.ParserHardwareFileEn)
		if err != nil {
			SLog.Error("http request failed, err :", err)
			return nil, err
		}
	default:
		return nil, fmt.Errorf("the input parameter `lang` must be zh or en")
	}
	if strings.Index(string(bytes), "<!DOCTYPE html>") > -1 || strings.Index(string(bytes), "<html>") > -1 {
		return nil, nil
	}
	var datas []cveSa.HardwareCompatibility
	err = json.Unmarshal(bytes, &datas)
	if err != nil {
		SLog.Error("json unmarshal failed,err :", err)
		return nil, err
	}
	return datas, nil
}

// assemblyHardware assembly oe_compatibility_hardware insert datas
func assemblyHardware(datas []cveSa.HardwareCompatibility, lang string) []models.OeCompatibilityHardware {
	var list = make([]models.OeCompatibilityHardware, 0, len(datas))
	timeStr := utils.GetCurTime()
	for _, v := range datas {
		if v.CompatibilityConfiguration == nil {
			list = append(list, *v.OeCompatibilityHardware)
			continue
		}
		list = append(list, *joinHardware(v, lang, timeStr))
	}
	return list
}

// joinHardware
func joinHardware(data cveSa.HardwareCompatibility, lang, timeStr string) *models.OeCompatibilityHardware {
	hardware := data.OeCompatibilityHardware
	hardware.ProductInformation = data.CompatibilityConfiguration.ProductInformation
	hardware.CertificationTime = data.CompatibilityConfiguration.CertificationTime
	hardware.CommitID = data.CompatibilityConfiguration.CommitID
	hardware.MotherBoardRevision = data.CompatibilityConfiguration.MotherBoardRevision
	hardware.BiosUefi = data.CompatibilityConfiguration.BiosUefi
	hardware.Cpu = data.CompatibilityConfiguration.Cpu
	hardware.Ram = data.CompatibilityConfiguration.Ram
	hardware.PortsBusTypes = data.CompatibilityConfiguration.PortsBusTypes
	hardware.VideoAdapter = data.CompatibilityConfiguration.VideoAdapter
	hardware.HostBusAdapter = data.CompatibilityConfiguration.HostBusAdapter
	hardware.HardDiskDrive = data.CompatibilityConfiguration.HardDiskDrive
	hardware.Updateime = timeStr
	hardware.Lang = lang
	return hardware
}

func SyncDriverCompatibility() (string, error) {
	var listZh, listEn []models.OeCompatibilityDriverResponse
	var err error
	listZh, err = parserOEDriverCompatibility("zh")
	if err != nil {
		return fmt.Sprint("获取数据失败:", err.Error()), err
	}
	tx := DB.Begin()
	if listZh != nil && len(listZh) > 0 {
		err = dao.DeleteDriverForLang("zh", tx)
		if err != nil {
			tx.Rollback()
			return "", err
		}
		drivers := assemblyDriver(listZh, "zh")
		err = dao.CreateDriver(drivers, tx)
		if err != nil {
			tx.Rollback()
			return "", err
		}
	}

	listEn, err = parserOEDriverCompatibility("en")
	if err != nil {
		return fmt.Sprint("获取数据失败:", err.Error()), err
	}
	if listEn != nil && len(listEn) > 0 {
		err = dao.DeleteDriverForLang("en", tx)
		if err != nil {
			tx.Rollback()
			return "", err
		}
		drivers := assemblyDriver(listEn, "en")
		err = dao.CreateDriver(drivers, tx)
		if err != nil {
			tx.Rollback()
			return "", err
		}
	}
	tx.Commit()
	return "success", nil
}

func parserOEDriverCompatibility(lang string) ([]models.OeCompatibilityDriverResponse, error) {
	var bytes []byte
	var err error
	switch lang {
	case "zh":
		bytes, err = utils.HTTPGet(_const.ParserDriverFileZh)
		if err != nil {
			SLog.Error("http request failed, err :", err)
			return nil, err
		}
	case "en":
		bytes, err = utils.HTTPGet(_const.ParserDriverFileEn)
		if err != nil {
			SLog.Error("http request failed, err :", err)
			return nil, err
		}
	default:
		return nil, fmt.Errorf("the input parameter `lang` must be zh or en")
	}
	if strings.Index(string(bytes), "<!DOCTYPE html>") > -1 {
		return nil, nil
	}
	var datas []models.OeCompatibilityDriverResponse
	err = json.Unmarshal(bytes, &datas)
	if err != nil {
		SLog.Error("json unmarshal failed,err :", err)
		return nil, err
	}
	return datas, nil
}

func assemblyDriver(datas []models.OeCompatibilityDriverResponse, lang string) []models.OeCompatibilityDriver {
	var list = make([]models.OeCompatibilityDriver, 0, len(datas))
	timeStr := utils.GetCurTime()
	for _, v := range datas {
		list = append(list, joinDriver(v, timeStr, lang))
	}
	return list
}

func joinDriver(data models.OeCompatibilityDriverResponse, timeStr, lang string) models.OeCompatibilityDriver {
	return models.OeCompatibilityDriver{
		Architecture: data.Architecture,
		BoardModel:   data.BoardModel,
		ChipModel:    utils.InterfaceToString(data.ChipModel),
		ChipVendor:   data.ChipVendor,
		Deviceid:     utils.InterfaceToString(data.Deviceid),
		DownloadLink: data.DownloadLink,
		DriverDate:   data.DriverDate,
		DriverName:   data.DriverName,
		DriverSize:   data.DriverSize,
		Item:         data.Item,
		Os:           data.Os,
		Sha256:       data.Sha256,
		SsID:         utils.InterfaceToString(data.SsID),
		SvID:         utils.InterfaceToString(data.SvID),
		Type:         data.Type,
		Vendorid:     utils.InterfaceToString(data.Vendorid),
		Version:      data.Version,
		Lang:         lang,
		Updateime:    timeStr,
	}
}

func TransferData(cve string) (string, error) {
	tx := DB.Begin()
	switch cve {
	case "SAreference":
		securitys, err := dao.FindAllSecurityNotice()
		if err != nil {
			return "", err
		}
		var insert = make([]models.CveSecurityNoticeReference, 0)
		fmt.Println(tx, securitys, insert)
	default:
		return "Parameter error", nil
	}
	return "success", nil
}
