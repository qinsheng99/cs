package web

import (
	"context"
	"encoding/json"
	"log"
	"strconv"
	"time"

	"cve-sa-backend/dao"
	"cve-sa-backend/iniconf"
	"cve-sa-backend/utils"
	cveSa "cve-sa-backend/utils/entity/cve_sa"
	"cve-sa-backend/utils/entity/es"

	"github.com/olivere/elastic/v7"
)

type EsHandle struct {
}

const (
	driverIndex = "gotty_user_log"
)

// Refresh Add gotty user operation logs
func (e *EsHandle) Refresh() error {
	var id, limit, page = 0, 50, 1
	t1 := time.Now()
	for {
		datas, err := dao.DefaultCompatibilityDriver.GetAllDataForId(id, limit, page)
		if err != nil {
			return err
		}

		var esdata []es.DriverEsData
		for _, v := range datas {
			esdata = append(esdata, es.DriverEsData{
				Id:           v.Id,
				Architecture: v.Architecture,
				BoardModel:   v.BoardModel,
				ChipModel:    v.ChipModel,
				ChipVendor:   v.ChipVendor,
				Deviceid:     v.Deviceid,
				DownloadLink: v.DownloadLink,
				DriverDate:   v.DriverDate,
				DriverName:   v.DriverName,
				DriverSize:   v.DriverSize,
				Item:         v.Item,
				Lang:         v.Lang,
				Os:           v.Os,
				Sha256:       v.Sha256,
				SsID:         v.SsID,
				SvID:         v.SvID,
				Type:         v.Type,
				Vendorid:     v.Vendorid,
				Version:      v.Version,
				Updateime:    v.Updateime,
			})
			id = int(v.Id)
		}
		if len(datas) == 0 {
			break
		}

		for _, v := range esdata {
			_, err = iniconf.GetEs().
				Index().
				Index(driverIndex).
				Id(strconv.Itoa(int(v.Id))).
				BodyJson(&v).
				Do(context.Background())
			if err != nil {
				log.Println(err)
				continue
			}
		}
	}

	log.Printf("refresh 耗时:%ss", time.Now().Sub(t1).String())
	return nil
}

func (e *EsHandle) Find(req cveSa.OeCompSearchRequest) (data *es.GottyEsDataResp, err error) {
	var sr *elastic.SearchResult
	var drivers []es.Gotty
	page, size := utils.GetPage(req.Pages)
	//h := elastic.NewHighlight()
	//h.Field("lang").PreTags([]string{"<font color='red'>"}...).PostTags([]string{"</font>"}...)
	sr, err = iniconf.GetEs().
		Search(driverIndex).
		Query(searchQuery(req)).
		//Sort("id", false).
		From((page - 1) * size).
		Size(size).
		//Highlight(h).
		Do(context.Background())
	if err != nil {
		return nil, err
	}

	for _, v := range sr.Hits.Hits {
		var driver es.Gotty
		_ = json.Unmarshal(v.Source, &driver)
		drivers = append(drivers, driver)
	}

	data = &es.GottyEsDataResp{Driver: drivers, Total: sr.Hits.TotalHits.Value}
	return
}

func (e *EsHandle) DeleteEs(id string) (err error) {
	_, err = iniconf.GetEs().Delete().Index(driverIndex).Id(id).Do(context.Background())
	if err != nil {
		return
	}
	return
}

func searchQuery(req cveSa.OeCompSearchRequest) *elastic.BoolQuery {
	bo := elastic.NewBoolQuery()

	if req.Lang != "" {
		bo.Must(elastic.NewTermQuery("lang", req.Lang))
	}

	if req.KeyWord != "" {
		opt := []elastic.Query{
			elastic.NewMatchPhraseQuery("driverName", req.KeyWord),
			elastic.NewMatchPhraseQuery("boardModel", req.KeyWord),
			elastic.NewMatchPhraseQuery("chipVendor", req.KeyWord),
		}
		kq := elastic.NewBoolQuery()
		kq.Should(opt...)

		bo.Must(kq)
	}

	return bo
}
