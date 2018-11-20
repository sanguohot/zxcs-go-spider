package spider

type Novel struct {
	NovelId int `json:"novel_id"`
	NovelHash string `json:"novel_hash"`
	DownloadUrl      string `json:"download_url"`
	VoteUrl      string `json:"vote_url"`
	NovelType      string `json:"novel_type"`
	Size      int    `json:"size"`
	Title     string `json:"title"`
	Detail    string `json:"detail"`
	XianCao   int    `json:"xian_cao"`
	LiangCao  int    `json:"liang_cao"`
	GanCao    int    `json:"gan_cao"`
	KuCao     int    `json:"ku_cao"`
	DuCao     int    `json:"du_cao"`
	Time      int64    `json:"time"`
}