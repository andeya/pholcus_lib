package spider_lib

// 基础包
import (
	"github.com/PuerkitoBio/goquery"                        //DOM解析
	"github.com/henrylee2cn/pholcus/app/downloader/context" //必需
	// "github.com/henrylee2cn/pholcus/logs"               //信息输出
	. "github.com/henrylee2cn/pholcus/app/spider" //必需
	// . "github.com/henrylee2cn/pholcus/app/spider/common"          //选用

	// net包
	// "net/http" //设置http.Header
	// "net/url"

	// 编码包
	// "encoding/xml"
	// "encoding/json"

	// 字符串处理包
	// "regexp"
	"strconv"
	"strings"

	// 其他包
	// "fmt"
	// "math"
	// "time"
)

func init() {
	GanjiGongsi.AddMenu()
}

var GanjiGongsi = &Spider{
	Name:        "企业名录-赶集网",
	Description: "企业名录-深圳-赶集网 [www.ganji.com/gongsi]",
	// Pausetime: [2]uint{uint(3000), uint(1000)},
	// Keyword:   USE,
	UseCookie: false,
	RuleTree: &RuleTree{
		Root: func(self *Spider) {
			self.AddQueue(map[string]interface{}{
				"Url":  "http://sz.ganji.com/gongsi/o1",
				"Rule": "请求列表",
				"Temp": map[string]interface{}{"p": 1},
			})
		},

		Trunk: map[string]*Rule{

			"请求列表": {
				ParseFunc: func(self *Spider, resp *context.Response) {
					curr := resp.GetTemp("p").(int)
					if resp.GetDom().Find(".linkOn span").Text() != strconv.Itoa(curr) {
						return
					}
					self.AddQueue(map[string]interface{}{
						"Url":  "http://sz.ganji.com/gongsi/o" + strconv.Itoa(curr+1),
						"Rule": "请求列表",
						"Temp": map[string]interface{}{"p": curr + 1},
					})

					// 用指定规则解析响应流
					self.Parse("获取列表", resp)
				},
			},

			"获取列表": {
				ParseFunc: func(self *Spider, resp *context.Response) {
					resp.GetDom().
						Find(".com-list-2 table a").
						Each(func(i int, s *goquery.Selection) {
						url, _ := s.Attr("href")
						self.AddQueue(map[string]interface{}{
							"Url":  url,
							"Rule": "输出结果",
						})
					})
				},
			},

			"输出结果": {
				//注意：有无字段语义和是否输出数据必须保持一致
				OutFeild: []string{
					"公司",
					"联系人",
					"地址",
					"简介",
					"行业",
					"类型",
					"规模",
				},
				ParseFunc: func(self *Spider, resp *context.Response) {
					query := resp.GetDom()

					var 公司, 规模, 行业, 类型, 联系人, 地址 string

					query.Find(".c-introduce li").Each(func(i int, s *goquery.Selection) {
						em := s.Find("em").Text()
						t := strings.Split(s.Text(), `   `)[0]
						t = strings.Replace(t, em, "", -1)
						t = strings.Trim(t, " ")

						switch em {
						case "公司名称：":
							公司 = t

						case "公司规模：":
							规模 = t

						case "公司行业：":
							行业 = t

						case "公司类型：":
							类型 = t

						case "联 系 人：":
							联系人 = t

						case "联系电话：":
							if img, ok := s.Find("img").Attr("src"); ok {
								self.AddQueue(map[string]interface{}{
									"Url":      "http://www.ganji.com" + img,
									"Rule":     "联系方式",
									"Temp":     map[string]interface{}{"n": 公司 + "(" + 联系人 + ").png"},
									"Priority": 1,
								})
							}

						case "公司地址：":
							地址 = t
						}
					})

					简介 := query.Find("#company_description").Text()

					// 结果存入Response中转
					resp.AddItem(map[string]interface{}{
						self.OutFeild(resp, 0): 公司,
						self.OutFeild(resp, 1): 联系人,
						self.OutFeild(resp, 2): 地址,
						self.OutFeild(resp, 3): 简介,
						self.OutFeild(resp, 4): 行业,
						self.OutFeild(resp, 5): 类型,
						self.OutFeild(resp, 6): 规模,
					})
				},
			},

			"联系方式": {
				ParseFunc: func(self *Spider, resp *context.Response) {
					resp.AddFile(resp.GetTemp("n").(string))
				},
			},
		},
	},
}
