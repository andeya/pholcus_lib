package spider_lib

// 基础包
import (
	"github.com/PuerkitoBio/goquery"                        //DOM解析
	"github.com/henrylee2cn/pholcus/app/downloader/context" //必需
	. "github.com/henrylee2cn/pholcus/app/spider"           //必需
	. "github.com/henrylee2cn/pholcus/app/spider/common"    //选用
	"github.com/henrylee2cn/pholcus/logs"                   //信息输出

	// net包
	"net/http" //设置http.Header
	// "net/url"

	// 编码包
	// "encoding/json"
	"encoding/xml"

	// 字符串处理包
	"regexp"
	// "strconv"
	"strings"

	// 其他包
	// "fmt"
	// "math"
	// "time"
)

func init() {
	BaiduNews.Register()
}

var rss_BaiduNews = map[string]string{
	"国内最新":  "http://news.baidu.com/n?cmd=4&class=civilnews&tn=rss",
	"国际最新":  "http://news.baidu.com/n?cmd=4&class=internews&tn=rss",
	"军事最新":  "http://news.baidu.com/n?cmd=4&class=mil&tn=rss",
	"财经最新":  "http://news.baidu.com/n?cmd=4&class=finannews&tn=rss",
	"互联网最新": "http://news.baidu.com/n?cmd=4&class=internet&tn=rss",
	"房产最新":  "http://news.baidu.com/n?cmd=4&class=housenews&tn=rss",
	"汽车最新":  "http://news.baidu.com/n?cmd=4&class=autonews&tn=rss",
	"体育最新":  "http://news.baidu.com/n?cmd=4&class=sportnews&tn=rss",
	"娱乐最新":  "http://news.baidu.com/n?cmd=4&class=enternews&tn=rss",
	"游戏最新":  "http://news.baidu.com/n?cmd=4&class=gamenews&tn=rss",
	"教育最新":  "http://news.baidu.com/n?cmd=4&class=edunews&tn=rss",
	"女人最新":  "http://news.baidu.com/n?cmd=4&class=healthnews&tn=rss",
	"科技最新":  "http://news.baidu.com/n?cmd=4&class=technnews&tn=rss",
	"社会最新":  "http://news.baidu.com/n?cmd=4&class=socianews&tn=rss",
}

var baiduNewsCountdownTimer = NewCountdownTimer([]float64{5, 10, 20, 30, 45, 60}, func() []string {
	src, i := make([]string, len(rss_BaiduNews)), 0
	for k := range rss_BaiduNews {
		src[i] = k
		i++
	}
	return src
}())

type BaiduNewsData struct {
	Item []BaiduNewsItem `xml:"item"`
}

type BaiduNewsItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
	Author      string `xml:"author"`
}

var BaiduNews = &Spider{
	Name:        "百度RSS新闻",
	Description: "百度RSS新闻，实现轮询更新 [Auto Page] [news.baidu.com]",
	// Pausetime: [2]uint{uint(3000), uint(1000)},
	// Keyword:     USE,
	EnableCookie: false,
	RuleTree: &RuleTree{
		Root: func(self *Spider, resp *context.Response) {
			for k, _ := range rss_BaiduNews {
				self.Aid("LOOP", map[string]interface{}{"loop": k})
			}
		},

		Trunk: map[string]*Rule{
			"LOOP": {
				AidFunc: func(self *Spider, aid map[string]interface{}) interface{} {
					k := aid["loop"].(string)
					v := rss_BaiduNews[k]

					self.AddQueue(&context.Request{
						Url:          v,
						Rule:         "XML列表页",
						Header:       http.Header{"Content-Type": []string{"text/html", "charset=GB2312"}},
						Temp:         map[string]interface{}{"src": k},
						DialTimeout:  -1,
						ConnTimeout:  -1,
						TryTimes:     -1,
						Duplicatable: true,
					})
					return nil
				},
			},
			"XML列表页": {
				ParseFunc: func(self *Spider, resp *context.Response) {
					src := resp.GetTemp("src").(string)
					defer func() {
						// 循环请求
						baiduNewsCountdownTimer.Wait(src)
						self.Aid("LOOP", map[string]interface{}{"loop": src})
					}()

					page := GBKToUTF8(resp.GetText())
					page = strings.TrimLeft(page, `<?xml version="1.0" encoding="gb2312"?>`)
					re, _ := regexp.Compile(`\<[\/]?rss\>`)
					page = re.ReplaceAllString(page, "")

					content := new(BaiduNewsData)
					if err := xml.Unmarshal([]byte(page), content); err != nil {
						logs.Log.Error("XML列表页: %v", err)
						return
					}

					for _, v := range content.Item {
						self.AddQueue(&context.Request{
							Url:  v.Link,
							Rule: "新闻详情",
							Temp: map[string]interface{}{
								"title":       CleanHtml(v.Title, 4),
								"description": CleanHtml(v.Description, 4),
								"src":         src,
								"releaseTime": CleanHtml(v.PubDate, 4),
								"author":      CleanHtml(v.Author, 4),
							},
						})
					}
				},
			},

			"新闻详情": {
				//注意：有无字段语义和是否输出数据必须保持一致
				OutFeild: []string{
					"标题",
					"描述",
					"内容",
					"发布时间",
					"分类",
					"作者",
				},
				ParseFunc: func(self *Spider, resp *context.Response) {
					// RSS标记更新
					baiduNewsCountdownTimer.Update(resp.GetTemp("src").(string))

					title := resp.GetTemp("title").(string)

					infoStr, isReload := baiduNewsFn.prase(resp)
					if isReload {
						return
					}
					// 结果存入Response中转
					self.Output(resp.GetRuleName(), resp, map[int]interface{}{
						0: title,
						1: resp.GetTemp("description"),
						2: infoStr,
						3: resp.GetTemp("releaseTime"),
						4: resp.GetTemp("src"),
						5: resp.GetTemp("author"),
					})
				},
			},
		},
	},
}

type baiduNews map[string]func(resp *context.Response) (infoStr string, isReload bool)

// @url 必须为含有协议头的地址
func (b baiduNews) prase(resp *context.Response) (infoStr string, isReload bool) {
	url := resp.Response.Request.URL.Host
	// Log.Println("域名", url)
	if _, ok := b[url]; ok {
		return b[url](resp)
	} else {
		return b.commonPrase(resp), false
	}
}

func (b baiduNews) commonPrase(resp *context.Response) (infoStr string) {
	body := resp.GetDom().Find("body")

	var info *goquery.Selection

	if h1s := body.Find("h1"); len(h1s.Nodes) != 0 {
		for i := 0; i < len(h1s.Nodes); i++ {
			info = b.findP(h1s.Eq(i))
		}
	} else if h2s := body.Find("h2"); len(h2s.Nodes) != 0 {
		for i := 0; i < len(h2s.Nodes); i++ {
			info = b.findP(h2s.Eq(i))
		}
	} else if h3s := body.Find("h3"); len(h3s.Nodes) != 0 {
		for i := 0; i < len(h3s.Nodes); i++ {
			info = b.findP(h3s.Eq(i))
		}
	} else {
		info = body.Find("body")
	}
	// 去除标签
	// info.RemoveFiltered("script")
	// info.RemoveFiltered("style")
	infoStr, _ = info.Html()

	// 清洗HTML
	infoStr = CleanHtml(infoStr, 5)
	return
}

func (b baiduNews) findP(html *goquery.Selection) *goquery.Selection {
	if html.Is("body") {
		return html
	} else if result := html.Parent().Find("p"); len(result.Nodes) == 0 {
		return b.findP(html.Parent())
	} else {
		return html.Parent()
	}
}

var baiduNewsFn = baiduNews{
	"yule.sohu.com": func(resp *context.Response) (infoStr string, isReload bool) {

		// 当有翻页等需求时，重新添加请求
		// req := resp.GetRequest()

		// 根据需要可能会对req进行某些修改
		// req.SetUrl("")

		// 添加请求到队列
		// scheduler.Sdl.Push(req)
		infoStr = resp.GetDom().Find("#contentText").Text()

		return
	},
	"news.qtv.com.cn": func(resp *context.Response) (infoStr string, isReload bool) {
		infoStr = resp.GetDom().Find(".zwConreally_z").Text()
		return
	},
}
