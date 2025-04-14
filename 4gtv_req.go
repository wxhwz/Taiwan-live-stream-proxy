package main

import (
	"bytes"
	"compress/gzip"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func Get4gtvAuth() {
	// 获取当前时间的format字符串（GMT格式化为yyyyMMdd）
	formatStr := getFormatString()
	LogDebug("当前时间的format：", formatStr)

	GetHeaderKey()

	// Base64解码headKey
	decodedHeaderKey, err := base64.StdEncoding.DecodeString(HeaderKey)
	if err != nil {
		LogError("Get4gtvAuth 失败：", err)
		return
	}
	decrypted, e := AESDecryptCBC(decodedHeaderKey, []byte(Key), []byte(IV))
	if e != nil {
		LogError("Get4gtvAuth 失败：", e)
		return
	}
	// 转换为UTF-8字符串
	decryptedStr := string(decrypted)
	LogDebug("decryptedStr: ", decryptedStr)
	// 拼接format和decrypted字符串
	combined := formatStr + decryptedStr
	// 计算SHA-512哈希并Base64编码
	hash := sha512.Sum512([]byte(combined))
	M4gtvAuth = base64.StdEncoding.EncodeToString(hash[:])
	LogInfo("M4gtvAuth: ", M4gtvAuth)

}

func GetHeaderKey() {
	_, respBody, err := MRequest(UrlGetAppConfig, "POST",
		map[string]interface{}{
			"fsDEVICE":  "Android",
			"fsVERSION": FsVersion,
		},
		map[string]string{
			"Content-Type":    "application/json; charset=UTF-8",
			"4gtv_auth":       M4gtvAuth,
			"fsdevice":        "Android",
			"fsversion":       FsVersion,
			"fsenc_key":       FsEncKey,
			"fsvalue":         FsValue,
			"Accept-Encoding": "gzip",
			"User-Agent":      UA,
		}, false)
	if err != nil {
		LogError("GetHeaderKey 失败：", err)
		return
	}

	var result map[string]interface{}
	e2 := json.Unmarshal([]byte(respBody), &result)
	if e2 != nil {
		LogError("GetHeaderKey 失败：", e2)
		return
	}
	if !result["Success"].(bool) {
		LogError("GetHeaderKey 失败：", result["ErrMessage"])
		return
	}
	data := result["Data"].(map[string]interface{})
	HeaderKey = data["header_key"].(string)
	LogInfo("HeaderKey ", HeaderKey)

}
func GetAllChannels() {
	_, respBody, err := MRequest(UrlGetAllChannel, "GET", nil,
		map[string]string{
			"Content-Type":    "application/json; charset=UTF-8",
			"4gtv_auth":       M4gtvAuth,
			"fsdevice":        "Android",
			"fsversion":       FsVersion,
			"fsenc_key":       FsEncKey,
			"fsvalue":         FsValue,
			"Accept-Encoding": "gzip",
			"User-Agent":      UA,
		}, false)
	if err != nil {
		LogError("GetAllChannels 失败：", err)
		return
	}
	var cnResponse ChannelResponse
	err = json.Unmarshal([]byte(respBody), &cnResponse)
	if err != nil {
		LogError("GetAllChannels 失败：", err)
		return
	}

	// 访问数据示例
	if cnResponse.Success {
		ChannelsData = cnResponse.Data
		LogInfo("请求成功，共获取频道数:", len(ChannelsData))

		// for _, channel := range ChannelsData {
		// 	if channel.FsCDN_ROUTE != "A" {
		// 		continue
		// 	}
		// 	fmt.Printf("\n频道名称: %s\n", channel.FsNAME)
		// 	fmt.Printf("FNID: %d\n", channel.FnID)
		// 	fmt.Printf("频道ID: %s\n", channel.Fs4GTV_ID)
		// 	//fmt.Printf("画质: %s\n", channel.FsQUALITY)
		// 	fmt.Printf("频道CDN: %s\n", channel.FsCDN_ROUTE)

		// 	// // 处理可能为nil的字段
		// 	// if channel.FsPICTURE_LOGO != nil {
		// 	// 	fmt.Printf("图片Logo: %s\n", *channel.FsPICTURE_LOGO)
		// 	// }
		//  }
	} else if cnResponse.ErrMessage != nil {
		LogError("GetAllChannels 请求失败：", *cnResponse.ErrMessage)
		return
	}
	BuildChannelMap()

}

// func PushDevice() {
// 	url := "https://service.4gtv.tv/4gtv/Data/PushDevice.ashx?DeviceID=356f0957-3aed-4d91-a177-48c14605bd58&Token=&Platform=Android&Account=wxh086573@gmail.com&Desc=2.7.2"
// 	_, respBody, err := SendHttpRequest(url, "GET", nil,
// 		map[string]string{
// 			"Accept-Encoding": "gzip",
// 			"User-Agent":      "Dalvik/2.1.0 (Linux; U; Android 12; M2011K2C Build/SKQ1.211230.001)",
// 		}, false)
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println(respBody)
// }

func SignIn() {
	url := "https://api2.4gtv.tv/AppAccount/SignIn"
	_, respBody, err := MRequest(url, "POST",
		map[string]interface{}{
			"fsUSER":     FsUSER,
			"fsPASSWORD": FsPASSWORD,
			"fsENC_KEY":  FsEncKey,
		},
		map[string]string{
			"Content-Type":    "application/json; charset=UTF-8",
			"4gtv_auth":       M4gtvAuth,
			"fsdevice":        "Android",
			"fsversion":       FsVersion,
			"fsenc_key":       FsEncKey,
			"fsvalue":         FsValue,
			"Accept-Encoding": "gzip",
			"User-Agent":      UA,
		}, false)
	if err != nil {
		LogError("SignIn 失败：", err)
		return
	}
	LogDebug("SignIn 结果：", respBody)
	var result map[string]interface{}
	e2 := json.Unmarshal([]byte(respBody), &result)
	if e2 != nil {
		LogError("SignIn 失败：", e2)
		return
	}
	if !result["Success"].(bool) {
		LogError("SignIn 失败：ErrMessage:", result["ErrMessage"])
		return
	}
	FsValue = result["Data"].(string)
	LogInfo("登录成功！")
	LogInfo("FsValue：", FsValue)
}

func GetPlayUrl(channelID, assetID string) (string, error) {
	requestBody := map[string]interface{}{
		"fnCHANNEL_ID":  channelID,
		"fsASSET_ID":    assetID,
		"fsDEVICE_TYPE": "mobile",
		"clsAPP_IDENTITY_VALIDATE_ARUS": map[string]interface{}{
			"fsVALUE":   FsValue,
			"fsENC_KEY": FsEncKey,
		},
	}
	// 将结构体序列化为 JSON
	jsonData, err := json.MarshalIndent(requestBody, "", "  ")
	if err != nil {
		LogError(err, requestBody)
		return "", err
	}
	// 创建请求主体
	reqBody := bytes.NewBuffer([]byte(jsonData))

	// 创建 HTTP 请求
	req, err := http.NewRequest("POST", UrlGetChannelUrl, reqBody)
	if err != nil {
		LogError(err)
		return "", err
	}

	// 设置请求头
	req.Header = http.Header{
		"Content-Type":    []string{"application/json; charset=UTF-8"},
		"4gtv_auth":       []string{M4gtvAuth},
		"fsdevice":        []string{"Android"},
		"fsversion":       []string{FsVersion},
		"fsenc_key":       []string{FsEncKey},
		"fsvalue":         []string{FsValue},
		"Accept-Encoding": []string{"gzip"},
		"User-Agent":      []string{"okhttp/4.9.2"},
	}

	resp, err := Client1.Do(req)
	if err != nil {
		LogError(err)
		return "", err
	}
	defer resp.Body.Close()
	var reader io.Reader
	if resp.Header.Get("Content-Encoding") == "gzip" {
		gzReader, err := gzip.NewReader(resp.Body)
		if err != nil {
			LogError(err)
			return "", err
		}
		defer gzReader.Close()
		reader = gzReader
	} else {
		reader = resp.Body
	}

	// 使用 strings.Builder 构建响应数据
	var body strings.Builder
	// 将响应体的内容复制到 body 中（忽略错误处理）
	_, _ = io.Copy(&body, reader)
	var response MainResponse
	err = json.Unmarshal([]byte(body.String()), &response)
	if err != nil {
		LogError("访问api时被cf拦截")
		return "", err

	}
	LogDebug("api 获取结果：", body.String())
	if !response.Success {
		LogError("通过api获取playurl失败 错误：assetID：", assetID, response.ErrMessage)
		return "", errors.New("通过api获取playurl失败")
	}
	data := response.Data
	flstURLs, ok := data["flstURLs"].([]interface{})
	if !ok {
		LogError("未找到flstURLs")
		return "", errors.New("未找到flstURLs")
	}

	if len(flstURLs) > 0 {
		LogDebug("api 获取url： ", flstURLs[0])
	} else {
		LogError("flstURLs为空")
		return "", errors.New("flstURLs为空")
	}
	playUrl := flstURLs[0].(string)
	parsedURL, err := url.Parse(playUrl)
	if err != nil {
		LogError("HandleM3u8Raw 失败：", err)
		return "", err
	}
	ProxyUrlWhiteList[parsedURL.Host] = true

	return playUrl, nil
}

func CheckPlayable() {
	LogInfo("开始测试4gtv")
	playUrl := ""
	if len(ChannelsData) == 0 {
		LogError("ChannelsData为空，跳过检测 4gtv api")
		goto NEXT
	}
	playUrl, _ = GetPlayUrl(strconv.Itoa(ChannelsData[0].FnID), ChannelsData[0].Fs4GTV_ID)
	if playUrl == "" {
		LogError("当前ip无法访问 4gtv api")
	} else {
		LogInfo("当前ip可访问 4gtv api")
	}

NEXT:
	//litv
	LogInfo("开始测试Litv")
	timestamp := int(time.Now().Unix()/4 - 355017625)
	t := timestamp * 4
	tsUrl := ""

	for id, ch := range LitvChannels {
		if ch.CachePrefix == "" {
			tsUrl = fmt.Sprintf("%s/live/pool/%s/litv-pc/%s-avc1_6000000=%d-mp4a_134000_zho=%d-begin=%d0000000-dur=40000000-seq=%d.ts",
				UrlLiTVServer,
				id,
				id,
				ch.Param1,
				ch.Param2,
				t,
				timestamp,
			)
		} else {
			tsUrl = fmt.Sprintf("%s/live/pool/%s/litv-pc/%s-begin=%d0000000-dur=40000000-seq=%d.ts",
				UrlLiTVServer,
				id,
				ch.CachePrefix,
				t,
				timestamp,
			)
		}
		break
	}
	LogInfo("当前测试url：", tsUrl)

	statusCode, _, err := MRequest(tsUrl, "GET", nil,
		map[string]string{
			"Referer":    "https://imasdk.googleapis.com",
			"origin":     "https://imasdk.googleapis.com",
			"User-Agent": "Mozilla/5.0 (Linux; Android 12; M2011K2C ; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/134.0.6998.135 Mobile Safari/537.36",
			//"Accept-Encoding": "identity",
		}, false)
	if err != nil || statusCode != 200 {
		LogError("当前ip无法访问 Litv", err)
	} else {
		LogInfo("当前ip可访问 Litv")
	}
}
