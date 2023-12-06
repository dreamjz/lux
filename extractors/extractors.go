package extractors

import (
	"net/url"
	"strings"
	"sync"

	"github.com/pkg/errors"

	"github.com/iawia002/lux/utils"
)

var lock sync.RWMutex
var extractorMap = make(map[string]Extractor)

// Register registers an Extractor.
func Register(domain string, e Extractor) {
	lock.Lock()
	extractorMap[domain] = e
	lock.Unlock()
}

// Extract is the main function to extract the data.
func Extract(u string, option Options) ([]*Data, error) { // 从 URL 中提取 资源信息
	u = strings.TrimSpace(u) // 去除首尾空格
	var domain string

	bilibiliShortLink := utils.MatchOneOf(u, `^(av|BV|ep)\w+`) // B 站短链
	if len(bilibiliShortLink) > 1 {
		bilibiliURL := map[string]string{
			"av": "https://www.bilibili.com/video/",
			"BV": "https://www.bilibili.com/video/",
			"ep": "https://www.bilibili.com/bangumi/play/",
		}
		domain = "bilibili"
		u = bilibiliURL[bilibiliShortLink[1]] + u
	} else {
		u, err := url.ParseRequestURI(u) // 解析 URL
		if err != nil {
			return nil, errors.WithStack(err)
		}
		if u.Host == "haokan.baidu.com" { // 百度好看
			domain = "haokan"
		} else {
			domain = utils.Domain(u.Host) // 获取域名
		}
	}
	extractor := extractorMap[domain] // 获取域名对应的 extractor
	if extractor == nil {
		extractor = extractorMap[""] // 通用 extractor
	}
	videos, err := extractor.Extract(u, option) // 解析 URL
	if err != nil {
		return nil, errors.WithStack(err)
	}
	for _, v := range videos {
		v.FillUpStreamsData() // 填充 资源 信息：大小，扩展名等
	}
	return videos, nil
}
