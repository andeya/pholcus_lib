package spider_lib

import (
	// 基础包
	"github.com/PuerkitoBio/goquery"                        //DOM解析
	"github.com/henrylee2cn/pholcus/app/downloader/context" //必需
	// "github.com/henrylee2cn/pholcus/logs"           //信息输出
	. "github.com/henrylee2cn/pholcus/app/spider" //必需
	// . "github.com/henrylee2cn/pholcus/app/spider/common" //选用

	// net包
	// "net/http" //设置http.Header
	// "net/url"

	// 编码包
	// "encoding/xml"
	// "encoding/json"

	// 字符串处理包
	// "regexp"
	"strconv"
	// "strings"

	// 其他包
	// "fmt"
	// "math"
	// "time"
)

func init() {
	Zolphone.AddMenu()
}

var Zolphone = &Spider{
	Name:        "中关村手机",
	Description: "中关村苹果手机数据 [Auto Page] [bbs.zol.com.cn/sjbbs/d544_p]",
	// Pausetime: [2]uint{uint(3000), uint(1000)},
	// Keyword:   USE,
	UseCookie: false,
	RuleTree: &RuleTree{
		Root: func(self *Spider) {
			self.Aid("生成请求", map[string]interface{}{"loop": [2]int{1, 950}, "Rule": "生成请求"})
		},

		Trunk: map[string]*Rule{

			"生成请求": {
				AidFunc: func(self *Spider, aid map[string]interface{}) interface{} {
					for loop := aid["loop"].([2]int); loop[0] < loop[1]; loop[0]++ {
						self.AddQueue(map[string]interface{}{
							"Url":  "http://bbs.zol.com.cn/sjbbs/d544_p" + strconv.Itoa(loop[0]) + ".html#c",
							"Rule": aid["Rule"],
						})
					}
					return nil
				},
				ParseFunc: func(self *Spider, resp *context.Response) {
					query := resp.GetDom()
					ss := query.Find("tbody").Find("tr[id]")
					ss.Each(func(i int, goq *goquery.Selection) {
						resp.SetTemp("html", goq)
						self.Parse("获取结果", resp)

					})
				},
			},

			"获取结果": {
				//注意：有无字段语义和是否输出数据必须保持一致
				OutFeild: []string{
					"机型",
					"链接",
					"主题",
					"发表者",
					"发表时间",
					"总回复",
					"总查看",
					"最后回复者",
					"最后回复时间",
				},
				ParseFunc: func(self *Spider, resp *context.Response) {

					selectObj := resp.GetTemp("html").(*goquery.Selection)
					//url
					outUrls := selectObj.Find("td").Eq(1)
					outUrl, _ := outUrls.Attr("data-url")
					outUrl = "http://bbs.zol.com.cn/" + outUrl

					//title type
					outTitles := selectObj.Find("td").Eq(1)
					outType := outTitles.Find(".iclass a").Text()
					outTitle := outTitles.Find("div a").Text()

					//author stime
					authors := selectObj.Find("td").Eq(2)
					author := authors.Find("a").Text()
					stime := authors.Find("span").Text()

					//reply read
					replys := selectObj.Find("td").Eq(3)
					reply := replys.Find("span").Text()
					read := replys.Find("i").Text()

					//ereply etime
					etimes := selectObj.Find("td").Eq(4)
					ereply := etimes.Find("a").Eq(0).Text()
					etime := etimes.Find("a").Eq(1).Text()

					// 结果存入Response中转
					resp.AddItem(map[string]interface{}{
						self.OutFeild(resp, 0): outType,
						self.OutFeild(resp, 1): outUrl,
						self.OutFeild(resp, 2): outTitle,
						self.OutFeild(resp, 3): author,
						self.OutFeild(resp, 4): stime,
						self.OutFeild(resp, 5): reply,
						self.OutFeild(resp, 6): read,
						self.OutFeild(resp, 7): ereply,
						self.OutFeild(resp, 8): etime,
					})
				},
			},
		},
	},
}
