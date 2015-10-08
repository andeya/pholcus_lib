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
	AlibabaProduct.AddMenu()
}

var AlibabaProduct = &Spider{
	Name:        "阿里巴巴产品搜索",
	Description: "阿里巴巴产品搜索 [s.1688.com/selloffer/offer_search.htm]",
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
					keyword := EncodeString(self.GetKeyword(), "GBK")
					for loop := aid["loop"].([2]int); loop[0] < loop[1]; loop[0]++ {
						self.AddQueue(map[string]interface{}{
							"Url":    "http://s.1688.com/selloffer/offer_search.htm?enableAsync=false&earseDirect=false&button_click=top&pageSize=60&n=y&offset=3&uniqfield=pic_tag_id&keywords=" + keyword + "&beginPage=" + strconv.Itoa(loop[0]+1),
							"Rule":   aid["Rule"],
							"Header": http.Header{"Content-Type": []string{"text/html", "charset=GBK"}},
						})
					}
					return nil
				},
				ParseFunc: func(self *Spider, resp *context.Response) {
					query := resp.GetDom()
					pageTag := query.Find("#sm-pagination div[data-total-page]")
					// 跳转
					if len(pageTag.Nodes) == 0 {
						logs.Log.Critical("[消息提示：| 任务：%v | 关键词：%v | 规则：%v] 由于跳转AJAX问题，目前只能每个子类抓取 1 页……\n", self.GetName(), self.GetKeyword(), resp.GetRuleName())
						query.Find(".sm-floorhead-typemore a").Each(func(i int, s *goquery.Selection) {
							if href, ok := s.Attr("href"); ok {
								self.AddQueue(map[string]interface{}{
									"Url":    href,
									"Header": http.Header{"Content-Type": []string{"text/html", "charset=GBK"}},
									"Rule":   "搜索结果",
								})
							}
						})
						return
					}
					total1, _ := pageTag.First().Attr("data-total-page")
					total1 = strings.Trim(total1, " \t\n")
					total, _ := strconv.Atoi(total1)
					if total > self.MaxPage {
						total = self.MaxPage
					} else if total == 0 {
						logs.Log.Critical("[消息提示：| 任务：%v | 关键词：%v | 规则：%v] 没有抓取到任何数据！！！\n", self.GetName(), self.GetKeyword(), resp.GetRuleName())
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
					"公司",
					"标题",
					"价格",
					"销量",
					"星级",
					"地址",
					"链接",
				},
				ParseFunc: func(self *Spider, resp *context.Response) {
					query := resp.GetDom()

					query.Find("#sm-offer-list > li").Each(func(i int, s *goquery.Selection) {

						// 获取公司
						company, _ := s.Find("a.sm-offer-companyName").First().Attr("title")

						// 获取标题
						t := s.Find(".sm-offer-title > a:nth-child(1)")
						title, _ := t.Attr("title")

						// 获取URL
						url, _ := t.Attr("href")

						// 获取价格
						price := s.Find(".sm-offer-priceNum").First().Text()

						// 获取成交量
						sales := s.Find("span.sm-offer-trade > em").First().Text()

						// 获取地址
						address, _ := s.Find(".sm-offer-location").First().Attr("title")

						// 获取信用年限
						level := s.Find("span.sm-offer-companyTag > a.sw-ui-flaticon-cxt16x16").First().Text()

						// 结果存入Response中转
						resp.AddItem(map[string]interface{}{
							self.OutFeild(resp, 0): company,
							self.OutFeild(resp, 1): title,
							self.OutFeild(resp, 2): price,
							self.OutFeild(resp, 3): sales,
							self.OutFeild(resp, 4): level,
							self.OutFeild(resp, 5): address,
							self.OutFeild(resp, 6): url,
						})
					})
				},
			},
		},
	},
}
