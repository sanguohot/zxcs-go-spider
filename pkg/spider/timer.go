package spider

import (
	"github.com/sanguohot/zxcs-go-spider/pkg/common/log"
	"time"
)

var step = 20

func timerTask()  {
	cnt := 0
	l := []Novel{}
	for k, v := range novelMap {
		if v.Time > 0 {
			cnt++
			if cnt > step {
				break
			}
			//dst := new(Novel)
			//if err := deepcopy.Copy(dst, v); err != nil {
			//	log.Logger.Error(err.Error())
			//}
			//l = append(l, *dst)
			// 这里无需克隆
			l = append(l, v)
			delete(novelMap, k)
		}
	}

	err := SqliteSetNovelList(l)
	if err != nil {
		log.Logger.Error(err.Error())
	}
}

func init() {
	// 一分钟写一百条
	ticks := time.NewTicker(1 * time.Minute)
	tick := ticks.C
	go func() {
		for _ = range tick {
			timerTask()
		}
	}()
}