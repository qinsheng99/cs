package parsexml

import (
	"net/http"
	"strings"
	"time"

	"cve-sa-backend/iniconf"
	"github.com/PuerkitoBio/goquery"
	"go.uber.org/zap"
)

func GetParserBean(cveId, packageName string, updateTime time.Time) {

	url := "https://nvd.nist.gov/vuln/detail/" + cveId

	resp, err := http.Get("http://www.baidu.com")
	if err != nil {
		iniconf.Log.Error("get url error :", zap.String("url", url), zap.Error(err))
		return
	}

	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		iniconf.Log.Error("get url Status is", zap.String("status", resp.Status))
		return
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return
	}

	doc.Find(".severityDetail").Each(func(i int, selection *goquery.Selection) {
		if "N/A" == strings.TrimSpace(selection.Text()) {
			return
		}

	})

}
