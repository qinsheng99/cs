package utils

import (
	"io/ioutil"
	"os"

	"cve-sa-backend/iniconf"

	"github.com/huaweicloud/huaweicloud-sdk-go-obs/obs"
)

func GetCvrfFile(key string) ([]byte, error) {
	var all []byte
	logConf, err := iniconf.Cfg.GetSection("obs")
	if err != nil {
		iniconf.SLog.Error("Fail to load section 'server': ", err)
		return nil, err
	}
	ak := logConf.Key("Ak").MustString("AK")
	sk := logConf.Key("Sk").MustString("SK")
	endpoint := logConf.Key("END_POINT").String()
	client, err := obs.New(os.Getenv(ak), os.Getenv(sk), endpoint, obs.WithSocketTimeout(60000), obs.WithConnectTimeout(60000))
	if err != nil {
		return nil, err
	}
	input := obs.GetObjectInput{
		GetObjectMetadataInput: obs.GetObjectMetadataInput{Bucket: "openeuler-cve-cvrf", Key: "cvrf/" + key},
	}
	object, err := client.GetObject(&input)
	defer client.Close()
	if err != nil {
		return nil, err
	}
	all, err = ioutil.ReadAll(object.Body)
	if err != nil {
		return nil, err
	}
	defer object.Body.Close()
	defer client.Close()
	return all, nil
}
