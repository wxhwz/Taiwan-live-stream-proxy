package main

import (
	"flag"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var m3uObj M3u
var m4gtvObj M4gtv

// 鉴权中间件
func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		if ValidToken != "" {
			token := c.Query("token")

			// 验证Token
			if token != ValidToken {
				c.JSON(401, gin.H{
					"error": "Invalid or missing token",
				})
				c.Abort()
				return
			}
		}

		// Token有效，继续处理请求
		c.Next()
	}
}

// 设置路由和处理逻辑
func setupRouter() *gin.Engine {
	// 设置Gin为发布模式
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// 创建需要鉴权的路由组
	authorized := r.Group("/")
	authorized.Use(authMiddleware())

	// 配置获取tv.m3u文件的路由
	authorized.GET("/", func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "application/octet-stream")
		c.Writer.Header().Set("Content-Disposition", "attachment; filename=4gtv.m3u")
		m3uObj.GetM3u(c)
	})

	authorized.GET("/4gtv/:rid", func(c *gin.Context) {
		ts := c.Query("ts")
		if ts == "" {
			rid := c.Param("rid")
			channelID := c.Query("channelid")
			assetID := strings.ReplaceAll(rid, ".m3u8", "")
			cdnType := c.Query("cdntype")

			if channelID != "" && assetID != "" {
				m4gtvObj.HandleMainRequest(c, channelID, assetID, cdnType)
			} else {
				c.JSON(400, gin.H{"error": "Missing required parameters"})
			}
		} else {
			fullPath := c.Request.URL.String()
			m4gtvObj.HandleTsRequest(c, fullPath[strings.Index(fullPath, "?ts")+4:])
		}

	})

	authorized.GET("/litv/:rid", func(c *gin.Context) {
		ts := c.Query("ts")
		if ts == "" {
			// assetID := c.Param("rid")
			// assetID = strings.ReplaceAll(assetID, ".m3u8", "")
			// _, exists := LitvChannels[assetID]
			// if exists {
			// 	c.Data(200, "application/vnd.apple.mpegurl", LitvGenerateM3U8(assetID))
			// } else {
			// 	c.String(404, "ID不存在")
			// }

			rid := c.Param("rid")
			assetID := strings.ReplaceAll(rid, ".m3u8", "")
			if assetID != "" {
				m4gtvObj.HandleMainRequest(c, "", assetID, "")
			} else {
				c.JSON(400, gin.H{"error": "Missing required parameters"})
			}
		} else {
			fullPath := c.Request.URL.String()
			m4gtvObj.HandleTsRequest(c, fullPath[strings.Index(fullPath, "?ts")+4:])
		}
	})

	return r
}
func main() {
	host := flag.String("host", "0.0.0.0", "host")
	port := flag.String("p", "18089", "port")
	flag.BoolVar(&ProxyTs, "proxyts", false, "开启TS代理")
	flag.StringVar(&ValidToken, "token", "", "设置token鉴权")
	flag.BoolVar(&DebugMode, "debug", false, "开启调试模式")
	flag.BoolVar(&PreCache, "precache", false, "开启预先缓存")
	flag.StringVar(&FsUSER, "user", "", "用户名")
	flag.StringVar(&FsPASSWORD, "password", "", "密码")
	flag.Parse()

	InitUID()
	//PushDevice()
	GetAllChannels()
	Get4gtvAuth()

	if FsUSER != "" && FsPASSWORD != "" {
		SignIn()
	}

	//fmt.Printf("http://127.0.0.1:%s\n", *port)
	LitvInitChannelsFromFile()

	// 创建一个通道用于停止定时任务
	done := make(chan bool)
	// 启动定时任务（goroutine）
	go timedFunction(done)

	r := setupRouter()
	LogInfo("\n")
	LogInfo("对当前ip进行检测")
	CheckPlayable()
	LogInfo("检测完毕\n")

	LogInfo("可通过 -h 查看帮助")
	LogInfo("Listen on "+*host+":"+*port, "...")
	if ValidToken != "" {
		LogInfo("使用浏览器访问 http://ip:" + *port + "?token=" + ValidToken + " 获取m3u文件")
	} else {
		LogInfo("使用浏览器访问 http://ip:" + *port + " 获取m3u文件")
	}

	r.Run(*host + ":" + *port)
	done <- true // 发送停止信号
}

// 定时执行的函数
func timedFunction(done <-chan bool) {

	i := 0
	// 使用 Timer 控制首次执行
	timer := time.NewTimer(1800 * time.Second)
	defer timer.Stop()

	// 等待首次执行
	select {
	case <-timer.C:
		// 并发执行首次任务
		//fmt.Println("定时任务执行")
		i += 1
		RefreshCache(PreCache, true)

	case <-done:
		//fmt.Println("程序在首次执行前退出")
		return
	}

	// 创建一个定时器，每隔 ? 秒触发一次
	ticker := time.NewTicker(10800 * time.Second)
	defer ticker.Stop() // 确保结束时释放资源

	for {
		select {
		case <-done:
			// 收到停止信号，退出函数
			return
		case <-ticker.C:
			// 这是定时执行的业务逻辑
			//fmt.Println("定时任务执行")
			i = i % 8
			if i == 0 {
				RefreshCache(PreCache, true)
			} else {
				RefreshCache(PreCache, false)
			}
			i += 1
		}
	}

}
