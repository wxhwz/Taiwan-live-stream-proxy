package main

import (
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type M4gtv struct {
}

var playUrlCache sync.Map

type PlayUrlCacheItem struct {
	playUrl    string
	Expiration int64
}

func (y *M4gtv) HandleMainRequest(c *gin.Context, channelID, assetID, cdnType string) {
	var (
		playUrl string
		found   bool
	)
	if cdnType == "A" || (channelID == "" && cdnType == "") {
		_, exists := LitvChannels[assetID]
		if exists {
			m3u8Raw := LitvGenerateM3U8(assetID)
			if !ProxyTs {
				c.Data(http.StatusOK, "application/vnd.apple.mpegurl", m3u8Raw)
				return
			} else {
				fmt.Println("ttttttttttt")
				m3u8Content := ReplaceM3u8Data(string(m3u8Raw), "http://"+c.Request.Host+c.Request.URL.Path+"?ts=")
				if m3u8Content == "" {
					LogError()
					c.String(http.StatusNotFound, "m3u8Content 为空")
					return
				}
				c.Header("Content-Type", "application/vnd.apple.mpegurl")
				c.String(200, m3u8Content)
				return
			}
		}
		goto NEXT
	}

	playUrl, found = getPlayUrlCache(assetID)

	if found {
		if !ProxyTs {
			c.Redirect(302, playUrl)
			return
		} else {
			m3u8Content := ReplaceM3u8Data2(playUrl, "http://"+c.Request.Host+c.Request.URL.Path)
			if m3u8Content == "" {
				LogError()
				c.String(http.StatusNotFound, "m3u8Content 为空")
				return
			}
			c.Header("Content-Type", "application/vnd.apple.mpegurl")
			c.String(200, m3u8Content)
			return
		}
	}
NEXT:
	playUrl, _ = GetPlayUrl(channelID, assetID)
	if playUrl == "" {
		c.String(http.StatusNotFound, "playUrl 为空")
		return
	}
	if cdnType == "A" {
		if !LitvUpdateChannel(assetID, playUrl) {
			c.String(http.StatusNotFound, "UpdateLitvChannel 失败")
			return
		}
		m3u8Raw := LitvGenerateM3U8(assetID)
		if !ProxyTs {
			c.Data(http.StatusOK, "application/vnd.apple.mpegurl", m3u8Raw)
			return
		} else {
			m3u8Content := ReplaceM3u8Data(string(m3u8Raw), "http://"+c.Request.Host+c.Request.URL.Path+"?ts=")
			if m3u8Content == "" {
				LogError()
				c.String(http.StatusNotFound, "m3u8Content 为空")
				return
			}
			c.Header("Content-Type", "application/vnd.apple.mpegurl")
			c.String(200, m3u8Content)
			return
		}

	} else if cdnType == "B" {
		playUrl = strings.Replace(playUrl, "/index.m3u8?", "/1080.m3u8?", 1)
	}
	setPlayUrlCache(assetID, playUrl)

	if !ProxyTs {
		c.Redirect(302, playUrl)
	} else {
		m3u8Content := ReplaceM3u8Data2(playUrl, "http://"+c.Request.Host+c.Request.URL.Path)
		if m3u8Content == "" {
			LogError()
			c.String(http.StatusNotFound, "m3u8Content 为空")
			return
		}
		c.Header("Content-Type", "application/vnd.apple.mpegurl")
		c.String(200, m3u8Content)
	}

}
func (y *M4gtv) HandleTsRequest(c *gin.Context, tsUrl string) {

	parsedURL, err := url.Parse(tsUrl)
	if err != nil {
		LogError(err)
		c.String(404, "", err)
		return
	}
	_, exist := ProxyUrlWhiteList[parsedURL.Host]
	if !exist {
		LogError("未知的ts", tsUrl)
		c.String(404, "未知的ts")
	}

	statusCode, respBody, err := MRequest(tsUrl, "GET", nil,
		map[string]string{
			"Referer":    "https://imasdk.googleapis.com",
			"origin":     "https://imasdk.googleapis.com",
			"User-Agent": "Mozilla/5.0 (Linux; Android 12; M2011K2C ; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/134.0.6998.135 Mobile Safari/537.36",
			//"Accept-Encoding": "identity",
		}, false)
	if err != nil || respBody == "" {
		LogError("HandleTsRequest 失败：", err)
		c.String(statusCode, respBody)
		return
	}
	// 设置响应头为视频流类型
	c.Header("Content-Type", "video/MP2T")
	// 返回视频数据
	c.String(statusCode, respBody)
}

func HandleM3u8Raw(m3u8Url string, returnType string) (string, error) {

	_, respBody, err := MRequest(m3u8Url, "GET", nil,
		map[string]string{
			"Referer":         "https://imasdk.googleapis.com",
			"origin":          "https://imasdk.googleapis.com",
			"User-Agent":      "Mozilla/5.0 (Linux; Android 12;) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/134.0.6998.135 Mobile Safari/537.36",
			"Accept-Encoding": "gzip",
		}, false)
	if err != nil || respBody == "" {
		LogError("HandleM3u8Raw 失败：", err)
		return "", err
	}
	parsedURL, err := url.Parse(m3u8Url)
	if err != nil {
		LogError("HandleM3u8Raw 失败：", err)
		return "", err
	}
	urlPath := path.Dir(parsedURL.Path)
	latestLine := getLastLineFromString(respBody)
	newURL := fmt.Sprintf("%s://%s%s/%s", parsedURL.Scheme, parsedURL.Host, urlPath, latestLine)

	hasM3u8 := strings.Contains(respBody, ".m3u8")

	if returnType == "url" {
		if hasM3u8 {
			return newURL, nil
		} else {
			return m3u8Url, nil
		}
	}
	if !hasM3u8 {
		return respBody, nil
	}
	return HandleM3u8Raw(newURL, returnType)
}

func ReplaceM3u8Data2(playUrl, sourceUrlPath string) string {
	lastSlash := strings.LastIndex(playUrl, "/")
	var playUrlPath string
	if lastSlash != -1 {
		playUrlPath = playUrl[:lastSlash+1]
		fmt.Println(playUrlPath)
	} else {
		LogError()
		return ""
	}
	playUrlPath = sourceUrlPath + "?ts=" + playUrlPath

	m3u8Content, err := HandleM3u8Raw(playUrl, "raw")
	if err != nil {
		LogError(err)
		return ""
	}
	return ReplaceM3u8Data(m3u8Content, playUrlPath)
}
func ReplaceM3u8Data(m3u8Content, prefix string) string {
	// 按行分割 m3u8 内容
	lines := strings.Split(m3u8Content, "\n")
	var builder strings.Builder

	// 逐行处理
	for i, line := range lines {
		// 不是第一行时添加换行符
		if i > 0 {
			builder.WriteString("\n")
		}

		// 跳过空行和以 # 开头的行（m3u8 元数据）
		if strings.TrimSpace(line) == "" || strings.HasPrefix(line, "#") {
			builder.WriteString(line)
			continue
		}

		// // 跳过已经是绝对路径的行（以 http:// 或 https:// 开头）
		// if strings.HasPrefix(line, "http://") || strings.HasPrefix(line, "https://") {
		// 	builder.WriteString(line)
		// 	continue
		// }

		// 添加 URL 前缀
		builder.WriteString(prefix)
		builder.WriteString(line)
	}
	return builder.String()
}

func getLastLineFromString(s string) string {
	if len(s) == 0 {
		return ""
	}

	lastNewLine := strings.LastIndex(s, "\n")
	if lastNewLine == -1 {
		return s // 没有换行符，返回整个字符串
	}

	// 如果最后一个字符是换行符，往前找
	if lastNewLine == len(s)-1 && lastNewLine > 0 {
		s = s[:lastNewLine]
		lastNewLine = strings.LastIndex(s, "\n")
		if lastNewLine == -1 {
			return s
		}
	}

	return s[lastNewLine+1:]
}

// 从缓存中获取数据
func getPlayUrlCache(key string) (string, bool) {
	// 查找缓存
	if item, found := playUrlCache.Load(key); found {
		cacheItem := item.(PlayUrlCacheItem)
		// 检查缓存是否过期
		if time.Now().Unix() < cacheItem.Expiration {
			return cacheItem.playUrl, true
		}
	}
	// 如果没有找到或缓存已过期，返回空
	return "", false
}

func setPlayUrlCache(key, playUrl string) {
	playUrlCache.Store(key, PlayUrlCacheItem{
		playUrl:    playUrl,
		Expiration: time.Now().Unix() + 10800,
	})
}

func RefreshCache(shouldCache bool, shouldUpdate bool) {
	LogInfo("刷新缓存中,shouldCache ", shouldCache, ",shouldUpdate ", shouldUpdate)

	if !shouldCache && !shouldUpdate {
		return
	}
	for _, channelData := range ChannelsData {
		time.Sleep(10 * time.Second)
		switch channelData.FsCDN_ROUTE {
		case "A":
			if shouldUpdate {
				playUrl, _ := GetPlayUrl(strconv.Itoa(channelData.FnID), channelData.Fs4GTV_ID)
				if playUrl == "" {
					continue
				}
				LitvUpdateChannel(channelData.Fs4GTV_ID, playUrl)
			}
		case "B":
			if shouldCache {
				playUrl, _ := GetPlayUrl(strconv.Itoa(channelData.FnID), channelData.Fs4GTV_ID)
				if playUrl == "" {
					continue
				}
				playUrl = strings.Replace(playUrl, "/index.m3u8?", "/1080.m3u8?", 1)
				setPlayUrlCache(channelData.Fs4GTV_ID, playUrl)
			}
		default:
			if shouldCache {
				playUrl, _ := GetPlayUrl(strconv.Itoa(channelData.FnID), channelData.Fs4GTV_ID)
				if playUrl == "" {
					continue
				}
				setPlayUrlCache(channelData.Fs4GTV_ID, playUrl)
			}
		}
	}
	if shouldUpdate {
		e := LitvSaveToFile(LiTVDataFilePath)
		if e != nil {
			LogError("更新LitvChannels到文件", LiTVDataFilePath, "失败", e)
		} else {
			LogInfo("更新LitvChannels到文件", LiTVDataFilePath, "成功")
		}

		GetAllChannels()
	}
	LogInfo("缓存刷新完毕")
}
