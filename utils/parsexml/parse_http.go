package parsexml

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"cve-sa-backend/iniconf"
	"cve-sa-backend/models"
	"cve-sa-backend/utils"
	_const "cve-sa-backend/utils/const"
	"github.com/PuerkitoBio/goquery"
	"go.uber.org/zap"
)

func GetParserBean(cveId, packageName string, updateTime time.Time) (models.RCveParser, error) {
	cveParser := models.RCveParser{}
	url := "https://nvd.nist.gov/vuln/detail/" + cveId

	resp, err := http.Get(url)
	if err != nil {
		iniconf.Log.Error("get url error :", zap.String("url", url), zap.Error(err))
		return cveParser, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		iniconf.Log.Error("get url Status is", zap.String("status", resp.Status))
		return cveParser, errors.New(resp.Status)
	}

	bodyRes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		iniconf.Log.Error("ioutil.ReadAll error", zap.Error(err))
	}
	bodyOne := ioutil.NopCloser(bytes.NewReader(bodyRes))
	defer bodyOne.Close()
	bodyTwo := ioutil.NopCloser(bytes.NewReader(bodyRes))
	defer bodyTwo.Close()

	doc, err := goquery.NewDocumentFromReader(bodyOne)
	if err != nil {
		return cveParser, err
	}

	s := doc.Find(".severityDetail")
	isV3 := true
	for k, _ := range s.Nodes {
		node := s.Eq(k)

		if _const.EMPTY == strings.TrimSpace(node.Text()) {
			isV3 = false
			continue
		}
		text, err := node.Parent().Html()
		if err != nil {
			continue
		}

		in, ok := utils.InterceptString(text, "<span", "</span>")
		if !ok {
			continue
		}
		text = "<span" + in + "</span>"

		href, ok := node.Children().Attr("href")
		if ok {
			q, ok := utils.InterceptString(href, "vector=", "&")
			if !ok {
				continue
			}
			cveParser.Vector = q
			cveParser.SeverityDetail = text
			cveParser.Score = strings.TrimSpace(node.Text())

			if isV3 {
				q, ok = utils.InterceptString(href, "version=", "&")
				if !ok {
					continue
				}
				cveParser.Cvss = "V" + q
			} else {
				cveParser.Cvss = "V2"
			}

		} else {
			continue
		}

	}

	if cveParser.Exception == "" && cveParser.Score == "" && cveParser.Cvss == "" {
		buf := new(bytes.Buffer)
		_, _ = buf.ReadFrom(bodyTwo)
		new := buf.String()
		cveParser.Exception = new
	}

	cveParser.PackageName = packageName
	cveParser.Updateime = updateTime.Format("2006-01-02 15:04:05")

	return cveParser, nil
}
