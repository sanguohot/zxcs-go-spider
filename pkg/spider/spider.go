package spider

import (
	"fmt"
	"github.com/gocolly/colly"
	"github.com/sanguohot/zxcs-go-spider/etc"
	"github.com/sanguohot/zxcs-go-spider/pkg/common/client"
	selfErr "github.com/sanguohot/zxcs-go-spider/pkg/common/err"
	"github.com/sanguohot/zxcs-go-spider/pkg/common/file"
	"github.com/sanguohot/zxcs-go-spider/pkg/common/hash"
	"github.com/sanguohot/zxcs-go-spider/pkg/common/log"
	"github.com/sanguohot/zxcs-go-spider/pkg/common/random"
	"go.uber.org/zap"
	"io/ioutil"
	"path"
	"strconv"
	"strings"
	"time"
)

var (
	defaultMaxDepth = 1
	defaultConcurrence = 1
	defaultMatchSelector = "a[href$='.rar']"
	defaultOutput = "./data"
	downloadTip = "线路一"
	novelMap = make(map[int]Novel)
)

type Spider struct {
	Url string
	Name string
	Collect *colly.Collector
}

func (s *Spider) Init() {
	//output := defaultOutput
	maxDepth := defaultMaxDepth
	concurrence := defaultConcurrence
	// Instantiate default collector
	c := colly.NewCollector(
		// MaxDepth is 1, so only the links on the scraped page
		// is visited, and no further links are followed
		colly.MaxDepth(maxDepth),
		colly.Async(true),
	)
	//c.AllowedDomains = []string{"zxcs8.com", "zxcs.me", "zxcs1.xyz"}
	c.Limit(&colly.LimitRule{
		DomainGlob: "*",
		Parallelism: concurrence,
		RandomDelay: time.Second * 10,
	})
	s.Collect = c
	s.ErrorHandler()
}

func (s *Spider) Visit()  {
	s.Collect.Visit(s.Url)
}

func (s *Spider) Wait()  {
	s.Collect.Wait()
}

func (s *Spider) ErrorHandler()  {
	// 异常处理
	s.Collect.OnError(func(r *colly.Response, err error) {
		log.Logger.Error(err.Error())
	})
}

func (s *Spider) GetIndexAndNumberFromUrlEnd(url string, split string) (int, int, error) {
	index := strings.LastIndex(url, split)
	if index < 0 {
		return index, 0, nil
	}
	max, err := strconv.Atoi(url[index+len(split):])
	return index, max, err
}

func (s *Spider) GetNextPageAndVisit()  {
	// 寻找下一页并访问
	s.Collect.OnHTML("div#pagenavi", func(e *colly.HTMLElement) {
		curPage, err := strconv.Atoi(e.ChildText("span"))
		if err != nil {
			log.Logger.Error(err.Error())
			return
		}
		e.ForEachWithBreak("a", func(i int, element *colly.HTMLElement) bool {
			pageUrl := element.Attr("href")
			split := "/page/"
			index, page, err := s.GetIndexAndNumberFromUrlEnd(pageUrl, split)
			if err != nil {
				log.Logger.Error(err.Error(), zap.String("pageUrl", pageUrl))
				return true
			}
			if index < 0 || page <= 0 {
				return true
			}
			if page > curPage{
				nextUrl := fmt.Sprintf("%s%s%d", pageUrl[0:index], split, page)
				log.Sugar.Infof("当前页 ===> %d, 下一页 ===> %d, 延时访问 ===> %s", curPage, page, nextUrl)
				s.Collect.Visit(nextUrl)
				return false
			}
			return true
		})
	})
}

func (s *Spider) GetRealDownloadUrlAndStartDownload()  {
	// 获取真实下载地址
	s.Collect.OnHTML(defaultMatchSelector, func(e *colly.HTMLElement) {
		if e.Text == downloadTip {
			//log.Sugar.Debugf("开始下载小说 %s", e.Attr("href"))
			// http://www.zxcs.me/download.php?id=11179
			var (
				data []byte
				err error
			)
			_, id, err := s.GetIndexAndNumberFromUrlEnd(e.Request.URL.String(), "?id=")
			if err != nil {
				log.Logger.Error(err.Error())
				return
			}
			novel, ok := novelMap[id]
			if !ok {
				log.Logger.Error(selfErr.ErrNovelNotExistInMap.Error(), zap.Int("id", id))
				return
			}
			// http://d5.zxcs1.xyz/201811/dzckddm,dbljs.rar
			realDownloadUrl := e.Attr("href")
			novel.DownloadUrl = realDownloadUrl
			if !file.IsFileExist(path.Join(etc.GetServerDir(), defaultOutput), strconv.Itoa(id)) {
				data, err = client.Download(realDownloadUrl)
				if err != nil {
					log.Logger.Error(err.Error())
					return
				}
				err = file.SaveToLocal(path.Join(etc.GetServerDir(), defaultOutput), strconv.Itoa(id), data)
				if err != nil {
					log.Logger.Error(err.Error())
					return
				}
			}else {
				data, err = ioutil.ReadFile(path.Join(etc.GetServerDir(), defaultOutput, strconv.Itoa(id)))
				if err != nil {
					log.Logger.Error(err.Error())
					return
				}
			}
			novel.NovelHash = hash.Sha256Hash(data).Hex()
			novel.Size = len(data)
			novel.Time = time.Now().Unix()
			novelMap[id] = novel
			log.Sugar.Infof("下载完毕 title=%s id=%d", novel.Title, id)
		}
	})
}

func (s *Spider) GetNovelInfoAndPreDownload()  {
	// 获取小说基本信息并进行请求
	s.Collect.OnHTML("#plist", func(e *colly.HTMLElement) {
		title := e.ChildText("dt")
		novelUrl := e.ChildAttr("dt a", "href")
		novelDesc := e.ChildText("dd.des")
		_, id, err := s.GetIndexAndNumberFromUrlEnd(novelUrl, "/")
		if err != nil {
			log.Logger.Error(err.Error())
			return
		}
		_, ok := novelMap[id]
		if ok {
			log.Sugar.Warnf("小说 %d 已经下载，跳过", id)
			return
		}
		novelMap[id] = Novel{
			NovelId: id,
			Title: title,
			Detail: novelDesc,
			NovelType: s.Name,
		}
		s.GetVote(id)
		downloadUrl := fmt.Sprintf("%s://%s/%s?id=%d",etc.Config.Zxcs.Scheme, etc.Config.Zxcs.DownloadAddress, etc.Config.Zxcs.DownloadPath, id)
		log.Sugar.Debugf("延时下载 ===> title=%s, id=%d", title, id)
		s.Collect.Visit(downloadUrl)
	})
}

func (s *Spider) GetVote(id int)  {
	// 获取投票, 仙草、粮草、干草、枯草、毒草
	novel, ok := novelMap[id]
	if  !ok {
		log.Logger.Error(selfErr.ErrNovelNotExistInMap.Error())
		return
	}
	voteUrl := fmt.Sprintf("%s://%s/%s&id=%d&m=%16f",etc.Config.Zxcs.Scheme, etc.Config.Zxcs.Address, etc.Config.Zxcs.VotePath, id, random.GenerateRandomFloat())
	data, err := client.Download(voteUrl)
	if err != nil {
		log.Logger.Error(err.Error())
		return
	}
	//log.Sugar.Debugf("%d ===> %s", id, string(data))
	novel.VoteUrl = voteUrl
	list := strings.Split(string(data), ",")
	if len(list) != 5 {
		log.Logger.Error(selfErr.ErrInvalidVoteTypeNumber.Error(), zap.String("vote", string(data)))
		return
	}
	for _, v := range list{
		if _, err := strconv.Atoi(v); err != nil {
			log.Logger.Error(err.Error())
			return
		}
	}
	novel.XianCao, _ = strconv.Atoi(list[0])
	novel.LiangCao, _ = strconv.Atoi(list[1])
	novel.GanCao, _ = strconv.Atoi(list[2])
	novel.KuCao, _ = strconv.Atoi(list[3])
	novel.DuCao, _ = strconv.Atoi(list[4])
	novelMap[id] = novel
}

func (s *Spider) Run() error {
	if s.Url == "" {
		log.Logger.Fatal(selfErr.ErrUrlRequired.Error())
	}
	s.Init()

	// 注意这里必须是顺序执行
	s.GetNextPageAndVisit()
	s.GetNovelInfoAndPreDownload()
	s.GetRealDownloadUrlAndStartDownload()

	s.Visit()
	s.Wait()
	return nil
}