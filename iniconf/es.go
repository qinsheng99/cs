package iniconf

import (
	"github.com/olivere/elastic/v7"
)

var EsClient *elastic.Client

func InitEs() error {
	options := []elastic.ClientOptionFunc{
		//elastic.SetURL("http://" + Es.Host + ":" + Es.Port),
		elastic.SetURL("http://192.168.1.218:9200"),
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

func GetEs() *elastic.Client {
	return EsClient
}

//if len(data) == 4 && data[0] == uint8(27) && data[1] == uint8(91) && data[2] == uint8(75) && data[3] == uint8(8) {
//							if index > 0 {
//								index -= 1
//							}
//							command = command[:index] + command[index+1:]
//							continue
//						}

//func EqualNR(data []byte) bool {
//	equal := true
//	if len(data) != 2 || len(data) != 7 {
//		return false
//	}
//
//	if len(data) == 2{
//		for i := 0; i < len(data); i++ {
//			if data[i] != nr[i] {
//				equal = false
//				break
//			}
//		}
//		return equal
//	} else {
//		for i := 0; i < len(data); i++ {
//			if data[i] != zshnr[i] {
//				equal = false
//				break
//			}
//		}
//		return equal
//	}
//}
// zsh whether it is \r\n
//	zshnr = []byte{27,91,63,49,108,27,62}
