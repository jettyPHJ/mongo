package config

import (
	"ToMongoAndRedis/utils/log"
	"os"

	"gopkg.in/yaml.v2"
)

type sysParam struct {
	Mgo struct {
		Uri     string `yaml:"uri"`
		MaxConn int    `yaml:"maxConn"`
	}
	Kafka             string `yaml:"kafka"`
	TrackfileSavepath string `yaml:"trackfileSavepath"`
	Redis             string `yaml:"redis"`
}

var Param sysParam
var cfgName = `./config.yaml`

//SysParamInit 初始化，读取配置文件到缓存
func sysParamInit() {
	var err error
	var n int
	var f *os.File

	buf := make([]byte, 1000)

	f, err = os.OpenFile(cfgName, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		goto _exit
	}
	defer f.Close()

	if n, err = f.Read(buf); err != nil {
		log.ErrorLogger.Println(err.Error())
		goto _exit
	}
	if err = yaml.Unmarshal(buf[0:n], &Param); err != nil {
		goto _exit
	}

	return

_exit:
	panic(err)
}

func init() {
	sysParamInit()
}
