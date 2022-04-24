package iniconf

import (
	"github.com/olivere/elastic/v7"
)

var EsClient *elastic.Client

func InitEs() error {
	options := []elastic.ClientOptionFunc{
		elastic.SetURL("http://" + Es.Host + ":" + Es.Port),
		elastic.SetSniff(false),
	}

	//if c.Password != "" {
	//	options = append(options, elastic.SetBasicAuth(c.Username, c.Password))
	//}

	client, err := elastic.NewClient(options...)
	if err != nil {
		return err
	}
	EsClient = client
	return nil
}
