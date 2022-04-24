package web

import (
	"context"
	"log"
	"strconv"
	"time"

	"cve-sa-backend/dao"
	"cve-sa-backend/iniconf"
	cveSa "cve-sa-backend/utils/entity/cve_sa"
	"cve-sa-backend/utils/entity/es"
)

type EsHandle struct {
}

const (
	driverIndex = "driver"
)

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

		if len(datas) == 0 {
			break
		}
	}

	log.Printf("refresh 耗时:%ss", time.Now().Sub(t1).String())
	return nil
}

func (e *EsHandle) Find(req cveSa.OeCompSearchRequest) (data []es.DriverEsData, err error) {
	return nil, err
}
