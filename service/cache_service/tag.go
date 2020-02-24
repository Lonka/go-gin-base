package cache_service

import (
	"go_gin_base/pkg/e"
	"strings"

	"github.com/unknwon/com"
)

type Tag struct {
	ID       int
	Name     string
	State    int
	PageNum  int
	PageSize int
}

func (t *Tag) GetTagsKey() string {
	keys := []string{
		e.CACHE_TAG,
		"LIST",
	}
	t.GetKeys(&keys)
	return strings.Join(keys, "_")
}

func (t *Tag) CountTagsKey() string {
	keys := []string{
		e.CACHE_TAG,
		"COUNT",
	}
	t.GetKeys(&keys)
	return strings.Join(keys, "_")
}

func (t *Tag) GetKeys(keys *[]string) {
	if t.Name != "" {
		*keys = append(*keys, t.Name)
	}
	if t.State >= 0 {
		*keys = append(*keys, com.ToStr(t.State))
	}
	if t.PageNum >= 0 {
		*keys = append(*keys, com.ToStr(t.PageNum))
	}
	if t.PageSize >= 0 {
		*keys = append(*keys, com.ToStr(t.PageSize))
	}

}

func GetTagsKey() string {
	keys := []string{
		e.CACHE_TAG,
	}
	return strings.Join(keys, "_") + "*"
}
