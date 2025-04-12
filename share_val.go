package main

import (
	"crypto/tls"
	"net/http"
	"net/url"
	"time"
)

// 主响应结构体
type MainResponse struct {
	Success    bool                   `json:"Success"`
	ErrMessage *string                `json:"ErrMessage"` // 使用指针处理null值
	Data       map[string]interface{} `json:"Data"`
}

type ChannelResponse struct {
	Success    bool          `json:"Success"`
	ErrMessage *string       `json:"ErrMessage"` // 使用指针处理null值
	Data       []ChannelData `json:"Data"`
}

// 频道数据结构体
type ChannelData struct {
	FnID              int      `json:"fnID"`
	FsTYPE            string   `json:"fsTYPE"`
	FsCDN_ROUTE       string   `json:"fsCDN_ROUTE"`
	FsTYPE_NAME       string   `json:"fsTYPE_NAME"`
	FsNAME            string   `json:"fsNAME"`
	Fs4GTV_ID         string   `json:"fs4GTV_ID"`
	FnCHANNEL_NO      int      `json:"fnCHANNEL_NO"`
	FsHEAD_FRAME      string   `json:"fsHEAD_FRAME"`
	FsPICTURE_LOGO    *string  `json:"fsPICTURE_LOGO"` // 使用指针处理null值
	FsLOGO_PC         string   `json:"fsLOGO_PC"`
	FsLOGO_MOBILE     string   `json:"fsLOGO_MOBILE"`
	FcOVERSEAS        bool     `json:"fcOVERSEAS"`
	LstEXCEPT_COUNTRY []string `json:"lstEXCEPT_COUNTRY"`
	LstALL_BITRATE    []string `json:"lstALL_BITRATE"` // null时会是nil切片
	FsFREE_PROFILE    *string  `json:"fsFREE_PROFILE"`
	LstSETs           []int    `json:"lstSETs"`
	FcFREE            bool     `json:"fcFREE"`
	FsQUALITY         string   `json:"fsQUALITY"`
	FnCHAT_ID         int      `json:"fnCHAT_ID"`
	FcHAS_PROGLIST    bool     `json:"fcHAS_PROGLIST"`
	FdLIKE_BEG_DATE   *string  `json:"fdLIKE_BEG_DATE"`
	FcLIKE            bool     `json:"fcLIKE"`
	FcRISTRICT        bool     `json:"fcRISTRICT"`
	FsEXPIRE_DATE     string   `json:"fsEXPIRE_DATE"`
	FnGRADED          int      `json:"fnGRADED"`
	FsVENDOR          string   `json:"fsVENDOR"`
	FsDESCRIPTION     string   `json:"fsDESCRIPTION"`
	FsSCHEDULE_DATE   *string  `json:"fsSCHEDULE_DATE"`
	FnORDER           int      `json:"fnORDER"`
}

var (
	DebugMode  = false
	PreCache   = false
	ValidToken = ""
	ProxyTs    = false
)
var ProxyUrlWhiteList = map[string]bool{
	//UrlLiTVServer: true,
	"ntd-tgc.cdn.hinet.net": true,
}

var (
	//CnResponse   ChannelResponse
	ChannelsData []ChannelData
)

var HeaderKey = "PyPJU25iI2IQCMWq7kblwh9sGCypqsxMp4sKjJo95SK43h08ff+j1nbWliTySSB+N67BnXrYv9DfwK+ue5wWkg=="

// key和iv位于apk文件根目录下的 koin.properties
const Key = "ilyB29ZdruuQjC45JhBBR7o2Z8WJ26Vg"
const IV = "JUMxvVMmszqUTeKn"

var M4gtvAuth = ""

const (
	UrlGetChannelUrl = "https://api2.4gtv.tv/App/GetChannelUrl2"
	UrlGetAppConfig  = "https://api2.4gtv.tv/App/GetAPPConfig"
	UrlGetAllChannel = "https://api2.4gtv.tv/Channel/GetAllChannel2/mobile"

	UrlLiTVServer = "https://ntd-tgc.cdn.hinet.net"
)

const FsVersion = "2.7.2"

var (
	FsEncKey   = ""
	FsValue    = ""
	FsUSER     = ""
	FsPASSWORD = ""
)

const UA = "okhttp/4.9.2"

var (
	RetryAttempts  = 3               //最大重试次数
	RetrySleepTime = time.Second * 2 //重试延迟时间
)

const (
	CapPacket      = false
	CapPacketProxy = "http://127.0.0.1:8888"
)

var Client = &http.Client{
	Timeout: 10 * time.Second,
	Transport: &http.Transport{
		Proxy: func(req *http.Request) (*url.URL, error) {
			if CapPacket {
				return url.Parse(CapPacketProxy)
			}
			return nil, nil // 不使用代理
		},
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
		ForceAttemptHTTP2:   false,
		MaxIdleConnsPerHost: 100,
		IdleConnTimeout:     30 * time.Second,
	},
}

var Client1 = &http.Client{
	Timeout: 10 * time.Second,
	Transport: &http.Transport{
		Proxy: func(req *http.Request) (*url.URL, error) {
			if CapPacket {
				return url.Parse(CapPacketProxy)
			}
			return nil, nil // 不使用代理
		},
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
		MaxIdleConnsPerHost: 100,
		IdleConnTimeout:     30 * time.Second,
		// 禁用 HTTP/2
		TLSNextProto: make(map[string]func(authority string, c *tls.Conn) http.RoundTripper),
	},
}

const LiTVDataFilePath = "./litv_data.json"

// Channel 定义频道结构
type Channel struct {
	AssetID   string
	Name      string
	Logo      string
	GroupName string
}

var LitvOnlychannels = []Channel{
	{
		AssetID:   "4gtv-4gtv010",
		Name:      "非凡新聞台",
		Logo:      "https://logo.doube.eu.org/非凡新闻.png",
		GroupName: "新聞財經",
	},
	{
		AssetID:   "4gtv-4gtv048",
		Name:      "非凡商業台",
		Logo:      "https://logo.doube.eu.org/非凡商业.png",
		GroupName: "新聞財經",
	},
	{
		AssetID:   "4gtv-4gtv051",
		Name:      "台視新聞台",
		Logo:      "https://logo.doube.eu.org/台视新闻台.png",
		GroupName: "新聞財經",
	},
	{
		AssetID:   "4gtv-4gtv056",
		Name:      "台視財經台",
		Logo:      "https://logo.doube.eu.org/台视财经台.png",
		GroupName: "新聞財經",
	},
	{
		AssetID:   "4gtv-4gtv066",
		Name:      "台視",
		Logo:      "https://logo.doube.eu.org/台视.png",
		GroupName: "綜合",
	},
	// {
	// 	ID:        "4gtv-4gtv072",
	// 	Name:      "TVBS新聞台",
	// 	Logo:      "https://logo.doube.eu.org/TVBS新闻.png",
	// 	GroupName: "新聞財經",
	// },
	{
		AssetID:   "4gtv-4gtv104",
		Name:      "第1商業台",
		Logo:      "https://logo.doube.eu.org/第1商业台.png",
		GroupName: "新聞財經",
	},
	{
		AssetID:   "4gtv-4gtv109",
		Name:      "中天亞洲台",
		Logo:      "https://logo.doube.eu.org/中天亚洲.png",
		GroupName: "綜合",
	},
	// {
	// 	ID:        "4gtv-4gtv155",
	// 	Name:      "民視",
	// 	Logo:      "https://logo.doube.eu.org/民视.png",
	// 	GroupName: "綜合",
	// },
	{
		AssetID:   "litv-longturn01",
		Name:      "龍華動畫台",
		Logo:      "https://logo.doube.eu.org/龙华卡通台.png",
		GroupName: "戲劇、電影與紀錄片",
	},
	{
		AssetID:   "litv-longturn02",
		Name:      "龍華洋片台",
		Logo:      "https://logo.doube.eu.org/龙华洋片台.png",
		GroupName: "戲劇、電影與紀錄片",
	},
	{
		AssetID:   "litv-longturn03",
		Name:      "龍華電影台",
		Logo:      "https://logo.doube.eu.org/龙华电影台.png",
		GroupName: "戲劇、電影與紀錄片",
	},
	{
		AssetID:   "litv-longturn04",
		Name:      "博斯魅力",
		Logo:      "https://logo.doube.eu.org/博斯魅力.png",
		GroupName: "運動休閒",
	},
	{
		AssetID:   "litv-longturn05",
		Name:      "博斯高球台",
		Logo:      "https://logo.doube.eu.org/博斯高球1.png",
		GroupName: "運動休閒",
	},
	{
		AssetID:   "litv-longturn06",
		Name:      "博斯高球二台",
		Logo:      "https://logo.doube.eu.org/博斯高球2.png",
		GroupName: "運動休閒",
	},
	{
		AssetID:   "litv-longturn07",
		Name:      "博斯運動台",
		Logo:      "https://logo.doube.eu.org/博斯运动1.png",
		GroupName: "運動休閒",
	},
	{
		AssetID:   "litv-longturn08",
		Name:      "博斯運動二台",
		Logo:      "https://logo.doube.eu.org/博斯运动2.png",
		GroupName: "運動休閒",
	},
	{
		AssetID:   "litv-longturn09",
		Name:      "博斯網球台",
		Logo:      "https://logo.doube.eu.org/博斯网球1.png",
		GroupName: "運動休閒",
	},
	{
		AssetID:   "litv-longturn10",
		Name:      "博斯無限台",
		Logo:      "https://logo.doube.eu.org/博斯无限1.png",
		GroupName: "運動休閒",
	},
	{
		AssetID:   "litv-longturn13",
		Name:      "博斯無限二台",
		Logo:      "https://logo.doube.eu.org/博斯无限2.png",
		GroupName: "運動休閒",
	},
	{
		AssetID:   "litv-longturn11",
		Name:      "龍華日韓台",
		Logo:      "https://logo.doube.eu.org/龙华日韩台.png",
		GroupName: "戲劇、電影與紀錄片",
	},
	{
		AssetID:   "litv-longturn12",
		Name:      "龍華偶像台",
		Logo:      "https://logo.doube.eu.org/龙华偶像台.png",
		GroupName: "戲劇、電影與紀錄片",
	},
	{
		AssetID:   "litv-longturn18",
		Name:      "龍華影劇台",
		Logo:      "https://logo.doube.eu.org/龙华戏剧台.png",
		GroupName: "戲劇、電影與紀錄片",
	},
	{
		AssetID:   "litv-longturn21",
		Name:      "龍華經典台",
		Logo:      "https://logo.doube.eu.org/龙华经典台.png",
		GroupName: "戲劇、電影與紀錄片",
	},
	{
		AssetID:   "litv-longturn22",
		Name:      "台灣戲劇台",
		Logo:      "https://logo.doube.eu.org/台湾戏剧台.png",
		GroupName: "戲劇、電影與紀錄片",
	},
	{
		AssetID:   "litv-longturn19",
		Name:      "Smart-知識台",
		Logo:      "https://logo.doube.eu.org/Smart知识台.png",
		GroupName: "戲劇、電影與紀錄片",
	},
}
