package etc

import (
	"github.com/sanguohot/zxcs-go-spider/pkg/common/zap"
	"github.com/spf13/viper"
	"os"
	"path"
)
// auto generate struct
// https://mholt.github.io/json-to-go/
// use mapstructure to replace json for '_' key words, e.g. rpc_port,big_data
type ZXCS struct {
	Scheme string `json:"scheme"`
	Address string `json:"address"`
	RootPath string `json:"root_path"`
	Map struct {
		HistoryAndMilitary string `json:"historyAndMilitary"`
		History string `json:"history"`
		Military string `json:"military"`
		City string `json:"city"`
		SwordsmanAndGod string `json:"swordsmanAndGod"`
		Swordsman string `json:"swordsman"`
		God string `json:"god"`
		Fantasy0And1 string `json:"fantasy0And1"`
		Fantasy0 string `json:"fantasy0"`
		Fantasy1 string `json:"fantasy1"`
		PublishAndGirl string `json:"publishAndGirl"`
		Publish string `json:"publish"`
		Girl string `json:"girl"`
	} `json:"map"`
}

var (
	defaultFilePath  = "/etc/config.json"
	ViperConfig *viper.Viper
	Config *ZXCS
	serverPath = os.Getenv("ZXCS_GO_SPIDER_PATH")
)

func init()  {
	if serverPath == "" {
		zap.Logger.Fatal("MEDICHAIN_PATH env required")
	}
	zap.Sugar.Infof("zxcs-go-spider path ===> %s", serverPath)
	InitConfig(path.Join(GetServerDir(), defaultFilePath))
}
func InitConfig(filePath string) {
	zap.Sugar.Infof("config: init config path %s", filePath)
	ViperConfig = viper.New()
	if filePath == "" {
		ViperConfig.SetConfigFile(defaultFilePath)
	}else {
		ViperConfig.SetConfigFile(filePath)
	}

	err := ViperConfig.ReadInConfig()
	if err != nil {
		if filePath != defaultFilePath {
			zap.Logger.Fatal(err.Error())
		}
	}
	err = ViperConfig.Unmarshal(&Config)
	if err != nil {
		zap.Logger.Fatal(err.Error())
	}
}
func GetServerDir() string {
	//return GetViperConfig().GetString("server.dir")
	return serverPath
}