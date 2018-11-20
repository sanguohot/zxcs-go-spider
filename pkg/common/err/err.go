package util

import "errors"

var (
	ErrUrlRequired = errors.New("URL不能为空")
	ErrNovelNotExistInMap = errors.New("哈希map中索引不到小说")
	ErrInvalidVoteTypeNumber = errors.New("投票类型数量不正确")
)
