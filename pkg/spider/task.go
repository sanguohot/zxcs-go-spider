package spider

import (
	"fmt"
	"github.com/sanguohot/zxcs-go-spider/etc"
	"github.com/sanguohot/zxcs-go-spider/pkg/common/log"
	"time"
)

// 频率限制
func SpideAll()  {
	// 每隔30分钟爬取一个类型
	channels := make(chan int, len(etc.Config.Zxcs.TypeList))
	for i, _ := range etc.Config.Zxcs.TypeList {
		channels <- i
	}
	close(channels)
	limiter := time.Tick(time.Minute * 30)
	for channel := range channels {
		visitUrl := fmt.Sprintf("%s://%s/%s/%d",etc.Config.Zxcs.Scheme, etc.Config.Zxcs.Address, etc.Config.Zxcs.RootPath, etc.Config.Zxcs.TypeList[channel].Value)
		log.Sugar.Infof("爬取页面 ===> %s", visitUrl)
		s := Spider{
			Url: visitUrl,
			Name: etc.Config.Zxcs.TypeList[channel].Name,
		}
		s.Run()
		<-limiter
	}
}

func init()  {
	SpideAll()
}