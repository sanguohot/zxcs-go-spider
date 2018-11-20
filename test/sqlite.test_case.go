package main

import (
	"github.com/sanguohot/zxcs-go-spider/pkg/common/log"
	"github.com/sanguohot/zxcs-go-spider/pkg/spider"
)

func main()  {
	l := []spider.Novel{spider.Novel{
		NovelId: 1,
		NovelHash: "1",
		DownloadUrl: "2",
		VoteUrl: "3",
		NovelType: "4",
		Size: 666,
		Title: "a",
		Detail: "b",
		XianCao: 11,
		LiangCao: 12,
		GanCao: 13,
		KuCao: 14,
		DuCao: 15,
		Time: 16,
	}}
	if err := spider.SqliteSetNovelList(l); err != nil {
		log.Logger.Fatal(err.Error())
	}
}