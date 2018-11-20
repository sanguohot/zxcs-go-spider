package spider

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sanguohot/zxcs-go-spider/etc"
	"github.com/sanguohot/zxcs-go-spider/pkg/common/file"
	"github.com/sanguohot/zxcs-go-spider/pkg/common/log"
	"io/ioutil"
	"os"
)

func init()  {
	if !file.FilePathExist(etc.GetSqliteDbDirPath()) {
		err := os.Mkdir(etc.GetSqliteDbDirPath(), os.ModePerm)
		if err != nil {
			log.Logger.Fatal(err.Error())
		}
	}
	if !file.FilePathExist(etc.GetSqliteDbFilePath()) {
		_, err := os.Create(etc.GetSqliteDbFilePath())
		if err != nil {
			log.Logger.Fatal(err.Error())
		}
	}
	// 确保sql可以重复执行 也就是包含IF NOT EXISTS
	if !file.FilePathExist(etc.GetSqliteTabelFilePath()) {
		log.Logger.Fatal(fmt.Errorf("sqlite: not found %s", etc.GetSqliteTabelFilePath()).Error())
	}
	db, err := sql.Open("sqlite3", etc.GetSqliteDbFilePath())
	defer db.Close()
	if err != nil {
		log.Logger.Fatal(err.Error())
	}
	data, err := ioutil.ReadFile(etc.GetSqliteTabelFilePath())
	if err != nil {
		log.Logger.Fatal(err.Error())
	}
	_, err = db.Exec(string(data))
	if err != nil {
		log.Logger.Fatal(err.Error())
	}
}

func SqliteSetNovelList(l []Novel) error {
	if l == nil || len(l) == 0 {
		return nil
	}
	log.Sugar.Debugf("写入数据库 %v", l)
	db, err := sql.Open("sqlite3", etc.GetSqliteDbFilePath())
	defer db.Close()
	if err != nil {
		return err
	}
	stmt, err := db.Prepare("INSERT  OR IGNORE INTO tbl_novel(NovelId, NovelHash, DownloadUrl, VoteUrl, NovelType, Size" +
		", Title, Detail, XianCao, LiangCao, GanCao, KuCao, DuCao, Time) values(?,?,?,?,?,?,?,?,?,?,?,?,?,?)")
	if err != nil {
		return err
	}
	var res sql.Result
	for _, n := range l {
		res, err = stmt.Exec(n.NovelId, n.NovelHash, n.DownloadUrl, n.VoteUrl, n.NovelType, n.Size, n.Title, n.Detail, n.XianCao, n.LiangCao,
			n.GanCao, n.KuCao, n.DuCao, n.Time)
		if err != nil {
			return err
		}
	}
	_, err = res.LastInsertId()
	if err != nil {
		return err
	}
	return nil
}