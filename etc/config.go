package etc

import (
	"github.com/sanguohot/zxcs-go-spider/pkg/common/log"
	"github.com/spf13/viper"
	"os"
	"path"
)
// auto generate struct
// https://mholt.github.io/json-to-go/
// use mapstructure to replace json for '_' key words, e.g. rpc_port,big_data
type ConfigStruct struct {
	Zxcs struct {
		Scheme   string `json:"scheme"`
		Address  string `json:"address"`
		DownloadAddress string `mapstructure:"download_address"`
		RootPath string `mapstructure:"root_path"`
		DownloadPath string `mapstructure:"download_path"`
		VotePath string `mapstructure:"vote_path"`
		TypeList []struct {
			Name  string `json:"name"`
			Value int    `json:"value"`
		} `mapstructure:"type_list"`
	} `json:"zxcs"`
	Sqlite struct {
		DbName string `mapstructure:"db_name"`
		TbName string `mapstructure:"tb_name"`
	} `json:"sqlite"`
}

var (
	defaultFilePath  = "/etc/config.json"
	ViperConfig *viper.Viper
	Config *ConfigStruct
	serverPath = os.Getenv("ZXCS_GO_SPIDER_PATH")
)

func init()  {
	if serverPath == "" {
		serverPath = "./"
		log.Sugar.Warn("ZXCS_GO_SPIDER_PATH env not set, use ./ as default")
	}
	log.Sugar.Infof("ZXCS_GO_SPIDER_PATH ===> %s", serverPath)
	InitConfig(path.Join(GetServerDir(), defaultFilePath))
}
func InitConfig(filePath string) {
	log.Sugar.Infof("config: init config path %s", filePath)
	ViperConfig = viper.New()
	if filePath == "" {
		ViperConfig.SetConfigFile(defaultFilePath)
	}else {
		ViperConfig.SetConfigFile(filePath)
	}

	err := ViperConfig.ReadInConfig()
	if err != nil {
		if filePath != defaultFilePath {
			log.Logger.Fatal(err.Error())
		}
	}
	err = ViperConfig.Unmarshal(&Config)
	if err != nil {
		log.Logger.Fatal(err.Error())
	}
}
func GetServerDir() string {
	//return GetViperConfig().GetString("server.dir")
	return serverPath
}
func GetSqliteTabelFilePath() string {
	return path.Join(serverPath, "sql", Config.Sqlite.TbName)
}
func GetSqliteDbFilePath() string {
	return path.Join(serverPath, "databases", Config.Sqlite.DbName)
}
func GetSqliteDbDirPath() string {
	return path.Join(serverPath, "databases")
}