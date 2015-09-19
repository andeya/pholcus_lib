package spider_lib

// 基础包
import (
	"github.com/PuerkitoBio/goquery"                        //DOM解析
	"github.com/henrylee2cn/pholcus/app/downloader/context" //必需
	. "github.com/henrylee2cn/pholcus/app/spider"           //必需
	"github.com/henrylee2cn/pholcus/logs"                   //信息输出
	// . "github.com/henrylee2cn/pholcus/app/spider/common"          //选用

	// net包
	// "net/http" //设置http.Header
	// "net/url"

	// 编码包
	// "encoding/xml"
	// "encoding/json"

	// 字符串处理包
	"regexp"
	"strconv"
	"strings"

	// 其他包
	// "fmt"
	"math"
	// "time"
)

func init() {
	GoogleSearch.AddMenu()
}

var googleIp = []string{
	"210.242.125.100",
	"210.242.125.96",
	"210.242.125.91",
	"210.242.125.95",
	"64.233.189.163",
	"58.123.102.5",
	"210.242.125.97",
	"210.242.125.115",
	"58.123.102.28",
	"210.242.125.70",
}

var GoogleSearch = &Spider{
	Name:        "谷歌搜索",
	Description: "谷歌搜索结果 [www.google.com镜像]",
	// Pausetime: [2]uint{uint(3000), uint(1000)},
	Keyword:   USE,
	UseCookie: false,
	RuleTree: &RuleTree{
		Root: func(self *Spider) {
			var url string
			var success bool
			logs.Log.Critical("正在查找可用的Google镜像，该过程可能需要几分钟……")
			for _, ip := range googleIp {
				url = "http://" + ip + "/search?q=" + self.GetKeyword() + "&newwindow=1&biw=1600&bih=398&start="
				if _, err := goquery.NewDocument(url); err == nil {
					success = true
					break
				}
			}
			if !success {
				logs.Log.Critical("没有可用的Google镜像IP！！")
				return
			}
			logs.Log.Critical("开始Google搜索……")
			self.AddQueue(map[string]interface{}{
				"Url":  url,
				"Rule": "获取总页数",
				"Temp": map[string]interface{}{
					"baseUrl": url,
				},
			})
		},

		Trunk: map[string]*Rule{

			"获取总页数": {
				AidFunc: func(self *Spider, aid map[string]interface{}) interface{} {
					for loop := aid["loop"].([2]int); loop[0] < loop[1]; loop[0]++ {
						self.AddQueue(map[string]interface{}{
							"Url":  aid["urlBase"].(string) + strconv.Itoa(10*loop[0]),
							"Rule": aid["Rule"],
						})
					}
					return nil
				},
				ParseFunc: func(self *Spider, resp *context.Response) {
					query := resp.GetDom()
					txt := query.Find("#resultStats").Text()
					re, _ := regexp.Compile(`,+`)
					txt = re.ReplaceAllString(txt, "")
					re, _ = regexp.Compile(`[\d]+`)
					txt = re.FindString(txt)
					num, _ := strconv.Atoi(txt)
					total := int(math.Ceil(float64(num) / 10))
					if total > self.MaxPage {
						total = self.MaxPage
					} else if total == 0 {
						logs.Log.Critical("[消息提示：| 任务：%v | 关键词：%v | 规则：%v] 没有抓取到任何数据！!!\n", self.GetName(), self.GetKeyword(), resp.GetRuleName())
						return
					}
					// 调用指定规则下辅助函数
					self.Aid("获取总页数", map[string]interface{}{
						"loop":    [2]int{1, total},
						"urlBase": resp.GetTemp("baseUrl"),
						"Rule":    "搜索结果",
					})
					// 用指定规则解析响应流
					self.Parse("搜索结果", resp)
				},
			},

			"搜索结果": {
				//注意：有无字段语义和是否输出数据必须保持一致
				OutFeild: []string{
					"标题",
					"内容",
					"链接",
				},
				ParseFunc: func(self *Spider, resp *context.Response) {
					query := resp.GetDom()
					query.Find("#ires li.g").Each(func(i int, s *goquery.Selection) {
						t := s.Find(".r > a")
						href, _ := t.Attr("href")
						href = strings.TrimLeft(href, "/url?q=")
						title := t.Text()
						content := s.Find(".st").Text()
						resp.AddItem(map[string]interface{}{
							self.OutFeild(resp, 0): title,
							self.OutFeild(resp, 1): content,
							self.OutFeild(resp, 2): href,
						})
					})
				},
			},
		},
	},
}
