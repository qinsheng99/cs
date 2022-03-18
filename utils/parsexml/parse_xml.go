package parsexml

import (
	"encoding/xml"
	"strings"
	"time"

	"cve-sa-backend/iniconf"
	"cve-sa-backend/models"
	"cve-sa-backend/utils"
	_const "cve-sa-backend/utils/const"
	cveSa "cve-sa-backend/utils/entity/cve_sa"
	"go.uber.org/zap"
)

func GetCVEDatabase(sa, affectedComponent string, child utils.Vulnerability, updateTime time.Time) (cveSa.DatabaseData, error) {

	cve := SetVulnerability(sa, affectedComponent, child, updateTime)
	bean, err := GetParserBean(cve.CveId, cve.PackageName, updateTime)
	if err != nil {
		return cve, err
	}
	cve.NationalCyberAwarenessSystem = bean.Cvss
	if cve.Type == "" {
		if bean.Score != "" {
			score, ok := utils.InterceptString(bean.Score, " ", "")
			if ok {
				score = strings.TrimSpace(score)
				score = strings.ToUpper(score[:1]) + strings.ToLower(score[1:])
				cve.Type = score
			}
		}
	}

	if bean.Score != "" {
		score, ok := utils.InterceptString(bean.Score, " ", "")
		if ok {
			cve.CvsssCorenvd = strings.TrimSpace(score)
		}

	}

	if bean.Vector != "" {
		vector := utils.GetVectorArr(bean.Vector)
		cve.AttackVectoroe = vector.AV
		cve.AttackComplexityoe = vector.AC
		cve.PrivilegesRequiredoe = vector.PR
		cve.UserInteractionoe = vector.UI
		cve.Scopeoe = vector.S
		cve.Confidentialityoe = vector.C
		cve.Integrityoe = vector.I
		cve.Availabilityoe = vector.A
	}
	cve.ParserBean = &bean
	return cve, nil
}

type Vulnerability struct {
	Text    string `xml:",innerxml"`
	Ordinal string `xml:"Ordinal,attr"`
	Xmlns   string `xml:"xmlns,attr"`
}

func SetVulnerability(sa, affectedComponent string, child utils.Vulnerability, updateTime time.Time) cveSa.DatabaseData {
	cve := cveSa.DatabaseData{}
	cve.AffectedProduct = sa

	var statusType string
	var productList []string
	if len(child.Notes.Note) > 0 {
		cve.Summary = child.Notes.Note[0].Note
	}

	cve.AnnouncementTime = utils.GetTextTrim(child.ReleaseDate)

	cve.CveId = utils.GetTextTrim(child.Cve)

	if len(child.ProductStatuses.Status) > 0 {
		for _, v := range child.ProductStatuses.Status {
			statusType = v.Type
			for _, nv := range v.ProductID {
				productList = append(productList, nv)
			}
		}
	}

	if len(child.Threats.Threat) > 0 {
		cve.Type = child.Threats.Threat[0].Description
	}

	if len(child.CVSSScoreSets.ScoreSet) > 0 {
		cve.CvsssCoreoe = child.CVSSScoreSets.ScoreSet[0].BaseScore
		if child.CVSSScoreSets.ScoreSet[0].Vector != "" {
			vector := utils.GetVectorArr(child.CVSSScoreSets.ScoreSet[0].Vector)
			cve.AttackVectoroe = vector.AV
			cve.AttackComplexityoe = vector.AC
			cve.PrivilegesRequiredoe = vector.PR
			cve.UserInteractionoe = vector.UI
			cve.Scopeoe = vector.S
			cve.Confidentialityoe = vector.C
			cve.Integrityoe = vector.I
			cve.Availabilityoe = vector.A
		}
	}

	for _, v := range child.Remediations.Remediation {
		if strings.ToLower(statusType) == strings.ToLower("Unaffected") {
			if v.Description != "" {
				affectedComponent = utils.GetTextTrim(v.Description)
			}
		} else {
			break
		}
	}

	var packageList []models.RCveProductPackage
	for _, v := range productList {
		packageObj := models.RCveProductPackage{}
		packageObj.CveId = cve.CveId
		packageObj.PackageName = affectedComponent
		packageObj.ProductName = v
		packageObj.Status = statusType
		packageObj.Updateime = updateTime.Format("2006-01-02 15:04:05")
		packageList = append(packageList, packageObj)
	}

	cve.PackageName = affectedComponent
	cve.PackageList = packageList
	cve.Updateime = updateTime.Format("2006-01-02 15:04:05")

	vu := Vulnerability{}
	vu.Text = child.Text
	vu.Ordinal = child.Ordinal
	vu.Xmlns = child.Xmlns

	by, err := xml.Marshal(&vu)
	if err != nil {
		iniconf.Log.Error("xml.Marshal error:", zap.String("error", err.Error()))
	}

	sb := new(strings.Builder)
	sb.WriteString(_const.TOP + "\n")
	sb.WriteString("<" + _const.CVRFDOC + " xmlns=" + _const.XMLNS + " xmlns:cvrf=" + _const.XmlnsCvrf + ">" + "\n")
	sb.WriteString(string(by) + "\n")
	sb.WriteString("</" + _const.CVRFDOC + ">")

	cvrf := models.RCveCvrf{}
	cvrf.Cvrf = sb.String()
	cvrf.CveId = cve.CveId
	cvrf.PackageName = cve.PackageName
	cvrf.Updateime = updateTime.Format("2006-01-02 15:04:05")
	cve.Cvrf = &cvrf

	return cve
}
