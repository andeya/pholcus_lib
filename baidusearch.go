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
	BaiduSearch.AddMenu()
}

var BaiduSearch = &Spider{
	Name:        "百度搜索",
	Description: "百度搜索结果 [www.baidu.com]",
	// Pausetime: [2]uint{uint(3000), uint(1000)},
	Keyword:   USE,
	UseCookie: false,
	RuleTree: &RuleTree{
		Root: func(self *Spider) {
			self.Aid("生成请求", map[string]interface{}{"loop": [2]int{0, 1}, "Rule": "生成请求"})
		},

		Trunk: map[string]*Rule{

			"生成请求": {
				AidFunc: func(self *Spider, aid map[string]interface{}) interface{} {
					for loop := aid["loop"].([2]int); loop[0] < loop[1]; loop[0]++ {
						self.AddQueue(map[string]interface{}{
							"Url":  "http://www.baidu.com/s?ie=utf-8&nojc=1&wd=" + self.GetKeyword() + "&rn=50&pn=" + strconv.Itoa(50*loop[0]),
							"Rule": aid["Rule"],
						})
					}
					return nil
				},
				ParseFunc: func(self *Spider, resp *context.Response) {
					query := resp.GetDom()
					total1 := query.Find(".nums").Text()
					re, _ := regexp.Compile(`[\D]*`)
					total1 = re.ReplaceAllString(total1, "")
					total2, _ := strconv.Atoi(total1)
					total := int(math.Ceil(float64(total2) / 50))
					if total > self.MaxPage {
						total = self.MaxPage
					} else if total == 0 {
						logs.Log.Critical("[消息提示：| 任务：%v | 关键词：%v | 规则：%v] 没有抓取到任何数据！!!\n", self.GetName(), self.GetKeyword(), resp.GetRuleName())
						return
					}
					// 调用指定规则下辅助函数
					self.Aid("生成请求", map[string]interface{}{"loop": [2]int{1, total}, "Rule": "搜索结果"})
					// 用指定规则解析响应流
					self.Parse("搜索结果", resp)
				},
			},

			"搜索结果": {
				//注意：有无字段语义和是否输出数据必须保持一致
				OutFeild: []string{
					"标题",
					"内容",
					"不完整URL",
					"百度跳转",
				},
				ParseFunc: func(self *Spider, resp *context.Response) {
					query := resp.GetDom()
					query.Find("#content_left .c-container").Each(func(i int, s *goquery.Selection) {

						title := s.Find(".t").Text()
						content := s.Find(".c-abstract").Text()
						href, _ := s.Find(".t >a").Attr("href")
						tar := s.Find(".g").Text()

						re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
						// title = re.ReplaceAllStringFunc(title, strings.ToLower)
						// content = re.ReplaceAllStringFunc(content, strings.ToLower)

						title = re.ReplaceAllString(title, "")
						content = re.ReplaceAllString(content, "")

						// 结果存入Response中转
						resp.AddItem(map[string]interface{}{
							self.OutFeild(resp, 0): strings.Trim(title, " \t\n"),
							self.OutFeild(resp, 1): strings.Trim(content, " \t\n"),
							self.OutFeild(resp, 2): tar,
							self.OutFeild(resp, 3): href,
						})
					})
				},
			},
		},
	},
}
