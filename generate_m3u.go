package main

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type M3u struct {
}

func (t *M3u) GetM3u(c *gin.Context) {

	host := c.Request.Host

	// var builder strings.Builder

	// // 写入头部
	// builder.WriteString("#EXTM3U x-tvg-url=\"https://assets.livednow.com/epg.xml\"\n")
	// for i, channel := range ChannelsData {

	// 	builder.WriteString("#EXTINF:-1,")
	// 	builder.WriteString("tvg-id=\"" + channel.FsNAME + "\" ")
	// 	builder.WriteString("tvg-name=\"" + channel.FsNAME + "\" ")
	// 	builder.WriteString("tvg-logo=\"" + channel.FsLOGO_MOBILE + "\" ")
	// 	builder.WriteString("group-title=\"" + channel.FsTYPE_NAME + "\",")
	// 	builder.WriteString(channel.FsNAME + "\n")

	// 	// URL 构建
	// 	builder.WriteString("http://")
	// 	builder.WriteString(host)
	// 	builder.WriteString("/4gtv/")
	// 	builder.WriteString(strconv.Itoa(i))
	// 	builder.WriteString(".m3u8?channelid=")
	// 	builder.WriteString(strconv.Itoa(channel.FnID))
	// 	builder.WriteString("&assetid=")
	// 	builder.WriteString(channel.Fs4GTV_ID)
	// 	builder.WriteString("&cdntype=")
	// 	builder.WriteString(channel.FsCDN_ROUTE)
	// 	builder.WriteString("\n")
	// }

	// for _, channel := range channels {
	// 	//litv
	// 	builder.WriteString("#EXTINF:-1,")
	// 	builder.WriteString("tvg-id=\"" + channel.Name + "\" ")
	// 	builder.WriteString("tvg-name=\"" + channel.Name + "\" ")
	// 	builder.WriteString("tvg-logo=\"" + channel.Logo + "\" ")
	// 	builder.WriteString("group-title=\"" + channel.GroupName + "\",")
	// 	builder.WriteString(channel.Name + "\n")

	// 	// URL 构建
	// 	builder.WriteString("http://")
	// 	builder.WriteString(host)
	// 	builder.WriteString("/litv/")
	// 	builder.WriteString(channel.AssetID)
	// 	builder.WriteString(".m3u8")
	// 	builder.WriteString("\n")
	// }

	// 一次性写入响应
	if _, err := c.Writer.Write([]byte(GenerateM3UPlaylist(host))); err != nil {
		//LogError()
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
}

// Channel 定义频道结构
type NChannel struct {
	OldChannel Channel
	ChannelID  int
	CDNType    string
}
type NChannels []NChannel

var MNChannels map[string]NChannels

func BuildChannelMap() {
	MNChannels = make(map[string]NChannels)
	for _, ch := range LitvOnlychannels {
		groupTitles := strings.Split(ch.GroupName, ",")
		for _, groupTitle := range groupTitles {
			MNChannels[groupTitle] = append(MNChannels[groupTitle], NChannel{
				OldChannel: ch,
				ChannelID:  -1,
				CDNType:    "",
			})
		}
	}

	for _, ch := range ChannelsData {
		groupTitles := strings.Split(ch.FsTYPE_NAME, ",")
		for _, groupTitle := range groupTitles {
			MNChannels[groupTitle] = append(MNChannels[groupTitle], NChannel{
				OldChannel: Channel{
					AssetID:   ch.Fs4GTV_ID,
					Name:      ch.FsNAME,
					Logo:      ch.FsLOGO_MOBILE,
					GroupName: groupTitle,
				},
				ChannelID: ch.FnID,
				CDNType:   ch.FsCDN_ROUTE,
			})
		}
	}
}

func GenerateM3UPlaylist(host string) string {

	//buildChannelMap()
	// // 遍历 MNChannels 并打印
	// for groupTitle, nChannels := range MNChannels {
	// 	fmt.Printf("Group: %s\n", groupTitle)
	// 	for i, nChannel := range nChannels {
	// 		fmt.Printf("  Channel %d:\n", i+1)
	// 		fmt.Printf("    AssetID: %s\n", nChannel.OldChannel.AssetID)
	// 		fmt.Printf("    Name: %s\n", nChannel.OldChannel.Name)
	// 		fmt.Printf("    Logo: %s\n", nChannel.OldChannel.Logo)
	// 		fmt.Printf("    GroupName: %s\n", nChannel.OldChannel.GroupName)
	// 		fmt.Printf("    ChannelID: %d\n", nChannel.ChannelID)
	// 		fmt.Printf("    CDNType: %s\n", nChannel.CDNType)
	// 	}
	// 	fmt.Println() // 分组间空行
	// }

	var builder strings.Builder
	var groupTitles = [10]string{
		"綜合",
		"兒童與青少年",
		"音樂綜藝",
		"戲劇、電影與紀錄片",
		"新聞財經",
		"生活旅遊時尚",
		"運動休閒",
		"熱播推薦",
		"國會頻道",
		"現場直擊",
	}
	for _, groupTitle := range groupTitles {
		appendGroupChannels(&builder, groupTitle, host)
	}
	for k := range MNChannels {
		contain := false
		for _, s := range groupTitles {
			if s == k {
				contain = true
				break
			}
		}
		if !contain {
			appendGroupChannels(&builder, k, host)
		}
	}
	return builder.String()
}
func appendGroupChannels(builder *strings.Builder, groupType, host string) *strings.Builder {
	nChannels, exist := MNChannels[groupType]
	if exist {
		for _, nChannel := range nChannels {
			builder.WriteString("#EXTINF:-1,")
			builder.WriteString("tvg-id=\"" + nChannel.OldChannel.Name + "\" ")
			builder.WriteString("tvg-name=\"" + nChannel.OldChannel.Name + "\" ")
			builder.WriteString("tvg-logo=\"" + nChannel.OldChannel.Logo + "\" ")
			builder.WriteString("group-title=\"" + nChannel.OldChannel.GroupName + "\",")
			builder.WriteString(nChannel.OldChannel.Name + "\n")

			// URL 构建
			builder.WriteString("http://")
			builder.WriteString(host)
			if nChannel.ChannelID != -1 {
				builder.WriteString("/4gtv/")
				builder.WriteString(nChannel.OldChannel.AssetID)
				builder.WriteString(".m3u8?channelid=")
				builder.WriteString(strconv.Itoa(nChannel.ChannelID))
				// builder.WriteString("&assetid=")
				// builder.WriteString(nChannel.OldChannel.AssetID)
				builder.WriteString("&cdntype=")
				builder.WriteString(nChannel.CDNType)
				if ValidToken != "" {
					builder.WriteString("&token=")
					builder.WriteString(ValidToken)
				}
			} else {
				builder.WriteString("/litv/")
				builder.WriteString(nChannel.OldChannel.AssetID)
				builder.WriteString(".m3u8")
				if ValidToken != "" {
					builder.WriteString("?token=")
					builder.WriteString(ValidToken)
				}
			}

			builder.WriteString("\n")
		}
	}
	//delete(MNChannels, groupType)
	return builder

}
