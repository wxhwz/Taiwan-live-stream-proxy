package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

// LitvChannelItem 定义结构体
type LitvChannelItem struct {
	Param1      int    `json:"param1"`
	Param2      int    `json:"param2"`
	CachePrefix string `json:"cacheprefix"`
}

var LitvChannelsMu sync.Mutex
var LitvChannels = map[string]LitvChannelItem{
	//新闻频道
	"4gtv-4gtv009":    {2, 7, ""},
	"4gtv-4gtv072":    {1, 2, ""},
	"4gtv-4gtv152":    {1, 6, ""},
	"litv-ftv13":      {1, 7, ""},
	"4gtv-4gtv075":    {1, 2, ""},
	"4gtv-4gtv010":    {1, 6, ""},
	"4gtv-4gtv051":    {1, 2, ""},
	"4gtv-4gtv052":    {1, 2, ""},
	"4gtv-4gtv074":    {1, 2, ""},
	"litv-longturn14": {1, 2, ""},
	"4gtv-4gtv156":    {1, 6, ""},
	"litv-ftv10":      {1, 7, ""},
	"litv-ftv03":      {1, 7, ""},

	//财经频道
	"4gtv-4gtv153": {1, 6, ""},
	"4gtv-4gtv048": {1, 2, ""},
	"4gtv-4gtv056": {1, 2, ""},
	"4gtv-4gtv104": {1, 7, ""},

	//综合频道
	"4gtv-4gtv073": {1, 2, ""},
	"4gtv-4gtv066": {1, 2, ""},
	"4gtv-4gtv040": {1, 6, ""},
	"4gtv-4gtv041": {1, 6, ""},
	"4gtv-4gtv002": {1, 10, ""},
	"4gtv-4gtv155": {1, 6, ""},
	"4gtv-4gtv001": {1, 6, ""},
	"4gtv-4gtv003": {1, 6, ""},
	"4gtv-4gtv109": {1, 7, ""},
	"4gtv-4gtv046": {1, 8, ""},
	"4gtv-4gtv063": {1, 6, ""},
	"4gtv-4gtv065": {1, 8, ""},
	"4gtv-4gtv043": {1, 6, ""},
	"4gtv-4gtv079": {1, 2, ""},
	"4gtv-4gtv084": {1, 6, ""},
	"4gtv-4gtv085": {1, 5, ""},

	//娱乐综艺频道
	"4gtv-4gtv068": {1, 7, ""},
	"4gtv-4gtv067": {1, 8, ""},
	"4gtv-4gtv070": {1, 7, ""},
	"4gtv-4gtv004": {1, 8, ""},
	"4gtv-4gtv039": {1, 7, ""},
	"4gtv-4gtv034": {1, 6, ""},
	"4gtv-4gtv054": {1, 8, ""},
	"4gtv-4gtv062": {1, 8, ""},
	"4gtv-4gtv064": {1, 8, ""},
	"4gtv-4gtv006": {1, 9, ""},

	//电影频道
	"4gtv-4gtv011":    {1, 6, ""},
	"4gtv-4gtv017":    {1, 6, ""},
	"4gtv-4gtv061":    {1, 7, ""},
	"4gtv-4gtv055":    {1, 8, ""},
	"4gtv-4gtv049":    {1, 8, ""},
	"litv-ftv09":      {1, 2, ""},
	"litv-longturn03": {5, 6, ""},
	"litv-longturn02": {5, 2, ""},

	//戏剧频道
	"4gtv-4gtv042":    {1, 6, ""},
	"4gtv-4gtv045":    {1, 6, ""},
	"4gtv-4gtv058":    {1, 8, ""},
	"4gtv-4gtv080":    {1, 6, ""},
	"4gtv-4gtv047":    {1, 8, ""},
	"litv-longturn18": {5, 6, ""},
	"litv-longturn11": {5, 2, ""},
	"litv-longturn12": {5, 2, ""},
	"litv-longturn21": {5, 2, ""},
	"litv-longturn22": {5, 2, ""},

	//体育频道
	"4gtv-4gtv014":    {1, 5, ""},
	"4gtv-4gtv053":    {1, 8, ""},
	"4gtv-4gtv101":    {1, 5, ""},
	"litv-longturn04": {5, 2, ""},
	"litv-longturn05": {5, 2, ""},
	"litv-longturn06": {5, 2, ""},
	"litv-longturn07": {5, 2, ""},
	"litv-longturn08": {5, 2, ""},
	"litv-longturn09": {5, 2, ""},
	"litv-longturn10": {5, 2, ""},
	"litv-longturn13": {4, 2, ""},

	//纪实/知识/旅游频道
	"4gtv-4gtv013":    {1, 6, ""},
	"4gtv-4gtv016":    {1, 6, ""},
	"4gtv-4gtv018":    {1, 6, ""},
	"4gtv-4gtv076":    {1, 2, ""},
	"litv-ftv07":      {1, 7, ""},
	"litv-longturn19": {5, 2, ""},

	//儿童/卡通频道
	"4gtv-4gtv044":    {1, 8, ""},
	"4gtv-4gtv057":    {1, 8, ""},
	"litv-longturn01": {4, 2, ""},

	//音乐/艺术频道
	"4gtv-4gtv059": {1, 6, ""},
	"4gtv-4gtv082": {1, 6, ""},
	"4gtv-4gtv083": {1, 6, ""},

	//教育/宗教频道
	"litv-ftv16":      {1, 2, ""},
	"litv-ftv17":      {1, 2, ""},
	"litv-longturn20": {5, 6, ""},
}

// LitvSaveToFile 将 LitvChannels 保存到文件
func LitvSaveToFile(filename string) error {
	data, err := json.MarshalIndent(LitvChannels, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

// LitvLoadFromFile 从文件加载 LitvChannels
func LitvLoadFromFile(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	LitvChannels = map[string]LitvChannelItem{}
	return json.Unmarshal(data, &LitvChannels)
}

func LitvGenerateM3U8(id string) []byte { //经测试，Litv最多回放10分钟，懒得写回看代码
	ch := LitvChannels[id]

	timestamp := int(time.Now().Unix()/4 - 355017625)
	t := timestamp * 4

	var buf bytes.Buffer
	buf.WriteString("#EXTM3U\n")
	buf.WriteString("#EXT-X-VERSION:3\n")
	buf.WriteString("#EXT-X-TARGETDURATION:4\n")
	buf.WriteString(fmt.Sprintf("#EXT-X-MEDIA-SEQUENCE:%d\n", timestamp))

	if ch.CachePrefix == "" {
		for i := 0; i < 3; i++ {
			// fmt.Printf("https://ntd-tgc.cdn.hinet.net/live/pool/%s/litv-pc/%s-avc1_6000000=%d-mp4a_134000_zho=%d-begin=%d0000000-dur=40000000-seq=%d.ts\n",
			// 	id,
			// 	id,
			// 	ch.param1,
			// 	ch.param2,
			// 	t,
			// 	timestamp,
			// )
			buf.WriteString("#EXTINF:4,\n")
			buf.WriteString(fmt.Sprintf("%s/live/pool/%s/litv-pc/%s-avc1_6000000=%d-mp4a_134000_zho=%d-begin=%d0000000-dur=40000000-seq=%d.ts\n",
				UrlLiTVServer,
				id,
				id,
				ch.Param1,
				ch.Param2,
				t,
				timestamp,
			))
			timestamp++
			t += 4
		}
	} else {
		for i := 0; i < 3; i++ {
			buf.WriteString("#EXTINF:4,\n")
			buf.WriteString(fmt.Sprintf("%s/live/pool/%s/litv-pc/%s-begin=%d0000000-dur=40000000-seq=%d.ts\n",
				UrlLiTVServer,
				id,
				ch.CachePrefix,
				t,
				timestamp,
			))
			timestamp++
			t += 4
		}
	}

	return buf.Bytes()
}
func LitvUpdateChannel(assetID, playUrl string) bool {
	newURL, err := HandleM3u8Raw(playUrl, "url")
	if err != nil {
		LogError(err)
		return false
	}
	part2 := newURL[strings.LastIndex(newURL, "/")+1:]
	cachePrefix := strings.Split(part2, ".")[0]
	cachePrefix = strings.ReplaceAll(cachePrefix, "video=2000000", "video=6000000")
	cachePrefix = strings.ReplaceAll(cachePrefix, "video=2936000", "video=5936000")
	cachePrefix = strings.ReplaceAll(cachePrefix, "video=3000000", "video=6000000")
	cachePrefix = strings.ReplaceAll(cachePrefix, "avc1_2000000=3", "avc1_6000000=1")
	cachePrefix = strings.ReplaceAll(cachePrefix, "avc1_2000000=6", "avc1_6000000=1")
	cachePrefix = strings.ReplaceAll(cachePrefix, "avc1_2936000=4", "avc1_6000000=5")
	cachePrefix = strings.ReplaceAll(cachePrefix, "avc1_3000000=3", "avc1_6000000=1")

	if cachePrefix != "" && cachePrefix != "index" {
		LogDebug(assetID, " cachePrefix:", cachePrefix)
		LitvChannelsMu.Lock()
		LitvChannels[assetID] = LitvChannelItem{
			CachePrefix: cachePrefix,
		}
		LitvChannelsMu.Unlock()
	} else {
		LogError("无法解析url")
		return false
	}
	return true

	// re := regexp.MustCompile(`=(\d+)-.+=(\d+)\.m3u8`)
	// matches := re.FindStringSubmatch(newURL)

	// if len(matches) >= 3 {
	// 	LogDebug(assetID, " param1:", matches[1], " param2:", matches[2])
	// 	param1, err := strconv.Atoi(matches[1])
	// 	if err != nil {
	// 		LogError(err)
	// 		return false
	// 	}
	// 	param2, err := strconv.Atoi(matches[2])
	// 	if err != nil {
	// 		LogError(err)
	// 		return false
	// 	}
	// 	LitvChannelsMu.Lock()
	// 	LitvChannels[assetID] = LitvChannelItem{
	// 		Param1: param1,
	// 		Param2: param2,
	// 	}
	// 	LitvChannelsMu.Unlock()
	// } else {
	// 	LogError(assetID, "No matches found, playurl:", playUrl)
	// 	return false
	// }
	// return true
}

func LitvInitChannelsFromFile() {
	e := LitvLoadFromFile(LiTVDataFilePath)
	if e != nil {
		LogError("读取LitvChannels失败", LiTVDataFilePath, e)
		e = LitvSaveToFile(LiTVDataFilePath)
		if e != nil {
			LogError("写入LitvChannels到文件", LiTVDataFilePath, "失败", e)
		}
		LogInfo("写入LitvChannels到文件", LiTVDataFilePath, "成功")
	} else {
		LogInfo("读取LitvChannels成功", LiTVDataFilePath)
	}
}
